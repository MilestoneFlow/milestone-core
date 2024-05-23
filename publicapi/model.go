package publicapi

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventTrack struct {
	ID             primitive.ObjectID `json:"-,omitempty" bson:"_id,omitempty"`
	WorkspaceID    string             `json:"workspaceId" bson:"workspaceId"`
	ExternalUserID string             `json:"externalUserId" bson:"externalUserId"`
	EntityID       string             `json:"entityId" bson:"entityId"`
	EventType      EventType          `json:"eventType" bson:"eventType"`
	Timestamp      int64              `json:"timestamp" bson:"timestamp"`
	Metadata       map[string]string  `json:"metadata" bson:"metadata"`
}

type EventType string

const (
	EventTypeHelperClick    EventType = "helper_click"
	EventTypeHelperHover    EventType = "helper_hover"
	EventTypeHelperClose    EventType = "helper_close"
	EventTypeFlowStepStart  EventType = "flow_step_start"
	EventTypeFlowStepFinish EventType = "flow_step_finish"
	EventTypeFlowSkipped    EventType = "flow_skipped"
)
