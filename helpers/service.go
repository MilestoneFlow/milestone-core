package helpers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) Get(publicId string, workspaceId string) (*Helper, error) {
	var helper Helper
	err := s.Collection.FindOne(context.Background(), bson.M{"publicId": publicId, "workspaceId": workspaceId}).Decode(&helper)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("helper not found")
	}
	if err != nil {
		return nil, err
	}

	return &helper, nil
}

func (s Service) List(workspaceId string) ([]Helper, error) {
	var helpers []Helper
	cursor, err := s.Collection.Find(context.Background(), bson.M{"workspaceId": workspaceId})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &helpers)
	if err != nil {
		return nil, err
	}

	if helpers == nil {
		helpers = []Helper{}
	}

	return helpers, nil
}

func (s Service) Create(workspaceId string, inputHelper Helper) (*Helper, error) {
	newHelper := s.createNewHelper(workspaceId, inputHelper)

	_, err := s.Collection.InsertOne(context.Background(), newHelper)
	return newHelper, err
}

func (s Service) Update(publicId string, workspaceId string, helper map[string]interface{}) error {
	_, err := s.Get(publicId, workspaceId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("helper not found")
	}
	if err != nil {
		return err
	}

	helper["updated"] = time.Now().Unix()
	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"publicId": publicId, "workspaceId": workspaceId}, bson.M{"$set": helper})
	return err
}

func (s Service) Delete(publicId string, workspaceId string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.M{"publicId": publicId, "workspaceId": workspaceId})
	return err
}

func (s Service) Publish(publicId string, workspaceId string) error {
	_, err := s.Get(publicId, workspaceId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("helper not found")
	}
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"publicId": publicId, "workspaceId": workspaceId}, bson.M{"$set": bson.M{"published": true, "publishedAt": time.Now().Unix()}})
	return err
}

func (s Service) Unpublish(publicId string, workspaceId string) error {
	_, err := s.Get(publicId, workspaceId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("helper not found")
	}
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(context.Background(), bson.M{"publicId": publicId, "workspaceId": workspaceId}, bson.M{"$set": bson.M{"published": false}, "$unset": bson.M{"publishedAt": true}})
	return err
}

func (s Service) createNewHelper(workspaceId string, override Helper) *Helper {
	created := time.Now().Unix()
	defaultHelper := Helper{
		PublicID:    uuid.New().String(),
		WorkspaceID: workspaceId,
		Name:        "New untitled helper",
		Data: HelperData{
			TargetUrl:          "",
			AssignedCssElement: "",
			ElementType:        HelperElementTypePopup,
			Placement:          HelperPlacementBottom,
			Blocks: []HelperBlock{
				{
					BlockID: "block_1",
					Type:    HelperBlockTypeText,
					Data:    "This is a new helper",
					Order:   1,
				},
			},
			ActionText: "Close",
			IconColor:  "#f5f5f5",
		},
		RenderAction: HelperRenderActionClick,
		Published:    false,
		Created:      created,
	}

	helper := defaultHelper
	if override.Name != "" {
		helper.Name = override.Name
	}
	if override.Data.TargetUrl != "" {
		helper.Data.TargetUrl = override.Data.TargetUrl
	}
	if override.Data.AssignedCssElement != "" {
		helper.Data.AssignedCssElement = override.Data.AssignedCssElement
	}
	if override.Data.ElementType != "" {
		helper.Data.ElementType = override.Data.ElementType
	}
	if override.Data.Placement != "" {
		helper.Data.Placement = override.Data.Placement
	}
	if override.Data.Blocks != nil && len(override.Data.Blocks) > 0 {
		helper.Data.Blocks = override.Data.Blocks
	}
	if override.Data.ActionText != "" {
		helper.Data.ActionText = override.Data.ActionText
	}
	if override.Data.IconColor != "" {
		helper.Data.IconColor = override.Data.IconColor
	}
	if override.RenderAction != "" {
		helper.RenderAction = override.RenderAction
	}

	return &helper
}
