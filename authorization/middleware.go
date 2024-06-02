package authorization

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"milestone_core/awsinternal"
	"milestone_core/rest"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func CognitoMiddleware(apiClientsCollection *mongo.Collection, workspaceCollection *mongo.Collection, region string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		servePublicApiRequest := func(w http.ResponseWriter, r *http.Request, authToken string) {
			workspaceID, err := GetWorkspaceIDByPublicApiToken(apiClientsCollection, authToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "client", PublicApiClientData{
				WorkspaceID: workspaceID,
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

			if strings.HasPrefix(r.URL.Path, "/public/") {
				servePublicApiRequest(w, r, authToken)
				return
			}

			cfg, err := awsinternal.GetConfiguration(region)
			if err != nil {
				log.Default().Print("authorizer: failed to load AWS configuration")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			client := cognitoidentityprovider.NewFromConfig(*cfg)
			user, err := client.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
				AccessToken: aws.String(authToken),
			})
			if err != nil {
				log.Default().Print(err)
				rest.SendErrorResponse(w, errors.New("invalid token"), http.StatusUnauthorized)
				return
			}

			userEmailIdentifier := ""
			for _, attr := range user.UserAttributes {
				if *attr.Name == "email" {
					userEmailIdentifier = *attr.Value
					break
				}
			}

			workspaceID, err := GetWorkspaceIDByUserIdentifier(workspaceCollection, userEmailIdentifier)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if workspaceID == "" {
				splitEmail := strings.Split(userEmailIdentifier, "@")
				inserted, err := workspaceCollection.InsertOne(context.Background(), bson.M{
					"userIdentifiers": []string{userEmailIdentifier},
					"name":            "My workspace",
					"baseUrl":         "https://" + splitEmail[1],
				})
				if err != nil {
					log.Default().Print(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				workspaceID = inserted.InsertedID.(primitive.ObjectID).Hex()
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
