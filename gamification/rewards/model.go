package rewards

import "encoding/json"

type Reward struct {
	ID          string           `json:"id"  db:"id"`
	WorkspaceID string           `json:"-" db:"workspace_id"`
	Name        string           `json:"name" db:"name"`
	Key         string           `json:"key" db:"key"`
	Type        RewardType       `json:"type" db:"type"`
	Metadata    *json.RawMessage `json:"metadata" db:"metadata"`
	Rules       []Rule           `json:"rules" db:"rules"`
	RawOptions  *json.RawMessage `json:"options" db:"options"`
	Options     RewardOptions    `json:"-" db:"-"`
}

type RewardType string

const (
	RewardTypeBadge  RewardType = "badge"
	RewardTypePoints RewardType = "points"
	RewardTypeLevel  RewardType = "level"
	RewardTypeCustom RewardType = "custom"
)

type Rule struct {
	ID        string        `json:"id"  db:"id"`
	Condition RuleCondition `json:"condition" bson:"condition"`
	Value     RuleValue     `json:"value" bson:"value"`
}

type RuleValue = interface{}
type RuleValueEventOccurred struct {
	RuleValue `json:"-"`
	EventKey  *string `json:"event_key"`
	Times     *int    `json:"times"`
}
type RuleValueBadgeUnlocked struct {
	RuleValue `json:"-"`
	BadgeKey  *string `json:"badge_key"`
}

type RuleCondition string

const (
	RuleConditionEventOccurred RuleCondition = "EventOccurred"
)

type RewardOptions interface {
	IsRewardOptions() bool
}

type CommonRewardOptions struct {
}

type LevelRewardOptions struct {
	Level *string `json:"level"`
}

type PointsRewardOptions struct {
	PointsAmount *int  `json:"points_amount"`
	Repeatable   *bool `json:"repeatable"`
}

func (c CommonRewardOptions) IsRewardOptions() bool {
	return true
}

func (l LevelRewardOptions) IsRewardOptions() bool {
	return true
}

func (p PointsRewardOptions) IsRewardOptions() bool {
	return true
}
