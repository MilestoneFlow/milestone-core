package publicapi

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"milestone_core/apiclient"
	"milestone_core/users"
	"time"
)

type UserStateService struct {
	UserStateCollection *mongo.Collection
	ApiClientService    apiclient.Service
	EnrolledUserService users.Service
}

func (s UserStateService) GetState(token string, externalUserId string) (*users.UserState, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	enrolledUser, err := s.EnrolledUserService.Get(apiClient.WorkspaceID, externalUserId)
	if err != nil {
		return nil, err
	}
	if enrolledUser == nil {
		return nil, errors.New("user not found")
	}

	var userState users.UserState
	opts := options.FindOne().SetSort(bson.D{{"created", -1}})
	err = s.UserStateCollection.FindOne(context.Background(), bson.M{"userId": enrolledUser.ID}, opts).Decode(&userState)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &userState, nil
}

func (s UserStateService) PutState(token string, externalUserId string, newState users.UserState) error {
	currentState, err := s.GetState(token, externalUserId)
	if err != nil {
		return err
	}

	if currentState == nil {
		_, err := s.createState(token, externalUserId, newState)
		if err != nil {
			return err
		}
		return nil
	}

	newState.WorkspaceID = currentState.WorkspaceID
	newState.UserID = currentState.UserID
	newState.UpdatedTimestamp = time.Now().Unix()
	res, err := s.UserStateCollection.UpdateByID(context.Background(), currentState.ID, newState)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user state not found")
	}

	return nil
}

func (s UserStateService) createState(token string, externalUserId string, newState users.UserState) (interface{}, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	enrolledUser, err := s.EnrolledUserService.Get(apiClient.WorkspaceID, externalUserId)
	if err != nil {
		return nil, err
	}
	if enrolledUser == nil {
		return nil, errors.New("user not found")
	}

	newState.WorkspaceID = apiClient.WorkspaceID
	newState.UserID = enrolledUser.ID.Hex()
	newState.UpdatedTimestamp = time.Now().Unix()
	result, err := s.UserStateCollection.InsertOne(context.Background(), newState)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
