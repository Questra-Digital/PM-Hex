package main

import (
	"context"
	"encoding/json"

	"fmt"

	"net/http"

	"github.com/graphql-go/graphql"
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
	mutationQuery := `
        mutation {
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
		OperationName:  "",
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

	fmt.Println("Listening on localhost:8080")

	// Call the example mutation query
	runExampleMutation()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
