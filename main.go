package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/gamification/events"
	"milestone_core/gamification/rewards"
	"milestone_core/identity/apiclient"
	"milestone_core/identity/authorization"
	"milestone_core/identity/users"
	"milestone_core/identity/workspace"
	"milestone_core/public/apigateway"
	"milestone_core/public/enrolledusers"
	"milestone_core/shared/awsinternal"
	"milestone_core/shared/rest"
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
	err := checkRequiredEnvVars()
	if err != nil {
		log.Panic(err)
		return
	}

	flowDbConnection := getFlowDbConnection()
	log.Default().Print("mongo connected")
	postgresConnection := getPostgresConnection()
	log.Default().Print("postgres connected")

	awsCfg, err := awsinternal.GetConfiguration("us-east-1")
	if err != nil {
		log.Panic(err)
		return
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsCfg)

	flowCollection := flowDbConnection.Collection("flows")
	flowArchiveCollection := flowDbConnection.Collection("flows_archive")
	usersCollection := flowDbConnection.Collection("enrolled_users")
	branchingCollection := flowDbConnection.Collection("branching")
	usersStateCollection := flowDbConnection.Collection("users_state")
	helpersCollection := flowDbConnection.Collection("helpers")
	trackerCollection := flowDbConnection.Collection("tracking_data")

	flowService := flows.Service{Collection: flowCollection, ArchiveCollection: flowArchiveCollection}
	enrolledUsersService := enrolledusers.Service{Collection: usersCollection, UserStateCollection: usersStateCollection}
	branchingService := flows.BranchingService{Collection: branchingCollection}
	apiClientService := apiclient.Service{DbConnection: postgresConnection}
	usersService := users.Service{DbConnection: postgresConnection, CognitoClient: cognitoClient}
	workspaceService := workspace.Service{DbConnection: postgresConnection, UsersService: usersService}
	helpersService := helpers.Service{Collection: helpersCollection}
	flowEnroller := flows.Enroller{Collection: flowCollection}
	publicapiService := apigateway.Service{
		ApiClientService:    apiClientService,
		FlowEnroller:        flowEnroller,
		EnrolledUserService: enrolledUsersService,
		HelpersService:      helpersService,
	}
	trackerService := tracker.Tracker{Collection: trackerCollection}
	flowAnalyticsService := flows.Analytics{Tracker: trackerService}

	eventsResource := events.Resource{
		EventsService: events.Service{
			DbConnection: postgresConnection,
		},
	}
	rewardsResource := rewards.Resource{Service: rewards.Service{DbConnection: postgresConnection}}

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

	authorizer := authorization.CognitoMiddleware(postgresConnection, cognitoClient)
	r.Use(authorizer)
	r.Use(rest.RequestLoggerMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount("/enrolled-users", enrolledusers.UsersResource{UsersService: enrolledUsersService}.Routes())
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
			EnrolledUserService: enrolledUsersService,
		},
		Tracker: trackerService,
	}.Routes())

	r.Mount("/events", eventsResource.Routes())
	r.Mount("/rewards", rewardsResource.Routes())

	r.Mount("/api/v1/events", eventsResource.PublicRoutes())

	port := "3333"
	err = http.ListenAndServe(":"+port, r)
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

func getPostgresConnection() *sqlx.DB {
	postgresURI := os.Getenv("POSTGRES_DB_URI")
	db, err := sqlx.Connect("postgres", postgresURI)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return db
}

func checkRequiredEnvVars() error {
	if os.Getenv("AWS_COGNITO_USER_POOL_ID") == "" {
		return errors.New("AWS_COGNITO_USER_POOL_ID is not set")
	}

	return nil
}
