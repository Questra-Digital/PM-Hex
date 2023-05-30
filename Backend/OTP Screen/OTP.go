package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	PhoneNumber string `bson:"phoneNumber"`
	OTP         string `bson:"otp"`
	ExpiredAt   int64  `bson:"expiredAt"`
}

var collection *mongo.Collection

func init() {
	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Set up MongoDB collection
	collection = client.Database("myDatabase").Collection("users")
}

// ...

func main() {
	// Define GraphQL schema
	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"phoneNumber": &graphql.Field{
				Type: graphql.String,
			},
			"otp": &graphql.Field{
				Type: graphql.String,
			},
			"expiredAt": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Fetch all user documents from MongoDB collection
					cursor, err := collection.Find(context.Background(), bson.M{})
					if err != nil {
						return nil, err
					}

					var users []User
					if err := cursor.All(context.Background(), &users); err != nil {
						return nil, err
					}

					return users, nil
				},
			},
		},
	})

	var rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"sendOTP": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Mutation to send OTP code to user's phone number",
				Args: graphql.FieldConfigArgument{
					"phoneNumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					phoneNumber, _ := p.Args["phoneNumber"].(string)
					name, _ := p.Args["name"].(string)

					// Generate random OTP code
					rand.Seed(time.Now().UnixNano())
					otp := strconv.Itoa(rand.Intn(9999-1000+1) + 1000)

					// Set OTP expiration time to 5 minutes from now
					expiredAt := time.Now().Add(5 * time.Minute).Unix()
					// Create user document
					user := User{
						ID:          phoneNumber,
						Name:        name,
						PhoneNumber: phoneNumber,
						OTP:         otp,
						ExpiredAt:   expiredAt,
					}

					// Insert or update user document in MongoDB collection
					update := bson.M{"$set": user}
					_, err := collection.UpdateOne(context.Background(), bson.M{"_id": phoneNumber}, update, options.Update().SetUpsert(true))
					if err != nil {
						return false, err
					}

					return true, nil
				},
			},
			"verifyOTP": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Mutation to verify OTP code sent to user's phone number",
				Args: graphql.FieldConfigArgument{
					"phoneNumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"otp": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					phoneNumber, _ := p.Args["phoneNumber"].(string)
					otp, _ := p.Args["otp"].(string)
					// Get user document from MongoDB collection
					filter := bson.M{"_id": phoneNumber, "otp": otp, "expiredAt": bson.M{"$gt": time.Now().Unix()}}

					var result User
					err := collection.FindOne(context.Background(), filter).Decode(&result)
					if err != nil {
						return false, err
					}

					// Update user document to mark OTP as used
					update := bson.M{"$unset": bson.M{"otp": ""}}
					_, err = collection.UpdateOne(context.Background(), filter, update)
					if err != nil {
						return false, err
					}

					return true, nil
				},
			},
			"updateUserProfile": &graphql.Field{
				Type:        userType,
				Description: "Mutation to update user profile information",
				Args: graphql.FieldConfigArgument{
					"phoneNumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					phoneNumber, _ := p.Args["phoneNumber"].(string)
					name, _ := p.Args["name"].(string)

					update := bson.M{"$set": bson.M{"name": name}}

					_, err := collection.UpdateOne(context.Background(), bson.M{"_id": phoneNumber}, update)
					if err != nil {
						return nil, err
					}

					// Fetch the updated user document from MongoDB collection
					var user User
					err = collection.FindOne(context.Background(), bson.M{"_id": phoneNumber}).Decode(&user)
					if err != nil {
						return nil, err
					}

					return user, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		// Handle the error appropriately, such as logging it or returning an error response.
		fmt.Println("Error creating GraphQL schema:", err)
		return
	}
	// Create a GraphQL handler
	graphQLHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Serve the GraphQL API and GraphiQL playground
	mux := http.NewServeMux()
	mux.Handle("/otp", graphQLHandler)

	// Serve the GraphiQL playground at the root path ("/")
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/otp", http.StatusSeeOther)
	}))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
	if err != nil {
		log.Fatal(err)
	}

}
