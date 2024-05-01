package publicapi

import (
	"errors"
	"log"
	"milestone_core/apiclient"
	"milestone_core/flow"
	"milestone_core/users"
	"slices"
)

type Service struct {
	ApiClientService    apiclient.Service
	FlowService         flow.Service
	EnrolledUserService users.Service
}

func (s Service) GetFlow(token string, id string) (*flow.Flow, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	resFlow, err := s.FlowService.Get(apiClient.WorkspaceID, id)
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

func (s Service) GetUserState(token string, userId string) (*users.UserState, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	lastUserState, err := s.EnrolledUserService.GetLastUserState(apiClient.WorkspaceID, userId)
	if err != nil {
		return nil, err
	}

	if lastUserState == nil || lastUserState.CurrentStepId == "" {
		finishedFlowIds, err := s.EnrolledUserService.GetFinishedFlowsForUser(apiClient.WorkspaceID, userId)
		if err != nil {
			return nil, err
		}

		flows, err := s.FlowService.ListLive(apiClient.WorkspaceID)
		if err != nil {
			log.Panic(err)
			return nil, nil
		}

		for _, inputFlow := range flows {
			if slices.Contains(finishedFlowIds, inputFlow.ID.Hex()) {
				continue
			}

			prerequisite, err := s.FlowService.GetFlowPrerequisite(apiClient.WorkspaceID, inputFlow.ID.Hex())
			if err != nil {
				return nil, err
			}

			for _, flowId := range prerequisite {
				if !slices.Contains(finishedFlowIds, flowId) {
					continue
				}
			}

			source := inputFlow.Steps[0]
			for _, step := range inputFlow.Steps {
				if step.Opts.IsSource == true {
					source = step
					break
				}
			}

			lastUserState = &users.UserState{
				UserId:                userId,
				CurrentEnrolledFlowId: inputFlow.ID.Hex(),
				CurrentStepId:         source.StepID,
				Times:                 make([]users.UserStateTime, 0),
			}
			_, err = s.EnrolledUserService.CreateUserState(apiClient.WorkspaceID, *lastUserState)
			return lastUserState, nil
		}

		return nil, nil
	}

	return lastUserState, nil
}

func (s Service) EnrollInNextFlow(token string, userId string) (*users.UserState, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	lastUserState, err := s.EnrolledUserService.GetLastUserState(apiClient.WorkspaceID, userId)
	if err != nil {
		return nil, err
	}

	if lastUserState == nil || lastUserState.CurrentStepId == "" {
		return nil, nil
	}

	return lastUserState, nil
}

func (s Service) UpdateUserState(token string, userState users.UserState) error {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return err
	}

	existingUser, err := s.EnrolledUserService.Get(apiClient.WorkspaceID, userState.UserId)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not enrolled")
	}

	_, err = s.EnrolledUserService.CreateUserState(apiClient.WorkspaceID, userState)
	if err != nil {
		return err
	}

	return nil
}
