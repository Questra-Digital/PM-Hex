package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"-"`
	ConfirmPassword string `json:"-"`
}

var users []*User

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":              &graphql.Field{Type: graphql.ID},
			"email":           &graphql.Field{Type: graphql.String},
			"password":        &graphql.Field{Type: graphql.String},
			"confirmpassword": &graphql.Field{Type: graphql.String},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello, world!", nil
				},
			},
		},
	},
)

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"user": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					confirmPassword := p.Args["confirmPassword"].(string)

					if password != confirmPassword {
						return nil, fmt.Errorf("passwords do not match")
					}
					// Check if the user exists and the password is correct
					for _, user := range users {
						if user.Email == email && user.Password == password {
							return user, nil
						}
					}
					// Initialize the users variable
					users := []User{}
					// Connect to the database
					client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
					if err != nil {
						log.Fatal(err)
					}

					// Connect to the MongoDB database
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					err = client.Connect(ctx)
					if err != nil {
						log.Fatal(err)
					}

					// Get a handle to the users collection
					collection := client.Database("myDatabase").Collection("Credentials")

					// Check if the email address is already in use
					filter := bson.M{"email": email}
					err = collection.FindOne(ctx, filter).Decode(&User{})
					if err == nil {
						return nil, fmt.Errorf("email address already in use")
					}

					// Hash the password
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					if err != nil {
						return nil, fmt.Errorf("error hashing password")
					}

					// Create a new user
					user := User{
						ID:       "1",
						Email:    email,
						Password: string(hashedPassword),
					}
					// Insert the user into the database
					_, err = collection.InsertOne(ctx, user)
					if err != nil {
						return nil, fmt.Errorf("error inserting user into database")
					}

					// Append the new user to the list of users
					users = append(users, user)

					// Return the new user
					return user, nil

				},
			},
		},
	},
)
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	// Initialize the users slice
	users = []*User{}
	// Create a new HTTP handler for the GraphQL endpoint
	handler := graphql.Handler{
		Schema: &schema,
	}
	// Start the server
	http.Handle("/graphql", &handler)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
