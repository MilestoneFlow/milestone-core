package authorization

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"milestone_core/shared/rest"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func CognitoMiddleware(postgresConnection *sqlx.DB, cognitoClient *cognitoidentityprovider.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		servePublicApiRequest := func(w http.ResponseWriter, r *http.Request, authToken string) {
			workspaceID, err := GetWorkspaceIDByPublicApiToken(postgresConnection, authToken)
			if workspaceID == nil {
				rest.SendErrorResponse(w, errors.New("invalid token"), http.StatusUnauthorized)
				return
			}
			if err != nil {
				log.Default().Print(err)
				rest.SendErrorResponse(w, errors.New("server error"), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "user", UserData{
				WorkspaceID: *workspaceID,
				Token:       authToken,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/health") {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			authToken := strings.TrimPrefix(authHeader, "Bearer ")
			if authToken == authHeader {
				http.Error(w, "Bearer token not found", http.StatusUnauthorized)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/public/") || strings.HasPrefix(r.URL.Path, "/api/v1/") {
				servePublicApiRequest(w, r, authToken)
				return
			}

			user, err := cognitoClient.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
				AccessToken: aws.String(authToken),
			})
			if err != nil {
				log.Default().Print(err)
				rest.SendErrorResponse(w, errors.New("invalid token"), http.StatusUnauthorized)
				return
			}

			cognitoId := *user.Username
			userEmailIdentifier := ""
			for _, attr := range user.UserAttributes {
				if *attr.Name == "email" {
					userEmailIdentifier = *attr.Value
					break
				}
			}

			workspaceID := ""
			headerWorkspaceId := r.Header.Get("Workspace-Id")
			if headerWorkspaceId != "" {
				_, err = uuid.Parse(headerWorkspaceId)
				if err != nil {
					rest.SendErrorResponse(w, errors.New("invalid workspace id provided"), http.StatusBadRequest)
					return
				}

				hasAccess, err := UserHasAccessToWorkspace(postgresConnection, cognitoId, headerWorkspaceId)
				if err != nil {
					log.Default().Print(err)
					rest.SendErrorResponse(w, errors.New("error checking access to workspace"), http.StatusInternalServerError)
					return
				}
				if !hasAccess {
					rest.SendErrorResponse(w, errors.New("user does not have access to workspace"), http.StatusForbidden)
					return
				}

				workspaceID = headerWorkspaceId
			} else {
				workspaceID, err = GetWorkspaceIDByUserIdentifier(postgresConnection, cognitoId)
				if err != nil {
					log.Default().Print(err)
					rest.SendErrorResponse(w, errors.New("error checking access to workspace"), http.StatusInternalServerError)
					return
				}
			}

			if workspaceID == "" {
				workspaceID, err = CreateDefaultWorkspaceForUser(cognitoId, postgresConnection)
				if err != nil {
					log.Default().Print(err)
					rest.SendErrorResponse(w, errors.New("error checking access to workspace"), http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", UserData{
				WorkspaceID: workspaceID,
				UserID:      *user.Username,
				Email:       userEmailIdentifier,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
