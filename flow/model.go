package flow

import "go.mongodb.org/mongo-driver/bson/primitive"

type Step struct {
	StepID       string   `json:"stepId" bson:"stepId"`
	Data         StepData `json:"data" bson:"data"`
	Opts         StepOpts `json:"opts,omitempty" bson:"opts,omitempty"`
	ParentNodeId string   `json:"parentNodeId,omitempty" bson:"parentNodeId,omitempty"`
}

type StepData struct {
	Name               string `json:"name" bson:"name"`
	Description        string `json:"description" bson:"description"`
	TargetUrl          string `json:"targetUrl,omitempty" bson:"targetUrl,omitempty"`
	AssignedCssElement string `json:"assignedCssElement,omitempty" bson:"assignedCssElement,omitempty"`
	ElementType        string `json:"elementType,omitempty" bson:"elementType,omitempty"`
}

type StepOpts struct {
	IsFinal   bool   `json:"isFinal,omitempty" bson:"isFinal,omitempty"`
	IsSource  bool   `json:"isSource,omitempty" bson:"isSource,omitempty"`
	SegmentID string `json:"segmentId,omitempty" bson:"segmentId,omitempty"`
}

type Opts struct {
	Segmentation bool `json:"segmentation,omitempty" bson:"segmentation,omitempty"`
}

type Flow struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	BaseURL   string             `json:"baseUrl,omitempty" bson:"baseUrl,omitempty"`
	Segments  []Segment          `json:"segments,omitempty" bson:"segments,omitempty"`
	Steps     []Step             `json:"steps" bson:"steps"`
	Relations []Relation         `json:"relations" bson:"relations"`
	Opts      Opts               `json:"opts,omitempty" bson:"opts,omitempty"`
}

type Relation struct {
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type Segment struct {
	SegmentID string `json:"segmentId" bson:"segmentId"`
	Name      string `json:"name" bson:"name"`
}
