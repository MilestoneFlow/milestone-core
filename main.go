package main

import (
	"context"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/apiclient"
	"milestone_core/authorization"
	"milestone_core/flow"
	"milestone_core/helpers"
	"milestone_core/publicapi"
	"milestone_core/template"
	"milestone_core/tracker"
	"milestone_core/users"
	"milestone_core/workspace"
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

	templateCollection := flowDbConnection.Collection("flows_templates")
	usersCollection := flowDbConnection.Collection("enrolled_users")
	branchingCollection := flowDbConnection.Collection("branching")
	apiClientsCollection := flowDbConnection.Collection("api_clients")
	usersStateCollection := flowDbConnection.Collection("users_state")
	workspaceCollection := flowDbConnection.Collection("workspaces")
	helpersCollection := flowDbConnection.Collection("helpers")
	trackerCollection := flowDbConnection.Collection("tracking_data")

	log.Default().Print("collections initialized")

	flowService := flow.Service{Collection: flowCollection, ArchiveCollection: flowArchiveCollection}
	templateService := template.Service{Collection: templateCollection, FlowCollection: flowCollection}
	usersService := users.Service{Collection: usersCollection, UserStateCollection: usersStateCollection}
	branchingService := flow.BranchingService{Collection: branchingCollection}
	apiClientService := apiclient.Service{Collection: apiClientsCollection}
	workspaceService := workspace.Service{Collection: workspaceCollection}
	helpersService := helpers.Service{Collection: helpersCollection}
	flowEnroller := flow.Enroller{Collection: flowCollection}
	publicapiService := publicapi.Service{
		ApiClientService:    apiClientService,
		FlowEnroller:        flowEnroller,
		EnrolledUserService: usersService,
		HelpersService:      helpersService,
	}
	trackerService := tracker.Tracker{Collection: trackerCollection}
	flowAnalyticsService := flow.Analytics{Tracker: trackerService}

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

	r.Mount("/enrolled-users", usersResource{usersService: usersService}.Routes())
	r.Mount("/todos", todosResource{}.Routes())
	r.Mount("/flows", flow.FlowsResource{
		FlowService: flowService,
		Analytics:   flowAnalyticsService,
	}.Routes())
	r.Mount("/helpers", helpers.Resource{
		Service: helpersService,
	}.Routes())
	r.Mount("/templates", TemplateResource{
		TemplateService: templateService,
	}.Routes())
	r.Mount("/branching", BranchingResource{
		BranchingService: branchingService,
	}.Routes())
	r.Mount("/workspaces", workspace.Resource{
		Service: workspaceService,
	}.Routes())
	r.Mount("/auth", AuthResource{}.Routes())
	r.Mount("/apiclients", ApiClientResource{
		Service: apiClientService,
	}.Routes())
	r.Mount("/public", PublicApiResource{
		Service: publicapiService,
		UserStateService: publicapi.UserStateService{
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
