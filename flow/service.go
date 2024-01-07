package flow

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) Get(id string) (*Flow, error) {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flow Flow
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": flowID}).Decode(&flow)
	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (s Service) GetChildrenStep(flowId string, parentStepId string) (*Step, error) {
	flow, err := s.Get(flowId)
	if err != nil {
		return nil, err
	}

	for i := range flow.Steps {
		if flow.Steps[i].ParentNodeId == parentStepId {
			return &flow.Steps[i], nil
		}
	}

	return nil, nil
}

func (s Service) GetRootStep(flowId string) (*Step, error) {
	flow, err := s.Get(flowId)
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

func (s Service) List() ([]*Flow, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var flows []*Flow
	for cursor.Next(context.Background()) {
		var resFlow Flow
		err := cursor.Decode(&resFlow)
		if err != nil {
			log.Fatal(err)
		}
		flows = append(flows, &resFlow)
	}

	return flows, nil
}

func (s Service) Update(id string, updateInput UpdateInput) error {
	flow, err := s.Get(id)
	if err != nil {
		return err
	}

	if nil != updateInput.Name {
		flow.Name = *updateInput.Name
	}

	for _, updatedStep := range updateInput.UpdatedSteps {
		s.updateStepData(flow, &updatedStep)
	}

	s.updateStepsRelations(flow)

	flowIndented, _ := json.MarshalIndent(flow, "", "\t")
	log.Default().Print(string(flowIndented))

	err = s.saveUpdatedFlow(flow)

	return err
}

func (s Service) UpdateStepData(id string, idStep string, updateStepDataInput StepData) error {
	flow, err := s.Get(id)
	if err != nil {
		return err
	}

	for i := range flow.Steps {
		if flow.Steps[i].StepID == idStep {
			if len(updateStepDataInput.AssignedCssElement) > 0 {
				flow.Steps[i].Data.AssignedCssElement = updateStepDataInput.AssignedCssElement
			}
			break
		}
	}

	err = s.saveUpdatedFlow(flow)

	return err
}

func (s Service) updateStepData(flow *Flow, updatedStep *Step) {
	for i := range flow.Steps {
		if flow.Steps[i].StepID == updatedStep.StepID {
			flow.Steps[i].Data = updatedStep.Data
			return
		}
	}
	flow.Steps = append(flow.Steps, Step{
		StepID:       updatedStep.StepID,
		ParentNodeId: updatedStep.ParentNodeId,
		Data: StepData{
			Name:               updatedStep.Data.Name,
			Description:        updatedStep.Data.Description,
			TargetUrl:          updatedStep.Data.TargetUrl,
			AssignedCssElement: updatedStep.Data.AssignedCssElement,
			ElementType:        updatedStep.Data.ElementType,
		},
	})
}

func (s Service) saveUpdatedFlow(flow *Flow) error {
	_, err := s.Collection.UpdateByID(context.Background(), flow.ID, bson.M{"$set": flow})

	return err
}

func (s Service) updateStepsRelations(flow *Flow) {
	adjList := make(map[string][]string)

	for _, step := range flow.Steps {
		if 0 == len(step.ParentNodeId) {
			continue
		}
		adjList[step.ParentNodeId] = append(adjList[step.ParentNodeId], step.StepID)
	}

	noOfRelations := len(flow.Steps) - 1
	relations := make([]Relation, 0, noOfRelations)
	for parentNodeId, childrenNodes := range adjList {
		for _, childrenNodeId := range childrenNodes {
			relations = append(relations, Relation{From: parentNodeId, To: childrenNodeId})
		}
	}

	for i := range flow.Steps {
		if _, ok := adjList[flow.Steps[i].StepID]; ok {
			flow.Steps[i].Opts.IsFinal = false
		} else {
			flow.Steps[i].Opts.IsFinal = true
		}
	}

	flow.Relations = relations
}
