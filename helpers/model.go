package helpers

import "go.mongodb.org/mongo-driver/bson/primitive"

type Helper struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	PublicID     string             `json:"publicId" bson:"publicId"`
	WorkspaceID  string             `json:"-" bson:"workspaceId"`
	Published    bool               `json:"published" bson:"published"`
	Name         string             `json:"name" bson:"name"`
	Data         HelperData         `json:"data" bson:"data"`
	RenderAction HelperRenderAction `json:"renderAction" bson:"renderAction"`
	Created      int64              `json:"created" bson:"created"`
	Updated      int64              `json:"updated" bson:"updated"`
	PublishedAt  int64              `json:"publishedAt" bson:"publishedAt"`
}

type HelperData struct {
	TargetUrl          string            `json:"targetUrl" bson:"targetUrl,omitempty"`
	AssignedCssElement string            `json:"assignedCssElement" bson:"assignedCssElement,omitempty"`
	ElementType        HelperElementType `json:"elementType" bson:"elementType,omitempty"`
	Placement          HelperPlacement   `json:"placement" bson:"placement,omitempty"`
	Blocks             []HelperBlock     `json:"blocks" bson:"blocks,omitempty"`
}

type HelperBlockType string

const (
	HelperBlockTypeText   HelperBlockType = "text"
	HelperBlockTypeImage  HelperBlockType = "image"
	HelperBlockTypeVideo  HelperBlockType = "video"
	HelperBlockTypeAvatar HelperBlockType = "avatar"
)

type HelperPlacement string

const (
	HelperPlacementBottom HelperPlacement = "bottom"
	HelperPlacementTop    HelperPlacement = "top"
	HelperPlacementLeft   HelperPlacement = "left"
	HelperPlacementRight  HelperPlacement = "right"
)

type HelperElementType string

const (
	HelperElementTypeTooltip HelperElementType = "tooltip"
	HelperElementTypePopup   HelperElementType = "popup"
)

type HelperRenderAction string

const (
	HelperRenderActionClick HelperRenderAction = "click"
	HelperRenderActionHover HelperRenderAction = "hover"
)

type HelperBlock struct {
	BlockID string          `json:"blockId" bson:"blockId"`
	Type    HelperBlockType `json:"type" bson:"type"`
	Data    string          `json:"data" bson:"data"`
	Order   int             `json:"order" bson:"order"`
}
