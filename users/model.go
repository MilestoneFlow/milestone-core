package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnrolledUser struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkspaceId     string             `json:"workspaceId" bson:"workspaceId"`
	Created         int64              `json:"created" bson:"created"`
	ExternalId      string             `json:"externalId" bson:"externalId"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	SignUpTimestamp int64              `json:"signUpTimestamp,omitempty" bson:"signUpTimestamp,omitempty"`
	Segment         string             `json:"segment,omitempty" bson:"segment,omitempty"`
}

type UserState struct {
	ID                    primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	WorkspaceId           string             `json:"-" bson:"workspaceId"`
	UserId                string             `json:"userId" bson:"userId"`
	CurrentEnrolledFlowId string             `json:"currentEnrolledFlowId" bson:"currentEnrolledFlowId"`
	CurrentStepId         string             `json:"currentStepId" bson:"currentStepId"`
	Times                 []UserStateTime    `json:"times" bson:"times"`
	Created               int64              `json:"-" bson:"created"`
}

type UserStateTime struct {
	StepId         string `json:"stepId" bson:"stepId"`
	StartTimestamp int64  `json:"startTimestamp" bson:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp" bson:"endTimestamp"`
}
