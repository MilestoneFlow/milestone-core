package publicapi

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserData struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Workspace string             `json:"workspace" bson:"workspace"`
	UserID    string             `json:"userId" bson:"userId"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Created   int64              `json:"created,omitempty" bson:"created,omitempty"`
	Segment   string             `json:"segment,omitempty" bson:"segment,omitempty"`
}
