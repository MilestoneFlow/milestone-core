package flow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"milestone_core/users"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) Get(workspace string, id string) (*Flow, error) {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flow Flow
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": flowID, "workspaceId": workspace}).Decode(&flow)
	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (s Service) GetChildrenStep(workspace string, flowId string, parentStepId string, segmentId string) (*Step, error) {
	flow, err := s.Get(workspace, flowId)
	if err != nil {
		return nil, err
	}

	for i := range flow.Steps {
		if flow.Steps[i].ParentNodeId == parentStepId && (len(segmentId) == 0 || flow.Steps[i].Opts.SegmentID == segmentId) {
			return &flow.Steps[i], nil
		}
	}

	return nil, nil
}

func (s Service) GetStep(workspace string, flow *Flow, stepId string) *Step {
	for i := range flow.Steps {
		if flow.Steps[i].StepID == stepId {
			return &flow.Steps[i]
		}
	}

	return nil
}

func (s Service) GetRootStep(workspace string, flowId string) (*Step, error) {
	flow, err := s.Get(workspace, flowId)
	if err != nil {
		return nil, err
	}

	if flow != nil {
		for i := range flow.Steps {
			if len(flow.Steps[i].ParentNodeId) == 0 {
				return &flow.Steps[i], nil
			}
		}
	}

	return nil, errors.New("invalid flow id / no root available")
}

func (s Service) List(workspace string) ([]*Flow, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	flows := make([]*Flow, 0)
	for cursor.Next(context.Background()) {
		var resFlow Flow
		err := cursor.Decode(&resFlow)
		if err != nil {
			log.Fatal(err)
		}

		if resFlow.WorkspaceID == workspace {
			flows = append(flows, &resFlow)
		}
	}

	return flows, nil
}

func (s Service) Publish(workspace string, id string) error {
	flow, err := s.Get(workspace, id)
	if err != nil {
		return err
	}

	flow.Live = true
	err = s.saveUpdatedFlow(flow)

	return err
}

func (s Service) UnPublish(workspace string, id string) error {
	flow, err := s.Get(workspace, id)
	if err != nil {
		return err
	}

	flow.Live = false
	err = s.saveUpdatedFlow(flow)

	return err
}

func (s Service) ListLive(workspace string) ([]*Flow, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{"live": true})
	if err != nil {
		return nil, err
	}

	flows := make([]*Flow, 0)
	for cursor.Next(context.Background()) {
		var resFlow Flow
		err := cursor.Decode(&resFlow)
		if err != nil {
			log.Fatal(err)
		}

		if resFlow.WorkspaceID == workspace {
			flows = append(flows, &resFlow)
		}
	}

	return flows, nil
}

func (s Service) Update(workspace string, id string, updateInput UpdateInput) error {
	flow, err := s.Get(workspace, id)
	if err != nil {
		return err
	}

	if nil != updateInput.Name {
		flow.Name = *updateInput.Name
	}

	if nil != updateInput.BaseURL {
		flow.BaseURL = *updateInput.BaseURL
	}

	for _, deletedStep := range updateInput.DeletedSteps {
		for i := range flow.Steps {
			if flow.Steps[i].StepID == deletedStep {
				prevId := flow.Steps[i].ParentNodeId
				nextId := ""
				for _, step := range flow.Steps {
					if step.ParentNodeId == deletedStep {
						nextId = step.StepID
						break
					}
				}

				if len(nextId) > 0 {
					for i := range flow.Steps {
						if flow.Steps[i].StepID == nextId {
							flow.Steps[i].ParentNodeId = prevId
							break
						}
					}
				}

				flow.Steps = append(flow.Steps[:i], flow.Steps[i+1:]...)

				break
			}
		}
	}

	for _, updatedStep := range updateInput.UpdatedSteps {
		s.updateStepData(flow, &updatedStep)
	}

	s.updateStepsRelations(flow)

	for _, updatedSegment := range updateInput.Segments {
		if len(updatedSegment.Name) > 0 {
			found := false
			for i := range flow.Segments {
				if flow.Segments[i].SegmentID == updatedSegment.SegmentID {
					flow.Segments[i].Name = updatedSegment.Name
					found = true
					break
				}
			}

			if !found {
				flow.Segments = append(flow.Segments, Segment{
					SegmentID: uuid.New().String(),
					Name:      updatedSegment.Name,
					IconURL:   updatedSegment.IconURL,
				})
			}
		}
	}

	if updateInput.Trigger.TriggerID != "" {
		flow.Opts.Trigger = updateInput.Trigger
	}
	if updateInput.Targeting.TargetingID != "" {
		flow.Opts.Targeting = updateInput.Targeting
	}
	if updateInput.FinishEffect != nil {
		flow.Opts.FinishEffect = *updateInput.FinishEffect
	}
	if len(updateInput.Segments) == 0 {
		flow.Opts.Segmentation = false
	}

	if nil != updateInput.Opts {
		if updateInput.Opts.ThemeColor != "" {
			flow.Opts.ThemeColor = updateInput.Opts.ThemeColor
		}
		//if updateInput.Opts.AvatarId != "" {
		//	flow.Opts.AvatarId = updateInput.Opts.AvatarId
		//}
		flow.Opts.AvatarId = updateInput.Opts.AvatarId

		if updateInput.Opts.ElementTemplate != "" {
			flow.Opts.ElementTemplate = updateInput.Opts.ElementTemplate
		}
	}

	flowIndented, _ := json.MarshalIndent(flow, "", "\t")
	log.Default().Print(string(flowIndented))

	err = s.saveUpdatedFlow(flow)

	return err
}

func (s Service) UpdateStep(workspace string, flow *Flow, stepID string, updateInput Step) error {
	step := s.GetStep(workspace, flow, stepID)
	if step == nil {
		return errors.New("step not found")
	}

	step.Data = updateInput.Data
	if len(updateInput.ParentNodeId) > 0 {
		step.ParentNodeId = updateInput.ParentNodeId
	}
	err := s.saveUpdatedFlow(flow)

	return err
}

func (s Service) Capture(workspace string, id string, input UpdateInput) (string, error) {
	flow := Flow{
		Name:        *input.Name,
		BaseURL:     *input.BaseURL,
		WorkspaceID: workspace,
		Steps:       input.NewSteps,
		Opts: Opts{
			Segmentation:    false,
			Targeting:       Targeting{},
			Trigger:         Trigger{},
			ThemeColor:      "#000000",
			AvatarId:        "",
			ElementTemplate: "light",
			FinishEffect:    FinishEffect{},
		},
	}

	s.updateStepsRelations(&flow)

	newId, err := s.Collection.InsertOne(context.Background(), flow)
	if err != nil {
		return "", err
	}

	return newId.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s Service) GetFlowAnalytics(workspace string, flowId string) (FlowAnalytics, error) {
	flow, err := s.Get(workspace, flowId)
	if err != nil {
		return FlowAnalytics{
			FlowID:       "",
			Views:        0,
			AvgTotalTime: 0,
			AvgStepTime:  nil,
		}, err
	}
	if flow == nil {
		return FlowAnalytics{
			FlowID:       "",
			Views:        0,
			AvgTotalTime: 0,
			AvgStepTime:  nil,
		}, nil
	}

	analytics := FlowAnalytics{
		FlowID:       flow.ID.Hex(),
		Views:        0,
		AvgTotalTime: 0,
		AvgStepTime:  make(map[string]int64),
	}

	result, err := s.Collection.Database().Collection("users_state").Find(context.Background(), bson.M{"currentEnrolledFlowId": flow.ID.Hex()})
	if err != nil {
		return analytics, err
	}

	statesArr := make([]*users.UserState, 0)
	for result.Next(context.Background()) {
		var statesObj users.UserState
		err := result.Decode(&statesObj)
		if err != nil {
			log.Fatal(err)
		}

		statesArr = append(statesArr, &statesObj)
	}

	if len(statesArr) == 0 {
		return analytics, nil
	}

	analytics.Views = len(statesArr)
	noOfSteps := len(flow.Steps)
	noOfValidStates := 0
	totalFlowAvgTime := int64(0)
	stepAvgTime := make(map[string]int64)
	for _, state := range statesArr {
		currentFlowAvgTime := int64(0)
		if len(state.Times) == noOfSteps {
			noOfValidStates += 1
			for _, time := range state.Times {
				stepAvgTime[time.StepId] += time.EndTimestamp - time.StartTimestamp
				currentFlowAvgTime += time.EndTimestamp - time.StartTimestamp
			}
		}
		totalFlowAvgTime += currentFlowAvgTime
	}
	totalFlowAvgTime = 0
	if noOfValidStates > 0 {
		totalFlowAvgTime = totalFlowAvgTime / int64(noOfValidStates)
	}
	analytics.AvgTotalTime = totalFlowAvgTime
	for stepId, time := range stepAvgTime {
		analytics.AvgStepTime[stepId] = 0
		if noOfValidStates > 0 {
			analytics.AvgStepTime[stepId] = time / int64(noOfValidStates)
		}
	}

	return analytics, nil
}

func (s Service) updateStepData(flow *Flow, updatedStep *Step) {
	for i := range flow.Steps {
		if flow.Steps[i].StepID == updatedStep.StepID {
			flow.Steps[i].Data = updatedStep.Data
			if len(updatedStep.ParentNodeId) > 0 && flow.Steps[i].ParentNodeId != updatedStep.ParentNodeId {
				flow.Steps[i].ParentNodeId = updatedStep.ParentNodeId
			}

			return
		}
	}

	// Create new step
	s.createNewStep(flow, updatedStep)
}

func (s Service) createNewStep(flow *Flow, updatedStep *Step) {
	var rightNode *Step
	for i := range flow.Steps {
		if flow.Steps[i].ParentNodeId == updatedStep.ParentNodeId {
			rightNode = &flow.Steps[i]
			break
		}
	}

	if rightNode != nil {
		rightNode.ParentNodeId = updatedStep.StepID
	}
	flow.Steps = append(flow.Steps, Step{
		StepID:       updatedStep.StepID,
		ParentNodeId: updatedStep.ParentNodeId,
		Data:         updatedStep.Data,
		Opts: StepOpts{
			SegmentID: updatedStep.Opts.SegmentID,
		},
	})
}

func (s Service) saveUpdatedFlow(flow *Flow) error {
	_, err := s.Collection.UpdateByID(context.Background(), flow.ID, bson.M{"$set": flow})

	return err
}

func (s Service) updateStepsRelations(flow *Flow) {
	adjList := make(map[string][]*Step)
	var sourceStep *Step

	for i := range flow.Steps {
		if 0 == len(flow.Steps[i].ParentNodeId) {
			flow.Steps[i].Opts.IsSource = true
			sourceStep = &flow.Steps[i]
			continue
		} else {
			flow.Steps[i].Opts.IsSource = false
		}
		adjList[flow.Steps[i].ParentNodeId] = append(adjList[flow.Steps[i].ParentNodeId], &flow.Steps[i])
	}

	relations := make([]Relation, 0, len(flow.Steps)-1)

	q := make([]*Step, 0)
	q = append(q, sourceStep)
	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		if children, ok := adjList[node.StepID]; ok {
			node.Opts.IsFinal = false
			for _, child := range children {
				q = append(q, child)
				relations = append(relations, Relation{From: node.StepID, To: child.StepID})

				if len(node.Opts.SegmentID) > 0 && len(child.Opts.SegmentID) == 0 {
					child.Opts.SegmentID = node.Opts.SegmentID
				}
			}
		} else {
			node.Opts.IsFinal = true
		}
	}

	flow.Relations = relations
}
