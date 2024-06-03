package flow

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Enroller struct {
	Collection *mongo.Collection
}

type EnrollmentOpts struct {
	SkippedIds          []string
	FinishedIds         []string
	CurrentEnrollmentId string
	SignUpTimestamp     int64
	UserSegment         string
}

func (s *Enroller) GetFlow(workspaceId string, opts EnrollmentOpts) (*Flow, error) {
	queryOpts, err := s.buildQueryOpts(workspaceId, opts)
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(queryOpts)
	log.Default().Print(string(jsonData))

	var flow Flow
	err = s.Collection.FindOne(context.Background(), queryOpts).Decode(&flow)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (s *Enroller) buildQueryOpts(workspaceId string, opts EnrollmentOpts) (bson.M, error) {
	queryOpts := bson.M{"$and": []bson.M{}}
	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), bson.M{
		"workspaceId": workspaceId,
		"live":        true,
	})

	if opts.CurrentEnrollmentId != "" {
		queryOpts = s.withCurrentEnrollmentIdCondition(queryOpts, opts.CurrentEnrollmentId)
		return queryOpts, nil
	}

	queryOpts = s.withDependsOnCondition(queryOpts, opts.FinishedIds)
	queryOpts = s.withExcludedFlows(queryOpts, append(opts.FinishedIds, opts.SkippedIds...))
	queryOpts = s.withUserElapsedTimeRule(queryOpts, opts.SignUpTimestamp)
	queryOpts = s.withUserSegment(queryOpts, opts.UserSegment)

	return queryOpts, nil
}

func (s *Enroller) withCurrentEnrollmentIdCondition(queryOpts bson.M, currentEnrollmentId string) bson.M {
	flowId, err := primitive.ObjectIDFromHex(currentEnrollmentId)
	if err != nil {
		return queryOpts
	}
	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), bson.M{"_id": flowId})

	return queryOpts
}

func (s *Enroller) withDependsOnCondition(queryOpts bson.M, ids []string) bson.M {
	inCondition := bson.M{"opts.dependsOn": bson.M{"$in": []string{}}}
	if len(ids) > 0 {
		inCondition = bson.M{"opts.dependsOn": bson.M{"$in": ids}}
	}
	emptyCondition := bson.M{"opts.dependsOn": bson.M{"$size": 0}}
	keyNotExist := bson.M{"opts.dependsOn": bson.M{"$exists": false}}

	clause := bson.M{"$or": []bson.M{inCondition, emptyCondition, keyNotExist}}
	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), clause)

	return queryOpts
}

func (s *Enroller) withExcludedFlows(queryOpts bson.M, excludedIds []string) bson.M {
	if len(excludedIds) == 0 {
		return queryOpts
	}

	primitiveIds := make([]primitive.ObjectID, len(excludedIds))
	for i, id := range excludedIds {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return queryOpts
		}
		primitiveIds[i] = oid
	}

	excludedCondition := bson.M{"$nin": primitiveIds}
	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), bson.M{"_id": excludedCondition})

	return queryOpts
}

func (s *Enroller) withUserElapsedTimeRule(queryOpts bson.M, signUpTimestamp int64) bson.M {
	keyNotExist := bson.M{"opts.targeting.rules": bson.M{"$exists": false}}
	emptyCondition := bson.M{"opts.targeting.rules": bson.M{"$size": 0}}

	ruleNotExist := bson.M{
		"opts.targeting.rules": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"condition": TargetingRuleUserElapsedDaysFromRegistration,
				},
			},
		},
	}

	clause := []bson.M{keyNotExist, emptyCondition, ruleNotExist}

	if signUpTimestamp != 0 {
		elapsedDays := calculateDaysBetweenTimestamps(signUpTimestamp, time.Now().Unix())
		conditionMatch := bson.M{
			"opts.targeting.rules": bson.M{
				"$elemMatch": bson.M{
					"condition": TargetingRuleUserElapsedDaysFromRegistration,
					"value":     bson.M{"$lte": elapsedDays},
				},
			},
		}
		clause = append(clause, conditionMatch)
	}

	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), bson.M{"$or": clause})

	return queryOpts
}

func (s *Enroller) withUserSegment(queryOpts bson.M, segment string) bson.M {
	keyNotExist := bson.M{"opts.targeting.rules": bson.M{"$exists": false}}
	emptyCondition := bson.M{"opts.targeting.rules": bson.M{"$size": 0}}
	ruleNotExist := bson.M{
		"opts.targeting.rules": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"condition": TargetingRuleUserSegment,
				},
			},
		},
	}

	clause := []bson.M{keyNotExist, emptyCondition, ruleNotExist}

	if segment != "" {
		conditionMatch := bson.M{
			"opts.targeting.rules": bson.M{
				"$elemMatch": bson.M{
					"condition": TargetingRuleUserSegment,
					"value":     segment,
				},
			},
		}
		clause = append(clause, conditionMatch)
	}

	queryOpts["$and"] = append(queryOpts["$and"].([]bson.M), bson.M{"$or": clause})

	return queryOpts
}

func calculateDaysBetweenTimestamps(providedTimestamp int64, currentTimestamp int64) int {
	providedTime := time.Unix(providedTimestamp, 0)
	currentTime := time.Unix(currentTimestamp, 0)
	duration := currentTime.Sub(providedTime)
	return int(duration.Hours() / 24)
}
