package progress

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"milestone_core/flow"
	"time"
)

type Service struct {
	Collection  *mongo.Collection
	FlowService flow.Service
}

func (s Service) GetProgress(flowId string, userId string) (*FlowProgress, error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"order", -1}})

	result := s.Collection.FindOne(context.Background(), bson.M{"flowId": flowId, "userId": userId}, findOptions)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}

	var progress *FlowProgress
	err := result.Decode(&progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s Service) MoveToNextStep(flowId string, userId string, completionDate int32) (*flow.Step, error) {
	currentProgress, err := s.GetProgress(flowId, userId)
	if err != nil {
		return nil, err
	}

	if nil == currentProgress {
		nextStep, err := s.FlowService.GetRootStep(flowId)
		if err != nil {
			return nil, err
		}

		_, err = s.Collection.InsertOne(context.Background(), FlowProgress{
			FlowId:      flowId,
			StepId:      nextStep.StepID,
			UserId:      userId,
			Status:      StatusStarted,
			Order:       1,
			EffectiveAt: primitive.Timestamp{T: uint32(completionDate)},
		})
		if err != nil {
			return nil, err
		}

		return nextStep, nil
	}

	if currentProgress.Status == StatusSkipped || currentProgress.Status == StatusCompleted {
		return nil, nil
	}

	flowProgressesToInsert := make([]interface{}, 0, 2)
	flowProgressesToInsert = append(flowProgressesToInsert, FlowProgress{
		FlowId:      flowId,
		StepId:      currentProgress.StepId,
		UserId:      currentProgress.UserId,
		Status:      StatusCompleted,
		Order:       currentProgress.Order + 1,
		EffectiveAt: primitive.Timestamp{T: uint32(completionDate)},
	})

	nextStep, err := s.FlowService.GetChildrenStep(flowId, currentProgress.StepId)
	if err != nil {
		return nil, err
	}
	if nextStep != nil {
		flowProgressesToInsert = append(flowProgressesToInsert, FlowProgress{
			FlowId:      flowId,
			StepId:      nextStep.StepID,
			UserId:      userId,
			Status:      StatusStarted,
			Order:       currentProgress.Order + 2,
			EffectiveAt: primitive.Timestamp{T: uint32(time.Now().Unix())},
		})
	}

	_, err = s.Collection.InsertMany(context.Background(), flowProgressesToInsert)
	if err != nil {
		return nil, err
	}

	return nextStep, nil
}

func (s Service) StartStep(flowId string, stepId string, userId string, timestamp uint32) (*flow.Step, error) {
	var nextStep *flow.Step
	var err error
	if "0" == stepId {
		nextStep, err = s.FlowService.GetRootStep(flowId)
	} else {
		var flowRes *flow.Flow
		flowRes, err = s.FlowService.Get(flowId)
		if flowRes != nil {
			for i := range flowRes.Steps {
				if flowRes.Steps[i].StepID == stepId {
					nextStep = &flowRes.Steps[i]
					break
				}
			}
		}
	}

	if err != nil {
		return nil, err
	}

	_, err = s.Collection.InsertOne(context.Background(), FlowProgress{
		FlowId:      flowId,
		StepId:      nextStep.StepID,
		UserId:      userId,
		Status:      StatusStarted,
		Order:       1,
		EffectiveAt: primitive.Timestamp{T: timestamp},
	})
	if err != nil {
		return nil, err
	}

	return nextStep, nil
}

func (s Service) CompleteStep(flowId string, stepId string, userId string, timestamp uint32) (*flow.Step, error) {
	nextStep, err := s.FlowService.GetChildrenStep(flowId, stepId)
	if err != nil {
		return nil, err
	}
	if nextStep == nil {
		return nil, nil
	}

	_, err = s.Collection.InsertOne(context.Background(), FlowProgress{
		FlowId:      flowId,
		StepId:      stepId,
		UserId:      userId,
		Status:      StatusCompleted,
		Order:       1,
		EffectiveAt: primitive.Timestamp{T: timestamp},
	})
	if err != nil {
		return nil, err
	}

	return nextStep, nil
}
