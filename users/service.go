package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	Collection          *mongo.Collection
	UserStateCollection *mongo.Collection
}

func (s Service) List(workspaceId string) ([]*EnrolledUser, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{"workspaceId": workspaceId})
	if err != nil {
		return nil, err
	}

	users := make([]*EnrolledUser, 0)
	for cursor.Next(context.Background()) {
		var resUser EnrolledUser
		err := cursor.Decode(&resUser)
		if err != nil {
			return nil, err
		}
		users = append(users, &resUser)
	}

	return users, nil
}

func (s Service) Get(workspace string, externalId string) (*EnrolledUser, error) {
	var user EnrolledUser
	err := s.Collection.FindOne(context.Background(), bson.M{"externalId": externalId, "workspaceId": workspace}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s Service) Create(user EnrolledUser) error {
	user.Created = time.Now().Unix()
	result, err := s.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	enrolledUserId := result.InsertedID.(primitive.ObjectID).Hex()
	_, err = s.UserStateCollection.InsertOne(context.Background(), UserState{
		UserID:           enrolledUserId,
		WorkspaceID:      user.WorkspaceId,
		FlowsData:        FlowsData{},
		Metadata:         map[string]string{},
		UpdatedTimestamp: time.Now().Unix(),
	})

	return err
}

func (s Service) Delete(id string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetState(workspace string, userId string) (*UserState, error) {
	var userState UserState

	err := s.UserStateCollection.FindOne(context.Background(), bson.M{"userId": userId, "workspaceId": workspace}).Decode(&userState)
	if err != nil {
		return nil, err
	}

	return &userState, nil
}

func (s Service) PutState(workspace string, userId string, state UserState) error {
	currentState, err := s.GetState(workspace, userId)
	if err != nil {
		return err
	}

	state.UpdatedTimestamp = time.Now().Unix()
	state.WorkspaceID = currentState.WorkspaceID
	state.UserID = currentState.UserID
	_, err = s.UserStateCollection.UpdateOne(context.Background(), bson.M{"_id": currentState.ID}, bson.M{"$set": state})

	return err
}
