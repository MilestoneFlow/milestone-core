package workspace

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workspace struct {
	ID              primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	BaseURL         string             `json:"baseUrl" bson:"baseUrl,omitempty"`
	UserIdentifiers []string           `json:"userIdentifiers" bson:"userIdentifiers,omitempty"`
	InviteToken     string             `json:"inviteToken" bson:"inviteToken,omitempty"`
}
