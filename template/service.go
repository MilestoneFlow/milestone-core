package template

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/flow"
	"strings"
)

type Service struct {
	Collection     *mongo.Collection
	FlowCollection *mongo.Collection
}

func (s Service) List() ([]*flow.Flow, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	flows := make([]*flow.Flow, 0)
	for cursor.Next(context.Background()) {
		var resFlow flow.Flow
		err := cursor.Decode(&resFlow)
		if err != nil {
			return nil, err
		}
		flows = append(flows, &resFlow)
	}

	return flows, nil
}

func (s Service) Get(id string) (*flow.Flow, error) {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flowTemplate flow.Flow
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": flowID}).Decode(&flowTemplate)
	if err != nil {
		return nil, err
	}

	return &flowTemplate, nil
}

func (s Service) CreateFromTemplate(id string, override flow.Flow) (interface{}, error) {
	flowTemplate, err := s.Get(id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if len(override.Name) > 0 {
		flowTemplate.Name = override.Name
	}
	if len(override.BaseURL) > 0 {
		flowTemplate.BaseURL = override.BaseURL
	}
	if len(override.Segments) > 0 {
		flowTemplate.Segments = make([]flow.Segment, 0)
		for _, segment := range override.Segments {
			flowTemplate.Segments = append(flowTemplate.Segments, flow.Segment{
				SegmentID: strings.ReplaceAll(strings.ToLower(segment.Name), " ", "_"),
				Name:      segment.Name,
			})
		}

		flowTemplate.Steps = flowTemplate.Steps[:1]
		for _, segment := range flowTemplate.Segments {
			flowTemplate.Steps = append(flowTemplate.Steps, flow.Step{
				ParentNodeId: flowTemplate.Steps[0].StepID,
				StepID:       segment.SegmentID,
				Data: flow.StepData{
					Name: segment.Name + " Step 1",
				},
				Opts: flow.StepOpts{
					SegmentID: segment.SegmentID,
					IsFinal:   true,
				},
			})
		}

		s.updateStepsRelations(flowTemplate)
	}

	flowTemplate.ID = primitive.NilObjectID
	result, err := s.FlowCollection.InsertOne(context.Background(), flowTemplate)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID, nil
}

func (s Service) updateStepsRelations(template *flow.Flow) {
	adjList := make(map[string][]*flow.Step)
	var sourceStep *flow.Step

	for i := range template.Steps {
		if 0 == len(template.Steps[i].ParentNodeId) {
			template.Steps[i].Opts.IsSource = true
			sourceStep = &template.Steps[i]
			continue
		}
		adjList[template.Steps[i].ParentNodeId] = append(adjList[template.Steps[i].ParentNodeId], &template.Steps[i])
	}

	relations := make([]flow.Relation, 0, len(template.Steps)-1)

	q := make([]*flow.Step, 0)
	q = append(q, sourceStep)
	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		if children, ok := adjList[node.StepID]; ok {
			node.Opts.IsFinal = false
			for _, child := range children {
				q = append(q, child)
				relations = append(relations, flow.Relation{From: node.StepID, To: child.StepID})

				if len(node.Opts.SegmentID) > 0 {
					child.Opts.SegmentID = node.Opts.SegmentID
				}
			}
		} else {
			node.Opts.IsFinal = true
		}
	}

	template.Relations = relations
}
