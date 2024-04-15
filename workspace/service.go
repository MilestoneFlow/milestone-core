package workspace

import (
	"context"
	"crypto/rand"
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

func (s Service) CreateForUser(workspace Workspace, userIdentifier string) error {
	inviteToken, err := s.generateInviteToken()
	if err != nil {
		return err
	}

	workspace.InviteToken = inviteToken
	workspace.UserIdentifiers = []string{userIdentifier}
	_, err = s.Collection.InsertOne(context.Background(), workspace)

	return err
}

func (s Service) Update(id string, workspace Workspace) error {
	workspace.ID, _ = primitive.ObjectIDFromHex(id)
	workspace.UserIdentifiers = nil
	workspace.InviteToken = ""
	_, err := s.Collection.UpdateOne(context.Background(), bson.M{"_id": workspace.ID}, bson.M{"$set": workspace})
	return err
}

func (s Service) AddUserIdentifiers(workspaceId string, userIdentifiers []string) error {
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	for _, userIdentifier := range userIdentifiers {
		_, err = s.Collection.UpdateOne(context.Background(), bson.M{"_id": primitiveId}, bson.M{"$addToSet": bson.M{"userIdentifiers": userIdentifier}})
	}
	return err
}

func (s Service) RemoveUserIdentifier(workspaceId string, userIdentifier string) error {
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"_id": primitiveId}, bson.M{"$pull": bson.M{"userIdentifiers": userIdentifier}})
	return err
}

func (s Service) RefreshInviteToken(workspaceId string) (string, error) {
	inviteURL, err := s.generateInviteToken()
	if err != nil {
		return "", err
	}
	primitiveId, err := primitive.ObjectIDFromHex(workspaceId)

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"_id": primitiveId}, bson.M{"$set": bson.M{"inviteToken": inviteURL}})
	return inviteURL, err
}

func (s Service) generateInviteToken() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
