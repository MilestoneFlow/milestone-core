package main

import (
	"context"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"milestone_core/flow"
	"milestone_core/progress"
	"milestone_core/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	flowDbConnection := getFlowDbConnection()
	flowCollection := flowDbConnection.Collection("flows")
	progressCollection := flowDbConnection.Collection("flows_progress")
	templateCollection := flowDbConnection.Collection("flows_templates")

	flowService := flow.Service{Collection: flowCollection}
	progressService := progress.Service{Collection: progressCollection, FlowService: flowService}
	templateService := template.Service{Collection: templateCollection, FlowCollection: flowCollection}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust this based on your specific requirements
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	})

	r.Use(corsMiddleware.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount("/users", usersResource{FlowDb: flowDbConnection}.Routes())
	r.Mount("/todos", todosResource{}.Routes())
	r.Mount("/flows", FlowsResource{
		FlowService:     flowService,
		ProgressService: progressService,
	}.Routes())
	r.Mount("/templates", TemplateResource{
		TemplateService: templateService,
	}.Routes())

	http.ListenAndServe(":3333", r)
}

func getFlowDbConnection() mongo.Database {
	mongoURI := "mongodb://flowAdmin:milestoneFlow123@localhost:27017"
	dbName := "flowDb"

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
