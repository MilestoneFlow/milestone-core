package publicapi

import "milestone_core/tracker"

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
