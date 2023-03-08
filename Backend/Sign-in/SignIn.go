package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

var users = []User{
	{
		ID:       "1",
		Email:    "ayeshaasmat26@gmail.com",
		Password: "password1",
	},
	{
		ID:       "2",
		Email:    "Hammad@gmail.com",
		Password: "password2",
	},
}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"signIn": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				email, _ := params.Args["email"].(string)
				password, _ := params.Args["password"].(string)

				// Check if the user exists and the password is correct
				for _, user := range users {
					if user.Email == email && user.Password == password {
						return user, nil
					}
				}

				// Return an error if the email or password is incorrect
				return nil, fmt.Errorf("incorrect email or password")
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    nil,
	Mutation: rootMutation,
})

func main() {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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
	for _, record := range users {
		doc := bson.M{
			"id":       record.ID,
			"email":    record.Email,
			"password": record.Password,
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

	fmt.Println("Server is running on port 8080")
	s.ListenAndServe()
}

func runExampleMutation() {
	panic("unimplemented")
}
