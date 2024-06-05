package wallets

import "go.mongodb.org/mongo-driver/bson/primitive"

type Point struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EnrolledUserID string             `json:"enrolledUserId" bson:"enrolledUserId"`
	Amount         int                `json:"amount" bson:"amount"`
	EffectiveAt    int64              `json:"effectiveAt" bson:"effectiveAt"`
	Created        int64              `json:"created" bson:"created"`
}
