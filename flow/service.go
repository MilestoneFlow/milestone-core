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

func (s Service) GetChildrenStep(flowId string, parentStepId string, segmentId string) (*Step, error) {
	flow, err := s.Get(flowId)
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

func (s Service) GetStep(flowId string, stepId string) (*Step, error) {
	flow, err := s.Get(flowId)
	if err != nil {
		return nil, err
	}

	for i := range flow.Steps {
		if flow.Steps[i].StepID == stepId {
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
			for i := range flow.Segments {
				if flow.Segments[i].SegmentID == updatedSegment.SegmentID {
					flow.Segments[i].Name = updatedSegment.Name
					break
				}
			}
		}
	}

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

func (s Service) Capture(id string, newSteps []Step) error {
	flow, err := s.Get(id)
	if err != nil {
		return err
	}

	flow.Steps = newSteps
	flow.Segments = make([]Segment, 0)
	s.updateStepsRelations(flow)

	err = s.saveUpdatedFlow(flow)

	return err
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

	flow.Steps = append(flow.Steps, Step{
		StepID:       updatedStep.StepID,
		ParentNodeId: updatedStep.ParentNodeId,
		Data: StepData{
			Name:               updatedStep.Data.Name,
			Description:        updatedStep.Data.Description,
			TargetUrl:          updatedStep.Data.TargetUrl,
			AssignedCssElement: updatedStep.Data.AssignedCssElement,
			ElementType:        updatedStep.Data.ElementType,
			Placement:          updatedStep.Data.Placement,
			ElementTemplate:    updatedStep.Data.ElementTemplate,
		},
		Opts: StepOpts{
			SegmentID:  updatedStep.Opts.SegmentID,
			Transition: updatedStep.Opts.Transition,
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
