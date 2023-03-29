package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              string `json:"id"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"-" validate:"required,min=8"`
	ConfirmPassword string `json:"-" validate:"required,eqfield=Password"`
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
					// Create a new validator
					validate := validator.New()
					// Parse the user argument from the input parameters
					var user User
					err := json.Unmarshal([]byte(p.Args["user"].(string)), &user)
					if err != nil {
						return nil, fmt.Errorf("error parsing user argument: %s", err)
					}
					// Validate the user argument
					err = validate.Struct(user)
					if err != nil {
						return nil, err
					}
					// Check if the passwords match
					if user.Password != user.ConfirmPassword {
						return nil, fmt.Errorf("passwords do not match")
					}
					// Check if the user exists and the password is correct
					for _, u := range users {
						if u.Email == user.Email && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil {
							return u, nil
						}
					}
					// // Initialize the users variable
					// users := []User{}
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
					filter := bson.M{"email": user.Email}
					err = collection.FindOne(ctx, filter).Decode(&User{})
					if err == nil {
						return nil, fmt.Errorf("email address already in use")
					}

					// Hash the password
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
					if err != nil {
						return nil, fmt.Errorf("error hashing password")
					}
					// Generate a new UUID
					id := uuid.New().String()
					// Create a new user
					newUser := &User{
						ID:       id,
						Email:    user.Email,
						Password: string(hashedPassword),
					}
					// Insert the user into the database
					_, err = collection.InsertOne(ctx, user)
					if err != nil {
						return nil, fmt.Errorf("error inserting user into database")
					}
					// Append the new user to the list of users
					users = append(users, newUser)

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
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Execute the GraphQL query and return the result
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
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
	collection := client.Database("myDatabase").Collection("Credentials")
	// Insert the data into the collection
	for _, record := range users {
		doc := bson.M{
			"id":              record.ID,
			"email":           record.Email,
			"password":        record.Password,
			"confirmpassword": record.ConfirmPassword,
		}
		_, err = collection.InsertOne(context.Background(), doc)
		if err != nil {
			log.Fatal(err)
		}
	}
	// // Create a new HTTP handler for the GraphQL endpoint
	// handler := graphql.NewHandler(graphql.HandlerConfig{
	// 	Schema: &schema,
	// })
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

	fmt.Println("Data inserted successfully")

	// Start the server
	//http.Handle("/graphql", handler)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
