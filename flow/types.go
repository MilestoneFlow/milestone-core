package flow

import "mime/multipart"

type UpdateInput struct {
	Name         *string                `json:"name,omitempty"`
	BaseURL      *string                `json:"baseUrl,omitempty"`
	Opts         *Opts                  `json:"opts,omitempty"`
	UpdatedSteps []Step                 `json:"updatedSteps,omitempty"`
	DeletedSteps []string               `json:"deletedSteps,omitempty"`
	NewSteps     []Step                 `json:"newSteps,omitempty"`
	Segments     []Segment              `json:"segments,omitempty"`
	Trigger      *Trigger               `json:"trigger,omitempty"`
	Targeting    *Targeting             `json:"targeting,omitempty"`
	FinishEffect *FinishEffect          `json:"finishEffect,omitempty"`
	MediaFiles   []multipart.FileHeader `json:"mediaFiles,omitempty"`
}
