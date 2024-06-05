package flows

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Service struct {
	Collection        *mongo.Collection
	ArchiveCollection *mongo.Collection
}

func (s Service) Get(workspace string, id string) (*Flow, error) {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flow Flow
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": flowID, "workspaceId": workspace}).Decode(&flow)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (s Service) Archive(workspace string, id string) error {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	flow, err := s.Get(workspace, id)
	if err != nil {
		return err
	}
	if flow == nil {
		return nil
	}

	_, err = s.ArchiveCollection.InsertOne(context.Background(), flow)
	if err != nil {
		return err
	}

	_, err = s.Collection.DeleteOne(context.Background(), bson.M{"_id": flowID})

	return err
}

func (s Service) ListArchivedFlows(workspace string) ([]*Flow, error) {
	cursor, err := s.ArchiveCollection.Find(context.Background(), bson.M{
		"workspaceId": workspace,
	})
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

		flows = append(flows, &resFlow)
	}

	return flows, nil
}

func (s Service) RestoreFlow(workspace string, id string) error {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var archivedFlow *Flow
	err = s.ArchiveCollection.FindOne(context.Background(), bson.M{"_id": flowID, "workspaceId": workspace}).Decode(&archivedFlow)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = s.Collection.InsertOne(context.Background(), archivedFlow)
	if err != nil {
		return err
	}

	_, err = s.ArchiveCollection.DeleteOne(context.Background(), bson.M{"_id": flowID, "workspaceId": workspace})
	return err
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

	if nil != updateInput.DependsOn {
		flow.Opts.DependsOn = updateInput.DependsOn
	}
	if updateInput.Trigger != nil {
		flow.Opts.Trigger = *updateInput.Trigger
	}
	if updateInput.Targeting != nil {
		flow.Opts.Targeting = *updateInput.Targeting
	}
	if updateInput.FinishEffect != nil {
		flow.Opts.FinishEffect = *updateInput.FinishEffect

		if flow.Opts.FinishEffect.Type == FinishEffectTypeFullScreenAnimation {
			baseUrL := "https://milestone-uploaded-flows-media.s3.amazonaws.com/assets/"
			var effectData FinishEffectDataFullScreenAnimation
			err = mapstructure.Decode(flow.Opts.FinishEffect.Data, &effectData)

			flow.Opts.FinishEffect.Data["position"] = FullScreenAnimationPositionMiddleScreen
			if effectData.Name == "fireworks_1" {
				flow.Opts.FinishEffect.Data["url"] = baseUrL + "fireworks_1.gif"
				flow.Opts.FinishEffect.Data["durationS"] = 4
				flow.Opts.FinishEffect.Data["position"] = FullScreenAnimationPositionBottomMiddle
			}
			if effectData.Name == "confetti_1" {
				flow.Opts.FinishEffect.Data["url"] = baseUrL + "confetti_1.gif"
				flow.Opts.FinishEffect.Data["durationS"] = 2
				flow.Opts.FinishEffect.Data["position"] = FullScreenAnimationPositionBottomMiddle
			}
			if effectData.Name == "congratulations_1" {
				flow.Opts.FinishEffect.Data["url"] = baseUrL + "congratulations_1.gif"
				flow.Opts.FinishEffect.Data["durationS"] = 5
			}
		}
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

	newId, err := s.Collection.InsertOne(context.Background(), flow)
	if err != nil {
		return "", err
	}

	return newId.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s Service) GetPossibleDependsOnListForFlow(workspace string, flowId string) ([]EssentialFlowInfo, error) {
	flows, err := s.List(workspace)
	if err != nil {
		return nil, err
	}

	isFlowInDependsOnList := func(flow Flow) bool {
		for _, dependsOnId := range flow.Opts.DependsOn {
			if dependsOnId == flowId {
				return true
			}
		}

		return false
	}

	dependsOnList := make([]EssentialFlowInfo, 0)
	for _, f := range flows {
		if f.ID.Hex() != flowId && !isFlowInDependsOnList(*f) {
			dependsOnList = append(dependsOnList, EssentialFlowInfo{
				ID:   f.ID.Hex(),
				Name: f.Name,
				Live: f.Live,
			})
		}
	}

	return dependsOnList, nil
}

type EssentialFlowInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Live bool   `json:"live"`
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
