package publicapi

import (
	"errors"
	"milestone_core/apiclient"
	"milestone_core/flow"
	"milestone_core/helpers"
	"milestone_core/users"
)

type Service struct {
	ApiClientService    apiclient.Service
	FlowEnroller        flow.Enroller
	EnrolledUserService users.Service
	HelpersService      helpers.Service
}

func (s Service) ValidateToken(token string) error {
	_, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetFlow(token string, id string) (*flow.Flow, error) {
	apiClient, err := s.ApiClientService.GetByToken(token)
	if err != nil {
		return nil, err
	}

	resFlow, err := s.FlowEnroller.GetFlow(apiClient.WorkspaceID, flow.EnrollmentOpts{
		CurrentEnrollmentId: id,
	})
	if err != nil {
		return nil, err
	}

	return resFlow, err
}

func (s Service) EnrollInFlow(workspaceId string, externalUserId string) (*flow.Flow, error) {
	enrolledUser, err := s.EnrolledUserService.Get(workspaceId, externalUserId)
	if err != nil {
		return nil, err
	}
	if enrolledUser == nil {
		return nil, errors.New("user not found")
	}

	userState, err := s.EnrolledUserService.GetState(workspaceId, enrolledUser.ID.Hex())
	if err != nil {
		return nil, err
	}

	resFlow, err := s.FlowEnroller.GetFlow(workspaceId, flow.EnrollmentOpts{
		CurrentEnrollmentId: userState.FlowsData.CurrentFlowID,
		FinishedIds:         userState.FlowsData.CompletedFlowsIds,
		SkippedIds:          userState.FlowsData.SkippedFlowsIds,
	})
	if err != nil {
		return nil, err
	}
	if resFlow == nil {
		return nil, nil
	}

	userState.FlowsData.CurrentFlowID = resFlow.ID.Hex()
	err = s.EnrolledUserService.PutState(workspaceId, enrolledUser.ID.Hex(), *userState)

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
	err = s.EnrolledUserService.Create(newUser)
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

func (s Service) UpdateUserStateFromTrackEvents(workspaceId string, externalUserId string, events []EventTrack) error {
	skippedFlowId := ""
	skippedTimestamp := int64(0)
	finishedFlowId := ""
	finishedTimestamp := int64(0)
	for _, event := range events {
		if event.EventType == EventTypeFlowSkipped {
			skippedFlowId = event.EntityID
			skippedTimestamp = event.Timestamp
		}
		if event.EventType == EventTypeFlowFinished {
			finishedFlowId = event.EntityID
			finishedTimestamp = event.Timestamp
		}
	}

	if skippedFlowId == "" && finishedFlowId == "" {
		return nil
	}

	enrolledUser, err := s.EnrolledUserService.Get(workspaceId, externalUserId)
	if err != nil {
		return err
	}
	if enrolledUser == nil {
		return errors.New("user not found")
	}

	currentState, err := s.EnrolledUserService.GetState(workspaceId, enrolledUser.ID.Hex())
	if err != nil {
		return err
	}

	if skippedFlowId != "" {
		currentState.FlowsData.SkippedFlowsIds = s.getUniqueValuesFromArr(append(currentState.FlowsData.SkippedFlowsIds, skippedFlowId))
		if currentState.FlowsData.CurrentFlowID == skippedFlowId {
			currentState.FlowsData.CurrentFlowID = ""
		}
		if skippedTimestamp > currentState.FlowsData.LastSubmittedFlowTimestamp {
			currentState.FlowsData.LastSubmittedFlowTimestamp = skippedTimestamp
			currentState.FlowsData.LastSubmittedFlowID = skippedFlowId
		}
	}
	if finishedFlowId != "" {
		currentState.FlowsData.CompletedFlowsIds = s.getUniqueValuesFromArr(append(currentState.FlowsData.CompletedFlowsIds, finishedFlowId))
		if currentState.FlowsData.CurrentFlowID == finishedFlowId {
			currentState.FlowsData.CurrentFlowID = ""
		}
		if finishedTimestamp > currentState.FlowsData.LastSubmittedFlowTimestamp {
			currentState.FlowsData.LastSubmittedFlowTimestamp = finishedTimestamp
			currentState.FlowsData.LastSubmittedFlowID = finishedFlowId
		}
	}
	currentState.FlowsData.SkippedFlowsIds = s.excludeValuesFromArr(currentState.FlowsData.SkippedFlowsIds, currentState.FlowsData.CompletedFlowsIds)

	err = s.EnrolledUserService.PutState(workspaceId, enrolledUser.ID.Hex(), *currentState)

	return err
}

func (s Service) getUniqueValuesFromArr(arr []string) []string {
	uniqueValues := make(map[string]bool)
	for _, v := range arr {
		uniqueValues[v] = true
	}

	var res []string
	for k := range uniqueValues {
		res = append(res, k)
	}

	return res
}

func (s Service) excludeValuesFromArr(arr []string, excluded []string) []string {
	excludedValues := make(map[string]bool)
	for _, v := range excluded {
		excludedValues[v] = true
	}

	var res []string
	for _, v := range arr {
		if !excludedValues[v] {
			res = append(res, v)
		}
	}

	return res
}
