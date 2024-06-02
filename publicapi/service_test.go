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
	SetupMockData(&mongoConnection)

	service := Service{
		ApiClientService: apiclient.Service{
			Collection: mongoConnection.Collection("api_clients"),
		},
		FlowEnroller:        flow.Enroller{Collection: mongoConnection.Collection("flows")},
		EnrolledUserService: users.Service{Collection: mongoConnection.Collection("enrolled_users")},
	}

	t.Run("sanity test", func(t *testing.T) {
		newId := primitive.NewObjectID()
		mongoConnection.Collection("flows").InsertOne(context.Background(), flow.Flow{
			ID:          newId,
			WorkspaceID: WorkspaceID(),
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

	t.Run("sanity test", func(t *testing.T) {

		resFlow, err := service.EnrollInFlow("token", "userId")
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

	CleanupMockData(&mongoConnection)
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
