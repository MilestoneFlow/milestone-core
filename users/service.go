package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) List() ([]*EnrolledUser, error) {
	cursor, err := s.Collection.Find(context.Background(), bson.M{})
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

func (s Service) Get(id string) (*EnrolledUser, error) {
	var user EnrolledUser
	err := s.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Service) Create(user EnrolledUser) (interface{}, error) {
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
