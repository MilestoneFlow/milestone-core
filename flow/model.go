package flow

import "go.mongodb.org/mongo-driver/bson/primitive"

type Step struct {
	StepID       string   `json:"stepId" bson:"stepId"`
	Data         StepData `json:"data" bson:"data"`
	Opts         StepOpts `json:"opts,omitempty" bson:"opts,omitempty"`
	ParentNodeId string   `json:"parentNodeId,omitempty" bson:"parentNodeId,omitempty"`
}

type StepData struct {
	Name               string      `json:"name" bson:"name"`
	Description        string      `json:"description" bson:"description"`
	TargetUrl          string      `json:"targetUrl,omitempty" bson:"targetUrl,omitempty"`
	AssignedCssElement string      `json:"assignedCssElement,omitempty" bson:"assignedCssElement,omitempty"`
	ElementType        string      `json:"elementType,omitempty" bson:"elementType,omitempty"`
	Placement          string      `json:"placement,omitempty" bson:"placement,omitempty"`
	ElementTemplate    string      `json:"elementTemplate,omitempty" bson:"elementTemplate,omitempty"`
	Blocks             []StepBlock `json:"blocks,omitempty" bson:"blocks,omitempty"`
}

type StepOpts struct {
	IsFinal    bool           `json:"isFinal,omitempty" bson:"isFinal,omitempty"`
	IsSource   bool           `json:"isSource,omitempty" bson:"isSource,omitempty"`
	SegmentID  string         `json:"segmentId,omitempty" bson:"segmentId,omitempty"`
	Transition StepTransition `json:"transition,omitempty" bson:"transition,omitempty"`
	Actionable bool           `json:"actionable,omitempty" bson:"actionable,omitempty"`
}

type StepTransition struct {
	InAnimation   string `json:"inAnimation,omitempty" bson:"inAnimation,omitempty"`
	OutAnimation  string `json:"outAnimation,omitempty" bson:"outAnimation,omitempty"`
	LoopAnimation string `json:"loopAnimation,omitempty" bson:"loopAnimation,omitempty"`
}

type Opts struct {
	Segmentation bool `json:"segmentation,omitempty" bson:"segmentation,omitempty"`
}

type Flow struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkspaceID string             `json:"workspaceId" bson:"workspaceId"`
	Name        string             `json:"name" bson:"name"`
	BaseURL     string             `json:"baseUrl,omitempty" bson:"baseUrl,omitempty"`
	Segments    []Segment          `json:"segments,omitempty" bson:"segments,omitempty"`
	Steps       []Step             `json:"steps" bson:"steps"`
	Relations   []Relation         `json:"relations" bson:"relations"`
	Opts        Opts               `json:"opts,omitempty" bson:"opts,omitempty"`
}

type Relation struct {
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type Segment struct {
	SegmentID string `json:"segmentId" bson:"segmentId"`
	Name      string `json:"name" bson:"name"`
	IconURL   string `json:"iconUrl" bson:"iconUrl"`
}

type StepBlock struct {
	Type StepBlockType `json:"type" bson:"type"`
	Data string        `json:"data" bson:"data"`
}

type StepBlockType string

const (
	StepBlockTypeText   StepBlockType = "text"
	StepBlockTypeImage  StepBlockType = "image"
	StepBlockTypeVideo  StepBlockType = "video"
	StepBlockTypeAvatar StepBlockType = "avatar"
)
