package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s Service) Create(user EnrolledUser) (interface{}, error) {
	user.Created = time.Now().Unix()
	result, err := s.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (s Service) Delete(id string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetLastUserState(workspace string, userId string) (*UserState, error) {
	var userState UserState

	opts := options.FindOne().SetSort(bson.D{{"created", -1}})
	err := s.UserStateCollection.FindOne(context.Background(), bson.M{"userId": userId, "workspaceId": workspace}, opts).Decode(&userState)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &userState, nil
}

func (s Service) GetFinishedFlowsForUser(workspace string, userId string) ([]string, error) {
	results, err := s.UserStateCollection.Distinct(context.Background(), "currentEnrolledFlowId", bson.M{"userId": userId, "workspaceId": workspace})
	if err != nil {
		return nil, err
	}

	flowsIds := make([]string, 0)
	for _, res := range results {
		flowsIds = append(flowsIds, res.(string))
	}

	return flowsIds, nil
}
