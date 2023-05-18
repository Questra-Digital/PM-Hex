package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	//"log"
	"net/http"
	//"os"
	"time"

	//"github.com/shurcooL/graphql/playground"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RemoteMeeting struct {
	ID            primitive.ObjectID `bson:"_id"`
	Title         string             `bson:"title"`
	TimeZone      string             `bson:"timeZone"`
	Participants  []string           `bson:"participants"`
	Questions     []string           `bson:"questions"`
	Responses     []string           `bson:"responses"`
	Status        string             `bson:"status"`
	ScheduledTime time.Time          `bson:"scheduledTime"`
	ReminderTime  time.Time          `bson:"reminderTime"`
	Messages      []Message          `bson:"messages"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	Sender    string             `bson:"sender"`
	Text      string             `bson:"text"`
	CreatedAt time.Time          `bson:"createdAt"`
}

var meetingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Meeting",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.String},
		"title":         &graphql.Field{Type: graphql.String},
		"timeZone":      &graphql.Field{Type: graphql.String},
		"participants":  &graphql.Field{Type: graphql.NewList(graphql.String)},
		"questions":     &graphql.Field{Type: graphql.NewList(graphql.String)},
		"responses":     &graphql.Field{Type: graphql.NewList(graphql.String)},
		"status":        &graphql.Field{Type: graphql.String},
		"scheduledTime": &graphql.Field{Type: graphql.DateTime},
		"reminderTime":  &graphql.Field{Type: graphql.DateTime},
		"messages":      &graphql.Field{Type: graphql.NewList(messageType)},
	},
})

var messageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Message",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"sender":    &graphql.Field{Type: graphql.String},
		"text":      &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},
	},
})

var collection *mongo.Collection

func main() {
	// slackToken := os.Getenv("SLACK_TOKEN")
	// if slackToken == "" {
	// 	log.Fatal("SLACK_TOKEN environment variable is not set")
	// }

	// Connect to the MongoDB database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection = client.Database("myDatabase").Collection("RemoteMeetings")

	// Define the mutation schema
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createMeeting": &graphql.Field{
				Type:        meetingType,
				Description: "Create a new remote meeting",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"timeZone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"participants": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"questions": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"responses": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"scheduledTime": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"reminderTime": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: createMeeting,
			},
		},
	})

	// Define the query schema
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"meeting": &graphql.Field{
				Type:        meetingType,
				Description: "Retrieve a meeting by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: getMeetingByID,
			},
			"meetings": &graphql.Field{
				Type:        graphql.NewList(meetingType),
				Description: "Retrieve all meetings",
				Resolve:     getAllMeetings,
			},
		},
	})

	// Define the schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}

	// Define the GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var requestBody struct {
				Query string `json:"query"`
			}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return

			}

			// Execute the GraphQL query
			result := executeQuery(requestBody.Query, schema)
			json.NewEncoder(w).Encode(result)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	// Define the GraphQL Playground handler
	// graphqlPlaygroundHandler := graphiql.New(&graphiql.Config{
	// 	Endpoint: "/graphql",
	// })

	// // Serve the GraphQL Playground at a specific endpoint
	// http.Handle("/playground", graphqlPlaygroundHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
	// Define the GraphQL Playground handler

}
func sendMessageToSlack(channel, message string) error {
	slackWebhookURL := "YOUR_SLACK_WEBHOOK_URL"
	payload := map[string]string{
		"channel": channel,
		"text":    message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack: %s", resp.Status)
	}

	return nil
}


// executeQuery executes the given GraphQL query against the provided schema.
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("Unexpected errors inside executeQuery: %v", result.Errors)
	}

	return result
}

// createMeeting creates a new remote meeting in the database.
func createMeeting(params graphql.ResolveParams) (interface{}, error) {
	mutation := `
    mutation {
        createMeeting(
            title: "My Meeting"
            timeZone: "UTC"
            participants: ["John", "Jane"]
            questions: ["What did you do yesterday?", "What are you planning to do today?", "Are there any roadblocks or obstacles?"]
            responses: []
            status: "active"
            scheduledTime: "2023-05-12T12:00:00Z"
            reminderTime: "2023-05-12T11:30:00Z"
        ) {
            id
            title
            scheduledTime
            participants
            questions
            responses
            status
            messages {
                id
                sender
                text
                createdAt
            }
        }
    }
    `
	// You can use the "mutation" variable here if needed
	fmt.Println(mutation)
	title := params.Args["title"].(string)
	timeZone := params.Args["timeZone"].(string)
	participants := params.Args["participants"].([]interface{})
	questions := params.Args["questions"].([]interface{})
	responses := params.Args["responses"].([]interface{})
	status := params.Args["status"].(string)
	scheduledTimeString := params.Args["scheduledTime"].(string)
	reminderTimeString := params.Args["reminderTime"].(string)

	// Parse scheduledTime and reminderTime strings to time.Time format
	scheduledTime, err := time.Parse(time.RFC3339, scheduledTimeString)
	if err != nil {
		return nil, err
	}
	reminderTime, err := time.Parse(time.RFC3339, reminderTimeString)
	if err != nil {
		return nil, err
	}

	// Create a new RemoteMeeting instance
	meeting := RemoteMeeting{
		ID:            primitive.NewObjectID(),
		Title:         title,
		TimeZone:      timeZone,
		Participants:  toStringSlice(participants),
		Questions:     toStringSlice(questions),
		Responses:     toStringSlice(responses),
		Status:        status,
		ScheduledTime: scheduledTime,
		ReminderTime:  reminderTime,
		Messages:      []Message{},
	}
	// Insert the meeting into the database
	insertResult, err := collection.InsertOne(context.Background(), meeting)
	if err != nil {
		return nil, err
	}

	// Fetch the inserted meeting document from the database
	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	var insertedMeeting RemoteMeeting
	err = collection.FindOne(context.Background(), filter).Decode(&insertedMeeting)
	if err != nil {
		return nil, err
	}

	return insertedMeeting, nil
}

// getMeetingByID retrieves a meeting from the database by its ID.
func getMeetingByID(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	// Parse the ID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Fetch the meeting document from the database
	filter := bson.D{{Key: "_id", Value: objectID}}
	var meeting RemoteMeeting
	err = collection.FindOne(context.Background(), filter).Decode(&meeting)
	if err != nil {
		return nil, err
	}

	return meeting, nil
}

// getAllMeetings retrieves all meetings from the database.
func getAllMeetings(params graphql.ResolveParams) (interface{}, error) {
	// Fetch all meeting documents from the database
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var meetings []RemoteMeeting
	for cursor.Next(context.Background()) {
		var meeting RemoteMeeting
		if err := cursor.Decode(&meeting); err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

// toStringSlice converts a slice of interface{} to a slice of strings.
func toStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = item.(string)
	}
	return result
}
