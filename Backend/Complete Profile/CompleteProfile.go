package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define user struct
type UserProfile struct {
	ID        int    `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Age       int    `bson:"age"`
}

// Counter for generating unique IDs
var idCounter = 0

var collection *mongo.Collection

// Create a new user
func createUser(params graphql.ResolveParams) (interface{}, error) {
	firstName, _ := params.Args["firstName"].(string)
	lastName, _ := params.Args["lastName"].(string)
	age, _ := params.Args["age"].(int)

	if firstName == "" {
		return nil, errors.New("firstName is required")
	}

	if lastName == "" {
		return nil, errors.New("lastName is required")
	}

	if age <= 0 {
		return nil, errors.New("age must be greater than 0")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Increment the ID counter and use the new value as the user's ID
	idCounter++
	user := UserProfile{
		ID:        idCounter,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
	// Insert the user into the database
	collection := client.Database("myDatabase").Collection("Profile")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by ID
func getUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(string)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Find the user by ID
	var user UserProfile
	collection := client.Database("myDatabase").Collection("Profile")
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Define GraphQL schema
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"firstName": &graphql.Field{Type: graphql.String},
		"lastName":  &graphql.Field{Type: graphql.String},
		"age":       &graphql.Field{Type: graphql.Int},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: getUser,
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"firstName": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"lastName":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"age":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: createUser,
		},
	},
})

// Initialize MongoDB connection
func initMongoDB() error {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	// Set user collection
	collection = client.Database("myDatabase").Collection("Profile")

	return nil
}

// Initialize GraphQL server
func initGraphQL() error {
	// Define schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		return err
	}
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

// Main function
func main() {
	// Initialize MongoDB connection
	err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize GraphQL server
	err = initGraphQL()
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	log.Printf("Server started on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
