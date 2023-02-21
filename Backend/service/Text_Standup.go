package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StandupRecord struct
type StandupRecord struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var standupRecords []StandupRecord

// StandupRecordType defines the GraphQL schema for a StandupRecord
var StandupRecordType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StandupRecord",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// QueryType defines the GraphQL schema for queries
// QueryType defines the GraphQL schema for queries
var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"standupRecords": &graphql.Field{
				Type: graphql.NewList(StandupRecordType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, idOK := p.Args["id"].(string)
					title, titleOK := p.Args["title"].(string)
					if idOK {
						for _, record := range standupRecords {
							if record.ID == id {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					} else if titleOK {
						for _, record := range standupRecords {
							if record.Title == title {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					}
					return standupRecords, nil
				},
			},
		},
	},
)

// MutationType defines the GraphQL schema for mutations
var MutationType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createStandupRecord": &graphql.Field{
				Type: StandupRecordType,
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					title, _ := p.Args["title"].(string)
					id := fmt.Sprintf("%d", len(standupRecords)+1)
					newRecord := StandupRecord{
						ID:    id,
						Title: title,
					}
					standupRecords = append(standupRecords, newRecord)
					return newRecord, nil
				},
			},
		},
	},
)

// Schema defines the GraphQL schema for the API
var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	},
)

func runExampleMutation() {

	mutationQuery := `mutation createStandupRecord {
		createStandupRecord(title: "My New Standup Record") {
			id
			title
		}
	}
	`
	params := graphql.Params{
		Schema:         Schema,
		RequestString:  mutationQuery,
		VariableValues: nil,
		OperationName:  "createStandupRecord",
		Context:        context.Background(),
	}

	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		fmt.Printf("mutation failed: %v\n", result.Errors)
	} else {
		fmt.Printf("mutation result: %v\n", result.Data)
	}
}

func main() {
	// Initialize StandupRecords array
	standupRecords = []StandupRecord{
		{
			ID:    "1",
			Title: "Record 1",
		},
		{
			ID:    "2",
			Title: "Record 2",
		},
	}

	// Initialize GraphQL schema
	Schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    QueryType,
			Mutation: MutationType,
		},
	)

	// Register a handler for the "/graphql" endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        Schema,
			RequestString: r.URL.Query().Get("query"),
		})
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Errors)
			return
		}
		json.NewEncoder(w).Encode(result)
	})

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a handle for your collection
	collection := client.Database("myDatabase").Collection("myCollection")

	// Query the collection
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	// Run an example mutation
	runExampleMutation()

	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
