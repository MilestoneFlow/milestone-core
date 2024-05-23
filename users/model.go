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
	ID               primitive.ObjectID `json:"-,omitempty" bson:"_id,omitempty"`
	WorkspaceID      string             `json:"workspaceId" bson:"workspaceId"`
	UserID           string             `json:"userId" bson:"userId"`
	FlowsData        FlowsData          `json:"flowsData" bson:"flowsData"`
	Metadata         map[string]string  `json:"metadata" bson:"metadata"`
	UpdatedTimestamp int64              `json:"updatedTimestamp" bson:"updatedTimestamp"`
}

type FlowsData struct {
	CompletedFlowsIds []string `json:"completedFlowsIds" bson:"completedFlowsIds"`
	SkippedFlowsIds   []string `json:"skippedFlowsIds" bson:"skippedFlowsIds"`
	CurrentFlowID     string   `json:"currentFlowId" bson:"currentFlowId"`
}
