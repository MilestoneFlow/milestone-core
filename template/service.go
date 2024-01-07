package template

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/flow"
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

	flowTemplate.ID = primitive.NilObjectID
	result, err := s.FlowCollection.InsertOne(context.Background(), flowTemplate)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID, nil
}
