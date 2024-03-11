package main

import (
	"context"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/authorization"
	"milestone_core/flow"
	"milestone_core/progress"
	"milestone_core/template"
	"milestone_core/users"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Default().Print("a mers")

	flowDbConnection := getFlowDbConnection()
	log.Default().Print("a mers si db")

	flowCollection := flowDbConnection.Collection("flows")
	progressCollection := flowDbConnection.Collection("flows_progress")
	templateCollection := flowDbConnection.Collection("flows_templates")
	usersCollection := flowDbConnection.Collection("enrolled_users")
	branchingCollection := flowDbConnection.Collection("branching")

	log.Default().Print("a mers si colectii")

	flowService := flow.Service{Collection: flowCollection}
	progressService := progress.Service{Collection: progressCollection, FlowService: flowService}
	templateService := template.Service{Collection: templateCollection, FlowCollection: flowCollection}
	usersService := users.Service{Collection: usersCollection}
	branchingService := flow.BranchingService{Collection: branchingCollection}

	r := chi.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust this based on your specific requirements
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	})

	r.Use(corsMiddleware.Handler)
	r.Use(authorization.CognitoMiddleware("us-east-1_zrIqQshjP", "us-east-1"))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount("/enrolled-users", usersResource{usersService: usersService}.Routes())
	r.Mount("/todos", todosResource{}.Routes())
	r.Mount("/flows", FlowsResource{
		FlowService:     flowService,
		ProgressService: progressService,
	}.Routes())
	r.Mount("/templates", TemplateResource{
		TemplateService: templateService,
	}.Routes())
	r.Mount("/branching", BranchingResource{
		BranchingService: branchingService,
	}.Routes())
	r.Mount("/auth", AuthResource{}.Routes())

	http.ListenAndServe(":3333", r)
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
