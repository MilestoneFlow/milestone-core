package flow

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Enroller struct {
	Collection *mongo.Collection
}

type EnrollmentOpts struct {
	SkippedIds          []string
	FinishedIds         []string
	CurrentEnrollmentId string
}

func (s *Enroller) GetFlow(workspaceId string, opts EnrollmentOpts) (*Flow, error) {
	queryOpts, err := s.buildQueryOpts(workspaceId, opts)
	if err != nil {
		return nil, err
	}

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
	queryOpts := bson.M{"workspaceId": workspaceId, "live": true}

	if opts.CurrentEnrollmentId != "" {
		queryOpts = s.withCurrentEnrollmentIdCondition(queryOpts, opts.CurrentEnrollmentId)
		return queryOpts, nil
	}

	queryOpts = s.withDependsOnCondition(queryOpts, opts.FinishedIds)
	s.withExcludedFlows(queryOpts, append(opts.FinishedIds, opts.SkippedIds...))

	return queryOpts, nil
}

func (s *Enroller) withCurrentEnrollmentIdCondition(queryOpts bson.M, currentEnrollmentId string) bson.M {
	flowId, err := primitive.ObjectIDFromHex(currentEnrollmentId)
	if err != nil {
		return queryOpts
	}
	queryOpts["_id"] = flowId

	return queryOpts
}

func (s *Enroller) withDependsOnCondition(queryOpts bson.M, ids []string) bson.M {
	inCondition := bson.M{"opts.dependsOn": bson.M{"$in": []string{}}}
	if len(ids) > 0 {
		inCondition = bson.M{"opts.dependsOn": bson.M{"$in": ids}}
	}
	emptyCondition := bson.M{"opts.dependsOn": bson.M{"$size": 0}}
	keyNotExist := bson.M{"opts.dependsOn": bson.M{"$exists": false}}

	queryOpts["$or"] = []bson.M{inCondition, emptyCondition, keyNotExist}

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
	queryOpts["_id"] = excludedCondition

	return queryOpts
}

func (s *Enroller) withUserElapsedTimeRule(queryOpts bson.M, currentTimestamp int64) bson.M {
	keyNotExist := bson.M{"opts.targeting.rules": bson.M{"$exists": false}}
	keyNotExist := bson.M{"opts.targeting.rules": bson.M{"$exists": false}}

	queryOpts["$or"] = []bson.M{keyNotExist, emptyCondition, keyNotExist}

	return queryOpts
}
