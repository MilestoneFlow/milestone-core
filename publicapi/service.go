package publicapi

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/apiclient"
	"milestone_core/flow"
	"milestone_core/helpers"
	"milestone_core/users"
)

type Service struct {
	ApiClientService    apiclient.Service
	FlowService         flow.Service
	EnrolledUserService users.Service
	HelpersService      helpers.Service
}

func (s Service) GetFlow(token string, id string) (*flow.Flow, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	resFlow, err := s.FlowService.Get(apiClient.WorkspaceID, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return resFlow, err
}

func (s Service) EnrollUser(token string, newUser users.EnrolledUser) error {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return err
	}

	existingUser, err := s.EnrolledUserService.Get(apiClient.WorkspaceID, newUser.ExternalId)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return nil
	}

	newUser.WorkspaceId = apiClient.WorkspaceID
	_, err = s.EnrolledUserService.Create(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetHelpers(token string) ([]helpers.Helper, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	resHelpers, err := s.HelpersService.ListPublished(apiClient.WorkspaceID)
	if err != nil {
		return nil, err
	}

	return resHelpers, nil
}
