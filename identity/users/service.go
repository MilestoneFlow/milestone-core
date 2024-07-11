package users

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/jmoiron/sqlx"
	"os"
)

type Service struct {
	DbConnection  *sqlx.DB
	CognitoClient *cognitoidentityprovider.Client
}

func (s Service) GetWorkspaceUsers(workspaceId string) ([]User, error) {
	var workspaceUsers []User
	query := "SELECT u.id, u.created_at FROM identity.platform_user u JOIN identity.workspace_user wu ON u.id = wu.user_id WHERE wu.workspace_id = $1"
	err := s.DbConnection.Select(&workspaceUsers, query, workspaceId)
	if err != nil {
		return nil, err
	}

	for i, _ := range workspaceUsers {
		workspaceUsers[i], err = s.getUserDetails(workspaceUsers[i])
		if err != nil {
			return nil, err
		}
	}

	return workspaceUsers, nil
}

func (s Service) GetWorkspaceInvitedUsers(workspaceId string) ([]InvitedUser, error) {
	var workspaceUsers []InvitedUser
	query := "SELECT workspace_id, email, token, created_at FROM identity.workspace_user_invite wu WHERE wu.workspace_id = $1"
	err := s.DbConnection.Select(&workspaceUsers, query, workspaceId)

	return workspaceUsers, err
}

func (s Service) getUserDetails(user User) (User, error) {
	cognitoPoolId := os.Getenv("AWS_COGNITO_USER_POOL_ID")

	userDetails, err := s.CognitoClient.AdminGetUser(context.TODO(), &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(cognitoPoolId),
		Username:   aws.String(user.ID),
	})
	if err != nil {
		return user, err
	}

	userName := ""
	userEmail := ""
	for _, attr := range userDetails.UserAttributes {
		if *attr.Name == "email" {
			userEmail = *attr.Value
		}
		if *attr.Name == "name" {
			userName = *attr.Value
		}
	}

	user.Details = UserDetails{
		Email: userEmail,
		Name:  userName,
	}

	return user, nil
}
