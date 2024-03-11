package flow

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BranchingService struct {
	Collection *mongo.Collection
}

func (s BranchingService) Get(workspace string, id string) (*Branching, error) {
	flowID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var branching Branching
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": flowID}).Decode(&branching)
	if err != nil {
		return nil, err
	}

	if branching.WorkspaceID != workspace {
		return nil, nil
	}

	return &branching, nil
}

func (s BranchingService) List(workspace string) ([]*Branching, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	branchings := make([]*Branching, 0)
	for cursor.Next(context.Background()) {
		var resFlow Branching
		err := cursor.Decode(&resFlow)
		if err != nil {
			return nil, err
		}

		if resFlow.WorkspaceID != workspace {
			continue
		}
		branchings = append(branchings, &resFlow)
	}

	return branchings, nil
}

func (s BranchingService) Create(workspace string, branching Branching) (string, error) {
	branching.WorkspaceID = workspace

	res, err := s.Collection.InsertOne(context.Background(), branching)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s BranchingService) Update(workspace string, id string, branching Branching) error {
	branchID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	branching.WorkspaceID = workspace
	_, err = s.Collection.ReplaceOne(context.Background(), bson.M{"_id": branchID}, branching)
	if err != nil {
		return err
	}

	return nil
}
