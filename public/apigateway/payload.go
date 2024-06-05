package apigateway

import (
	"milestone_core/tours/tracker"
)

type TrackEventsRequest struct {
	Events         []tracker.EventTrack `json:"data"`
	ExternalUserID string               `json:"externalUserId"`
}

type FlowStateUpdateRequest struct {
	FlowID        string `json:"flowId"`
	CurrentStepID string `json:"currentStepId"`
	Finished      bool   `json:"finished"`
	Skipped       bool   `json:"skipped"`
}
