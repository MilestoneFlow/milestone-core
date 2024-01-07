package progress

import "go.mongodb.org/mongo-driver/bson/primitive"

type Status string

const (
	StatusStarted   Status = "STARTED"
	StatusCompleted Status = "COMPLETED"
	StatusSkipped   Status = "SKIPPED"
)

type FlowProgress struct {
	FlowId      string              `json:"flowId" bson:"flowId"`
	StepId      string              `json:"stepId" bson:"stepId"`
	UserId      string              `json:"userId" bson:"userId"`
	Status      Status              `json:"status" bson:"status"`
	Order       uint32              `json:"order" bson:"order"`
	EffectiveAt primitive.Timestamp `json:"effectiveAt" bson:"effectiveAt"`
}
