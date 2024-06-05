package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                 primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	CognitoID          string             `json:"cognitoId" bson:"cognitoId"`
	WorkspacesEnrolled []string           `json:"workspacesEnrolled" bson:"workspacesEnrolled"`
	Created            int64              `json:"created" bson:"created"`
}
