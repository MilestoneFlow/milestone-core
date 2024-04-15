package authorization

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"milestone_core/awsinternal"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func CognitoMiddleware(workspaceCollection *mongo.Collection, region string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
				ctx := context.WithValue(r.Context(), "token", authToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			cfg, err := awsinternal.GetConfiguration(region)
			if err != nil {
				http.Error(w, "failed to load AWS configuration", http.StatusInternalServerError)
				return
			}

			client := cognitoidentityprovider.NewFromConfig(*cfg)
			user, err := client.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
				AccessToken: aws.String(authToken),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
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

			splitEmail := strings.Split(userEmailIdentifier, "@")

			if workspaceID == "" {
				workspaceCollection.InsertOne(context.Background(), bson.M{
					"userIdentifiers": []string{userEmailIdentifier},
					"name":            "My workspace",
					"baseUrl":         "https://" + splitEmail[1],
				})
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

func GetEmailDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[1] == "" {
		return ""
	}
	return parts[1]
}
