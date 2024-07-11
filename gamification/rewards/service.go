package rewards

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"milestone_core/shared/formats"
	"milestone_core/shared/sql"
)

type Service struct {
	DbConnection *sqlx.DB
}

func (s Service) CreateReward(workspaceId string, reward Reward) error {
	existingId, err := s.KeyToId(workspaceId, reward.Key)
	if err != nil {
		return err
	}
	if existingId != nil && *existingId != "" {
		return Errors.KeyExistsError
	}
	if !s.isValidType(reward.Type) {
		return NewRewardError(400, "invalid reward type")
	}

	rawRules := make([]RawRule, len(reward.Rules))
	for i, rule := range reward.Rules {
		rawRules[i] = RawRule{
			Rule:  rule,
			Value: formats.EncodeJson(rule.Value),
		}
	}
	if reward.Rules, err = s.mapRawRules(rawRules); err != nil {
		return err
	}

	decodedOptions, err := s.mapRawOptions(reward.Type, reward.RawOptions)
	if err != nil {
		return Errors.InvalidFormat
	}

	tx := s.DbConnection.MustBegin()
	var rewardId string
	err = tx.QueryRow("INSERT INTO game_engine.reward (workspace_id, key, name, type, metadata, options) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		workspaceId, reward.Key, reward.Name, reward.Type, reward.Metadata, s.encodeOptions(decodedOptions),
	).Scan(&rewardId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(reward.Rules) == 0 {
		tx.Commit()
		return nil
	}

	valuesToInsert := make([]interface{}, 0)
	for _, rule := range reward.Rules {
		valuesToInsert = append(valuesToInsert, rewardId, rule.Condition, formats.EncodeJson(rule.Value))
	}
	query := "INSERT INTO game_engine.reward_rule (reward_id, condition, value) VALUES "
	for i := 0; i < len(reward.Rules); i++ {
		query += fmt.Sprintf("($%d, $%d, $%d),", i*3+1, i*3+2, i*3+3)
	}
	query = query[:len(query)-1]
	_, err = tx.Exec(query, valuesToInsert...)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (s Service) GetRewards(workspaceId string) ([]Reward, error) {
	rewards, err := sql.FetchMultiple[Reward](s.DbConnection, "SELECT id, key, name, type, metadata, options FROM game_engine.reward WHERE workspace_id = $1 AND deleted_at IS NULL", workspaceId)
	return rewards, err
}

func (s Service) KeyToId(workspaceId string, key string) (*string, error) {
	id, err := sql.FetchOne[string](s.DbConnection, "SELECT id FROM game_engine.reward WHERE workspace_id = $1 AND key = $2 AND deleted_at IS NULL LIMIT 1", workspaceId, key)
	return id, err
}

func (s Service) GetRewardById(workspaceId string, id string) (*Reward, error) {
	reward, err := sql.FetchOne[Reward](s.DbConnection, "SELECT id, key, name, type, metadata, options FROM game_engine.reward WHERE workspace_id = $1 AND id = $2 AND deleted_at IS NULL", workspaceId, id)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, nil
	}

	rawRules, err := sql.FetchMultiple[RawRule](s.DbConnection, "SELECT id, condition, value FROM game_engine.reward_rule WHERE reward_id = $1", reward.ID)
	if err != nil {
		return nil, err
	}
	reward.Rules, err = s.mapRawRules(rawRules)

	return reward, err
}

func (s Service) UpdateReward(workspaceId string, id string, updatedReward *Reward) error {
	row, err := sql.FetchOne[string](s.DbConnection, "SELECT id FROM game_engine.reward WHERE workspace_id = $1 AND id = $2 AND deleted_at IS NULL LIMIT 1", workspaceId, id)
	if err != nil {
		return err
	}
	if row == nil {
		return Errors.NotFound
	}

	existingId, err := s.KeyToId(workspaceId, updatedReward.Key)
	if err != nil {
		return err
	}
	if existingId != nil && *existingId != id {
		return Errors.KeyExistsError
	}
	if !s.isValidType(updatedReward.Type) {
		return NewRewardError(400, "invalid reward type")
	}

	decodedOptions, err := s.mapRawOptions(updatedReward.Type, updatedReward.RawOptions)
	if err != nil {
		return Errors.InvalidFormat
	}

	tx := s.DbConnection.MustBegin()
	_, err = tx.Exec("UPDATE game_engine.reward SET key = $1, name = $2, type = $3, metadata = $4, options = $5 WHERE workspace_id = $6 AND id = $7 AND deleted_at IS NULL",
		updatedReward.Key, updatedReward.Name, updatedReward.Type, updatedReward.Metadata, s.encodeOptions(decodedOptions), workspaceId, id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if updatedReward.Rules == nil {
		tx.Commit()
		return nil
	}

	rawRules := make([]RawRule, len(updatedReward.Rules))
	for i, rule := range updatedReward.Rules {
		rawRules[i] = RawRule{
			Rule:  rule,
			Value: formats.EncodeJson(rule.Value),
		}
	}
	if updatedReward.Rules, err = s.mapRawRules(rawRules); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM game_engine.reward_rule WHERE reward_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(updatedReward.Rules) == 0 {
		tx.Commit()
		return nil
	}

	valuesToInsert := make([]interface{}, 0)
	for _, rule := range updatedReward.Rules {
		valuesToInsert = append(valuesToInsert, id, rule.Condition, formats.EncodeJson(rule.Value))
	}
	query := "INSERT INTO game_engine.reward_rule (reward_id, condition, value) VALUES "
	for i := 0; i < len(updatedReward.Rules); i++ {
		query += fmt.Sprintf("($%d, $%d, $%d),", i*3+1, i*3+2, i*3+3)
	}
	query = query[:len(query)-1]
	_, err = tx.Exec(query, valuesToInsert...)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (s Service) DeleteReward(workspaceId string, id string) error {
	currentReward, err := s.GetRewardById(workspaceId, id)
	if err != nil {
		return err
	}
	if currentReward == nil {
		return nil
	}

	tx := s.DbConnection.MustBegin()
	_, err = tx.Exec("SELECT game_engine.delete_reward_record($1)", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s Service) mapRawRules(rawRules []RawRule) ([]Rule, error) {
	rules := make([]Rule, len(rawRules))
	var err error
	for i, rawRule := range rawRules {
		rules[i] = rawRule.Rule

		switch rawRule.Condition {
		case RuleConditionEventOccurred:
			rules[i].Value, err = formats.DecodeJson[RuleValueEventOccurred](rawRule.Value)
			val := rules[i].Value.(*RuleValueEventOccurred)
			if err != nil || val.EventKey == nil || val.Times == nil {
				return nil, NewRewardError(400, string("invalid value for rule: "+rules[i].Condition))
			}
		default:
			return nil, NewRewardError(400, string("unknown rule condition: "+rawRule.Condition))
		}
	}
	return rules, nil
}

func (s Service) mapRawOptions(rewardType RewardType, rawOptions *json.RawMessage) (RewardOptions, error) {
	if rawOptions == nil {
		return nil, NewRewardError(400, "options are required")
	}

	switch rewardType {
	case RewardTypePoints:
		decodedVal, err := formats.DecodeJson[PointsRewardOptions](*rawOptions)
		if err != nil || decodedVal.Repeatable == nil || decodedVal.PointsAmount == nil {
			return nil, NewRewardError(400, string("invalid options for reward type: "+rewardType))
		}
		return decodedVal, nil
	case RewardTypeLevel:
		decodedVal, err := formats.DecodeJson[LevelRewardOptions](*rawOptions)
		if err != nil || decodedVal.Level == nil {
			return nil, NewRewardError(400, string("invalid options for reward type: "+rewardType))
		}
		return decodedVal, nil
	case RewardTypeBadge, RewardTypeCustom:
		return &CommonRewardOptions{}, nil
	default:
		break
	}

	return nil, NewRewardError(400, "unknown reward type")
}

func (s Service) encodeOptions(options RewardOptions) *json.RawMessage {
	switch v := options.(type) {
	case *PointsRewardOptions:
		encodedVal := formats.EncodeJson(v)
		return (*json.RawMessage)(&encodedVal)
	case *LevelRewardOptions:
		encodedVal := formats.EncodeJson(v)
		return (*json.RawMessage)(&encodedVal)
	case *CommonRewardOptions:
		emptyStruct := struct{}{}
		encodedVal := formats.EncodeJson(emptyStruct)
		return (*json.RawMessage)(&encodedVal)
	default:
		break
	}

	return nil
}

func (s Service) isValidType(rewardType RewardType) bool {
	switch rewardType {
	case RewardTypeBadge, RewardTypePoints, RewardTypeLevel, RewardTypeCustom:
		return true
	}
	return false
}

type RawRule struct {
	Rule
	Value json.RawMessage `db:"value"`
}
