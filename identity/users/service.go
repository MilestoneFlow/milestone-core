package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) GetByCognitoId(cognitoId string) (*User, error) {
	var user User
	err := s.Collection.FindOne(context.Background(), bson.M{"cognitoId": cognitoId}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Service) Create(user User) error {
	_, err := s.Collection.InsertOne(context.Background(), user)
	return err
}

func (s Service) Update(cognitoId string, user User) error {
	_, err := s.Collection.UpdateOne(context.Background(), bson.M{"cognitoId": cognitoId}, bson.M{"$set": user})
	return err
}

func (s Service) AddWorkspaceEnrolled(cognitoId string, workspaceId string) error {
	_, err := s.Collection.UpdateOne(context.Background(), bson.M{"cognitoId": cognitoId}, bson.M{"$addToSet": bson.M{"workspacesEnrolled": workspaceId}})
	return err
}

func (s Service) RemoveWorkspaceEnrolled(cognitoId string, workspaceId string) error {
	_, err := s.Collection.UpdateOne(context.Background(), bson.M{"cognitoId": cognitoId}, bson.M{"$pull": bson.M{"workspacesEnrolled": workspaceId}})
	return err
}
