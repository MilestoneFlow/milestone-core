package workspace

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func TestGroup(t *testing.T) {
	mongoConnection := getFlowDbConnection()
	workspacesCollection := mongoConnection.Collection("workspaces")
	workspacesCollection.Drop(context.Background())
	service := Service{
		Collection: workspacesCollection,
	}

	insertedWorkspaceRow, err := workspacesCollection.InsertOne(context.Background(), Workspace{
		Name:            "Test Workspace",
		BaseURL:         "email.com",
		UserIdentifiers: []string{"test@email.com", "hello@email.com"},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	insertedWorkspaceID := insertedWorkspaceRow.InsertedID.(primitive.ObjectID).Hex()

	t.Run("Get workspace by ID", func(t *testing.T) {
		resWorkspace, err := service.Get(insertedWorkspaceID)

		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if resWorkspace == nil {
			t.Fatalf("Workspace is nil")
		}

		if resWorkspace.Name != "Test Workspace" {
			t.Fatalf("Workspace name is not the expected one")
		}
	})

	t.Run("Get workspace for existing user identified", func(t *testing.T) {
		resWorkspace, err := service.GetByUserIdentifier("hello@email.com")

		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if resWorkspace == nil {
			t.Fatalf("Workspace is nil")
		}

		if resWorkspace.Name != "Test Workspace" {
			t.Fatalf("Workspace name is not the expected one")
		}
	})

	t.Run("User identifier is not in workspace, return nil", func(t *testing.T) {
		resWorkspace, err := service.GetByUserIdentifier("nonexistent@email.com")

		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if resWorkspace != nil {
			t.Fatalf("Workspace is not nil")
		}
	})

	t.Run("Add user identifier in workspace", func(t *testing.T) {
		newUserIdentifier := "newuser@email.com"
		err := service.AddUserIdentifier(insertedWorkspaceID, newUserIdentifier)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		resWorkspace, err := service.GetByUserIdentifier(newUserIdentifier)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if resWorkspace == nil {
			t.Fatalf("Workspace is nil")
		}

		if resWorkspace.Name != "Test Workspace" {
			t.Fatalf("Workspace name is not the expected one")
		}
	})

	t.Run("Remove user identifier from workspace", func(t *testing.T) {
		userIdentifierToRemove := "test@email.com"
		resWorkspace, err := service.GetByUserIdentifier(userIdentifierToRemove)
		if resWorkspace == nil {
			t.Fatalf("Workspace is nil")
		}

		err = service.RemoveUserIdentifier(insertedWorkspaceID, userIdentifierToRemove)

		resWorkspace, err = service.GetByUserIdentifier(userIdentifierToRemove)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if resWorkspace != nil {
			t.Fatalf("User wasn't deleted from workspace")
		}
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