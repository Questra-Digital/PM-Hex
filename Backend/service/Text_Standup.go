package main

import (
	"context"
	"encoding/json"
	"strings"

	//"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StandupRecord struct
type StandupRecord struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Date         string   `json:"date"`
	Timing       string   `json:"timing"`
	Participants string   `json:"participants" bson:"participants"`
	Updates      []string `json:"updates" bson:"updates"`
	Email        string   `json:"email" bson:"email"`
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
			"date": &graphql.Field{
				Type: graphql.String,
			},
			"timing": &graphql.Field{
				Type: graphql.String,
			},
			"participants": &graphql.Field{
				Type: graphql.String,
			},
			"updates": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

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
					"date": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"timing": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"participants": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"updates": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, idOK := p.Args["id"].(string)
					title, titleOK := p.Args["title"].(string)
					date, dateOK := p.Args["date"].(string)
					timing, timingOK := p.Args["timing"].(string)
					participants, participantsOK := p.Args["participants"].(string)
					updates, updatesOK := p.Args["updates"].([]string)
					email, emailOK := p.Args["email"].(string)

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
					} else if dateOK {
						for _, record := range standupRecords {
							if record.Date == date {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					} else if timingOK {
						for _, record := range standupRecords {
							if record.Timing == timing {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					} else if participantsOK {
						for _, record := range standupRecords {
							if record.Participants == participants {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					} else if updatesOK {
						for _, record := range standupRecords {
							found := false
							for _, update := range record.Updates {
								if update == updates[0] { // Assuming you only want to match the first update in the input
									found = true
									break
								}
							}
							if found {
								return []StandupRecord{record}, nil
							}
						}
						return []StandupRecord{}, nil
					} else if emailOK {
						for _, record := range standupRecords {
							if record.Email == email {
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
					"date": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"timing": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"participants": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"updates": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					title, _ := p.Args["title"].(string)
					participants, _ := p.Args["participants"].(string)
					updates, _ := p.Args["updates"].([]string)
					email, _ := p.Args["email"].(string)
					timing := time.Now().Format("3:04 PM")
					date := time.Now().Format("Jan 2, 2006")
					id := fmt.Sprintf("%d", len(standupRecords)+1)
					// Combine the ID, timing, and date into a single string
					// recordID := fmt.Sprintf("%s_%s_%s", id, timing, date)
					newRecord := StandupRecord{
						ID:           id,
						Title:        title,
						Date:         date,
						Timing:       timing,
						Participants: participants,
						Email:        email,
						Updates:      updates,
					}
					standupRecords = append(standupRecords, newRecord)
					return newRecord, nil
				},
			},
			"updateStandupRecord": &graphql.Field{
				Type: StandupRecordType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"date": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"timing": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"participants": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"updates": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Get the ID of the record to be edited
					id := params.Args["id"].(string)

					// Find the record with the specified ID
					for i, record := range standupRecords {
						if record.ID == id {
							// Update the fields of the record if they are present in the mutation arguments
							if title, ok := params.Args["title"].(string); ok {
								standupRecords[i].Title = title
							}
							if date, ok := params.Args["date"].(string); ok {
								standupRecords[i].Date = date
							}
							if timing, ok := params.Args["timing"].(string); ok {
								standupRecords[i].Timing = timing
							}
							if participants, ok := params.Args["participants"].(string); ok {
								standupRecords[i].Participants = participants
							}
							if updates, ok := params.Args["updates"].(string); ok {
								standupRecords[i].Updates = strings.Split(updates, ",")
							}
							if email, ok := params.Args["email"].(string); ok {
								standupRecords[i].Email = email
							}
							return standupRecords[i], nil
						}
					}

					// Return an error if the record with the specified ID was not found
					return nil, fmt.Errorf("standup record with ID %s not found", id)
				},
			},
		},
	},

	// "deleteStandupRecord": &graphql.Field{
	// 	Type: StandupRecordType,
	// 	Args: graphql.FieldConfigArgument{
	// 		"id": &graphql.ArgumentConfig{
	// 			Type: graphql.NewNonNull(graphql.String),
	// 		},
	// 	},
	// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
	// 		id := params.Args["id"].(string)
	// 		for i, record := range standupRecords {
	// 			if record.ID == id {
	// 				// remove the record from the slice
	// 				standupRecords = append(standupRecords[:i], standupRecords[i+1:]...)
	// 				return record, nil
	// 			}
	// 		}
	// 		return nil, errors.New(fmt.Sprintf("Could not find StandupRecord with id %s", id))
	// 	},
	// },

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
			createStandupRecord(
			   title: "My New Standup Record",
			   participants: "Jawad, Hammad, Ayesha",
			   email: "bsce19014@itu.edu.pk, bsce19040@itu.edu.pk",
			   updates: "Text-based standup is done",
			   timing: "10:00 AM",
               date: "2022-03-05") 
			   {
			   id
			   title
			   participants
			   email
			   updates
			   timing
			   date
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
			ID:           "1",
			Title:        "Record 1",
			Date:         "feb 22, 2023",
			Timing:       "5:22 p.m.",
			Participants: " Ayesha, Hammad",
			Updates:      []string{"Initial Setup of meeting"},
			Email:        "ayeshaasmat26@gmail.com, malikhammad90002@gmail.com",
		},
		{
			ID:           "2",
			Title:        "Record 2",
			Date:         "march 22, 2023",
			Timing:       "6:22 p.m.",
			Participants: "Jawad, Junaid",
			Updates:      []string{"Meeting record 2"},
			Email:        "Jawad@questra.digital, Junaid12345@gmail.com ",
		},
	}

	// Initialize GraphQL schema
	Schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    QueryType,
			Mutation: MutationType,
		},
	)
	// // Read the JSON file
	// jsonFile, err := http.Get("http://localhost:8080/graphql")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer jsonFile.Body.Close()

	// // Parse the JSON file into StandupRecord array
	// err = json.NewDecoder(jsonFile.Body).Decode(&standupRecords)
	// if err != nil {
	// 	log.Fatal(err)
	//	}
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
	// Insert the data into the collection
	for _, record := range standupRecords {
		doc := bson.M{
			"id":           record.ID,
			"title":        record.Title,
			"date":         record.Date,
			"timing":       record.Timing,
			"participants": record.Participants,
			"updates":      record.Updates,
			"email":        record.Email,
		}
		_, err = collection.InsertOne(context.Background(), doc)
		if err != nil {
			log.Fatal(err)
		}
	}

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
	fmt.Println("Data inserted successfully")
	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
