package workspace

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) Get(workspaceId string) (*Workspace, error) {
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	var workspace Workspace
	err = s.Collection.FindOne(context.Background(), bson.M{"_id": primitiveId}).Decode(&workspace)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (s Service) GetByUserIdentifier(userIdentifier string) (*Workspace, error) {
	var workspace Workspace
	err := s.Collection.FindOne(context.Background(), bson.M{"userIdentifiers": bson.M{"$in": []string{userIdentifier}}}).Decode(&workspace)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (s Service) Create(workspace Workspace) error {
	_, err := s.Collection.InsertOne(context.Background(), workspace)
	return err
}

func (s Service) Update(workspace Workspace) error {
	_, err := s.Collection.ReplaceOne(context.Background(), bson.M{"_id": workspace.ID}, workspace)
	return err
}

func (s Service) AddUserIdentifier(workspaceId string, userIdentifier string) error {
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"_id": primitiveId}, bson.M{"$push": bson.M{"userIdentifiers": userIdentifier}})
	return err
}

func (s Service) RemoveUserIdentifier(workspaceId string, userIdentifier string) error {
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"_id": primitiveId}, bson.M{"$pull": bson.M{"userIdentifiers": userIdentifier}})
	return err
}
