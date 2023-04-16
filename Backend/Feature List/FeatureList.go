package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/graphql-go/handler"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feature struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var currentID int

var client *mongo.Client
var featureCollection *mongo.Collection

func initMongoDB() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Set up feature collection
	featureCollection = client.Database("myDatabase").Collection("FeatureList")
}

var featureType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Feature",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"features": &graphql.Field{
			Type: graphql.NewList(featureType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Query MongoDB for all features
				cursor, err := featureCollection.Find(context.Background(), bson.M{})
				if err != nil {
					return nil, err
				}
				defer cursor.Close(context.Background())

				// Convert the cursor to a slice of Features
				var features []Feature
				for cursor.Next(context.Background()) {
					var feature Feature
					err := cursor.Decode(&feature)
					if err != nil {
						return nil, err
					}
					features = append(features, feature)
				}
				fmt.Println(featureCollection) // added statement

				return features, nil
			},
		},
	},
})
var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createFeature": &graphql.Field{
			Type: featureType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Extract the name argument
				name, ok := p.Args["name"].(string)
				if !ok {
					return nil, fmt.Errorf("missing or invalid 'name' parameter")
				}

				// Increment the current ID count
				currentID++

				// Create a new Feature object
				feature := Feature{
					ID:   strconv.Itoa(currentID),
					Name: name,
				}

				// Insert the Feature into MongoDB
				_, err := featureCollection.InsertOne(context.Background(), feature)
				if err != nil {
					return nil, err
				}

				return feature, nil
			},
		},
	},
})
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    queryType,
	Mutation: mutationType,
})

// Initialize GraphQL server
func initGraphQL() error {

	// Define GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Define HTTP server and route handler
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})

	return nil
}

func main() {
	// Initialize MongoDB
	initMongoDB()

	// Initialize GraphQL server
	err := initGraphQL()
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	log.Printf("Server started on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
