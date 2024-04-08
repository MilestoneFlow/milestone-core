package apiclient

import "go.mongodb.org/mongo-driver/bson/primitive"

type ApiClient struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkspaceID string             `json:"workspaceId" bson:"workspaceId"`
	Token       string             `json:"token" bson:"token"`
}
