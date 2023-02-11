package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Helper function to import json from file to map
func importJSONDataFromFile(filename string, result interface{}) (isOK bool) {
	isOK = true
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false

	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}

var StandupList []Standup
var _ = importJSONDataFromFile("./StandupData.json", &StandupList)

type Standup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// define custom GraphQL ObjectType `StandupType` for our Golang struct `Standup`
// Note that
// - the fields in our todoType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var StandupType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Standup Text",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var currentMaxId = 4

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"addStandup": &graphql.Field{
			Type:        StandupType, // the return type for this field
			Description: "add a new Standup Entry",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				name, _ := params.Args["name"].(string)
				// figure out new id
				newID := currentMaxId + 1
				currentMaxId = currentMaxId + 1

				// perform mutation operation here
				// for e.g. create a Standup and save to DB.
				newStandup := Standup{
					ID:   newID,
					Name: name,
				}

				StandupList = append(StandupList, newStandup)

				// return the new Standup object that we supposedly save to DB
				// Note here that
				// - we are returning a `Standup` struct instance here
				// - we previously specified the return Type to be `StandupType`
				// - `Standup` struct maps to `StandupType`, as defined in `StandupType` ObjectConfig`
				return newStandup, nil
			},
		},
		"updateStandup": &graphql.Field{
			Type:        StandupType, // the return type for this field
			Description: "Update existing Standup ",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				affectedStandup := Standup{}

				// Search list for Standup with id
				for i := 0; i < len(StandupList); i++ {
					if StandupList[i].ID == id {
						if _, ok := params.Args["name"]; ok {
							StandupList[i].Name = params.Args["name"].(string)
						}
						// Assign updated Standup so we can return it
						affectedStandup = StandupList[i]
						break
					}
				}
				// Return affected Standup
				return affectedStandup, nil
			},
		},
	},
})

// root query
// test with Sandbox at localhost:8080/sandbox
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"Standup": &graphql.Field{
			Type:        StandupType,
			Description: "Get single Standup",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				nameQuery, isOK := params.Args["name"].(string)
				if isOK {
					// Search for el with name
					for _, Standup := range StandupList {
						if Standup.Name == nameQuery {
							return Standup, nil
						}
					}
				}

				return Standup{}, nil
			},
		},

		"StandupList": &graphql.Field{
			Type:        graphql.NewList(StandupType),
			Description: "List of Text Standups",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return StandupList, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var StandupSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	h := handler.New(&handler.Config{
		Schema:   &StandupSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	http.Handle("/graphql", h)

	http.Handle("/sandbox", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(sandboxHTML)
	}))

	http.ListenAndServe(":8080", nil)

}

var sandboxHTML = []byte(`<!DOCTYPE html>
<html lang="en">
<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
<div id="sandbox" style="height:100vh; width:100vw;"></div>
<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
<script>
 new window.EmbeddedSandbox({
   target: "#sandbox",
   // Pass through your server href if you are embedding on an endpoint.
   // Otherwise, you can pass whatever endpoint you want Sandbox to start up with here.
   initialEndpoint: "http://localhost:8080/graphql",
 });
 // advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
</html>`)
