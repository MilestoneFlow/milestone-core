package main

import (
	"context"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/identity/apiclient"
	"milestone_core/identity/authorization"
	"milestone_core/identity/workspace"
	"milestone_core/public/apigateway"
	"milestone_core/public/enrolledusers"
	"milestone_core/tours/branching"
	"milestone_core/tours/flows"
	"milestone_core/tours/helpers"
	"milestone_core/tours/tracker"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	flowDbConnection := getFlowDbConnection()
	log.Default().Print("flowDb connected")

	flowCollection := flowDbConnection.Collection("flows")
	flowArchiveCollection := flowDbConnection.Collection("flows_archive")

	usersCollection := flowDbConnection.Collection("enrolled_users")
	branchingCollection := flowDbConnection.Collection("branching")
	apiClientsCollection := flowDbConnection.Collection("api_clients")
	usersStateCollection := flowDbConnection.Collection("users_state")
	workspaceCollection := flowDbConnection.Collection("workspaces")
	helpersCollection := flowDbConnection.Collection("helpers")
	trackerCollection := flowDbConnection.Collection("tracking_data")

	log.Default().Print("collections initialized")

	flowService := flows.Service{Collection: flowCollection, ArchiveCollection: flowArchiveCollection}
	usersService := enrolledusers.Service{Collection: usersCollection, UserStateCollection: usersStateCollection}
	branchingService := flows.BranchingService{Collection: branchingCollection}
	apiClientService := apiclient.Service{Collection: apiClientsCollection}
	workspaceService := workspace.Service{Collection: workspaceCollection}
	helpersService := helpers.Service{Collection: helpersCollection}
	flowEnroller := flows.Enroller{Collection: flowCollection}
	publicapiService := apigateway.Service{
		ApiClientService:    apiClientService,
		FlowEnroller:        flowEnroller,
		EnrolledUserService: usersService,
		HelpersService:      helpersService,
	}
	trackerService := tracker.Tracker{Collection: trackerCollection}
	flowAnalyticsService := flows.Analytics{Tracker: trackerService}

	r := chi.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust this based on your specific requirements
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	})

	r.Use(corsMiddleware.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authorizer := authorization.CognitoMiddleware(apiClientsCollection, workspaceCollection, "us-east-1")
	r.Use(authorizer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount("/enrolled-users", enrolledusers.UsersResource{UsersService: usersService}.Routes())
	r.Mount("/flows", flows.FlowsResource{
		FlowService: flowService,
		Analytics:   flowAnalyticsService,
	}.Routes())
	r.Mount("/helpers", helpers.Resource{
		Service: helpersService,
	}.Routes())
	r.Mount("/branching", branching.BranchingResource{
		BranchingService: branchingService,
	}.Routes())
	r.Mount("/workspaces", workspace.Resource{
		Service: workspaceService,
	}.Routes())
	r.Mount("/apiclients", apiclient.ApiClientResource{
		Service: apiClientService,
	}.Routes())
	r.Mount("/public", apigateway.PublicApiResource{
		Service: publicapiService,
		UserStateService: apigateway.UserStateService{
			UserStateCollection: usersStateCollection,
			ApiClientService:    apiClientService,
			EnrolledUserService: usersService,
		},
		Tracker: trackerService,
	}.Routes())

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Default().Print(err)
		return
	}
}

func getFlowDbConnection() mongo.Database {
	mongoURI := os.Getenv("FLOW_DB_CONNECTION_URL")
	dbName := os.Getenv("FLOW_DB_NAME")

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a handle for your collection
	return *client.Database(dbName)
}
