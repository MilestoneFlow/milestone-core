package publicapi

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/apiclient"
	"milestone_core/flow"
	"milestone_core/users"
	"testing"
)

func TestGroup(t *testing.T) {
	mongoConnection := getFlowDbConnection()
	mongoConnection.Collection("flows").Drop(context.Background())
	mongoConnection.Collection("api_clients").Drop(context.Background())

	mongoConnection.Collection("api_clients").InsertOne(context.Background(), apiclient.ApiClient{
		Token:       "token",
		WorkspaceID: "testWorkspaceId",
	})

	service := Service{
		ApiClientService: apiclient.Service{
			Collection: mongoConnection.Collection("api_clients"),
		},
		FlowService:         flow.Service{Collection: mongoConnection.Collection("flows")},
		EnrolledUserService: users.Service{Collection: mongoConnection.Collection("enrolled_users")},
	}

	t.Run("sanity test", func(t *testing.T) {
		t.Cleanup(func() {
			mongoConnection.Collection("flows").DeleteMany(context.Background(), primitive.M{})
		})

		newId := primitive.NewObjectID()
		mongoConnection.Collection("flows").InsertOne(context.Background(), flow.Flow{
			ID:          newId,
			WorkspaceID: "testWorkspaceId",
			Name:        "testName",
			BaseURL:     "testBaseURL",
			Segments:    []flow.Segment{},
			Steps:       []flow.Step{},
			Relations:   []flow.Relation{},
			Opts:        flow.Opts{},
			Live:        true,
		})

		resFlow, err := service.GetFlow("token", newId.Hex())
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		if resFlow == nil {
			t.Fatalf("Flow is nil")
		}
		if resFlow.Name != "testName" {
			t.Fatalf("Flow name is not the expected one")
		}
	})
	t.Run("B", func(t *testing.T) {
		t.Logf("Passed")
	})
}

func getFlowDbConnection() mongo.Database {
	mongoURI := "mongodb://flowAdmin:milestoneFlow123@localhost:27018"
	dbName := "flowDb"

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return *client.Database(dbName)
}
