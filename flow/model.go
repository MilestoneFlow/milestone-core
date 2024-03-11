package flow

import "go.mongodb.org/mongo-driver/bson/primitive"

type Step struct {
	StepID       string   `json:"stepId" bson:"stepId"`
	Data         StepData `json:"data" bson:"data"`
	Opts         StepOpts `json:"opts,omitempty" bson:"opts,omitempty"`
	ParentNodeId string   `json:"parentNodeId,omitempty" bson:"parentNodeId,omitempty"`
}

type StepData struct {
	TargetUrl          string          `json:"targetUrl" bson:"targetUrl,omitempty"`
	AssignedCssElement string          `json:"assignedCssElement" bson:"assignedCssElement,omitempty"`
	ElementType        StepElementType `json:"elementType" bson:"elementType,omitempty"`
	Placement          StepPlacement   `json:"placement" bson:"placement,omitempty"`
	Blocks             []StepBlock     `json:"blocks" bson:"blocks,omitempty"`
	Transition         StepTransition  `json:"transition" bson:"transition,omitempty"`
	ActionType         StepActionType  `json:"actionType" bson:"actionType,omitempty"`
}

type StepOpts struct {
	IsFinal   bool   `json:"isFinal" bson:"isFinal,omitempty"`
	IsSource  bool   `json:"isSource" bson:"isSource,omitempty"`
	SegmentID string `json:"segmentId" bson:"segmentId,omitempty"`
}

type StepTransition struct {
	InAnimation   string `json:"inAnimation" bson:"inAnimation,omitempty"`
	OutAnimation  string `json:"outAnimation" bson:"outAnimation,omitempty"`
	LoopAnimation string `json:"loopAnimation" bson:"loopAnimation,omitempty"`
}

type Opts struct {
	Segmentation    bool                `json:"segmentation,omitempty" bson:"segmentation,omitempty"`
	Targeting       Targeting           `json:"targeting,omitempty" bson:"targeting,omitempty"`
	Trigger         Trigger             `json:"trigger,omitempty" bson:"trigger,omitempty"`
	ThemeColor      string              `json:"themeColor,omitempty" bson:"themeColor,omitempty"`
	AvatarId        string              `json:"avatarId,omitempty" bson:"avatarId,omitempty"`
	ElementTemplate StepElementTemplate `json:"elementTemplate" bson:"elementTemplate,omitempty"`
	FinishEffect    FinishEffect        `json:"finishEffect,omitempty" bson:"finishEffect,omitempty"`
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
	BlockID string        `json:"blockId" bson:"blockId"`
	Type    StepBlockType `json:"type" bson:"type"`
	Data    string        `json:"data" bson:"data"`
	Order   int           `json:"order" bson:"order"`
}

type Trigger struct {
	TriggerID string        `json:"triggerId,omitempty" bson:"triggerId,omitempty"`
	Rules     []TriggerRule `json:"rules,omitempty" bson:"rules,omitempty"`
}

type TriggerRule struct {
	Condition string `json:"condition,omitempty" bson:"condition,omitempty"`
	Value     string `json:"value,omitempty" bson:"value,omitempty"`
}

type Targeting struct {
	TargetingID string          `json:"targetingId,omitempty" bson:"targetingId,omitempty"`
	Rules       []TargetingRule `json:"rules,omitempty" bson:"rules,omitempty"`
}

type TargetingRule struct {
	Condition string `json:"condition,omitempty" bson:"condition,omitempty"`
	Value     string `json:"value,omitempty" bson:"value,omitempty"`
}

type FinishEffect struct {
	Type FinishEffectType `json:"type,omitempty" bson:"type,omitempty"`
	Data string           `json:"data,omitempty" bson:"data,omitempty"`
}

type FinishEffectType string

const (
	FinishEffectTypePopup               FinishEffectType = "popup"
	FinishEffectTypeFullScreenAnimation FinishEffectType = "full_screen_animation"
)

type StepBlockType string

const (
	StepBlockTypeText   StepBlockType = "text"
	StepBlockTypeImage  StepBlockType = "image"
	StepBlockTypeVideo  StepBlockType = "video"
	StepBlockTypeAvatar StepBlockType = "avatar"
)

type StepElementType string

const (
	StepElementTypeTooltip   StepElementType = "tooltip"
	StepElementTypePopup     StepElementType = "popup"
	StepElementTypeBranching StepElementType = "branching"
)

type StepPlacement string

const (
	StepPlacementBottom StepPlacement = "bottom"
	StepPlacementTop    StepPlacement = "top"
	StepPlacementLeft   StepPlacement = "left"
	StepPlacementRight  StepPlacement = "right"
)

type StepElementTemplate string

const (
	StepElementTemplateLight StepElementTemplate = "light"
	StepElementTemplateDark  StepElementTemplate = "dark"
)

type StepActionType string

const (
	StepActionTypeNpAction StepActionType = "no_action"
	StepActionTypeAction   StepActionType = "action"
)

type Branching struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkspaceID string             `json:"workspaceId" bson:"workspaceId"`
	Name        string             `json:"name" bson:"name"`
	Content     string             `json:"content" bson:"content"`
	BaseURL     string             `json:"baseUrl,omitempty" bson:"baseUrl,omitempty"`
	Variants    []BranchingVariant `json:"variants" bson:"variants"`
	TargetURL   string             `json:"targetUrl,omitempty" bson:"targetUrl,omitempty"`
}

type BranchingVariant struct {
	VariantID string `json:"variantId" bson:"variantId"`
	FlowID    string `json:"flowId" bson:"flowId"`
	Name      string `json:"name" bson:"name"`
}
