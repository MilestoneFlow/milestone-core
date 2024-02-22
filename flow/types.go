package flow

type UpdateInput struct {
	Name         *string   `json:"name,omitempty"`
	BaseURL      *string   `json:"baseUrl,omitempty"`
	UpdatedSteps []Step    `json:"updatedSteps,omitempty"`
	DeletedSteps []string  `json:"deletedSteps,omitempty"`
	NewSteps     []Step    `json:"newSteps,omitempty"`
	Segments     []Segment `json:"segments,omitempty"`
}
