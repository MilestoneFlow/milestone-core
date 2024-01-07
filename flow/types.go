package flow

type UpdateInput struct {
	Name         *string  `json:"name,omitempty"`
	UpdatedSteps []Step   `json:"updatedSteps,omitempty"`
	DeletedSteps []string `json:"deletedSteps,omitempty"`
	NewSteps     []Step   `json:"newSteps,omitempty"`
}
