package main

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
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

var jwtKey = []byte("secret_key")

const tokenExpireDuration = time.Hour * 3

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

func generateToken(user *User) (string, error) {
	// Set the expiration time for the token
	expirationTime := time.Now().Add(tokenExpireDuration)

	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	})

	// Generate the signed token string using the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

func validateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func sanitizeInput(str string) string {
	return strings.TrimSpace(str)
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, idOk := params.Args["id"].(string)
				email, emailOk := params.Args["email"].(string)

				if !idOk && !emailOk {
					return nil, fmt.Errorf("id or email argument is required")
				}

				// Implement the logic to fetch user information based on ID or email
				var user *User
				for _, u := range users {
					if (idOk && u.ID == id) || (emailOk && u.Email == email) {
						user = &u
						break
					}
				}

				if user == nil {
					return nil, fmt.Errorf("user not found")
				}

				return user, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"signIn": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "SignInResponse",
				Fields: graphql.Fields{
					"user": &graphql.Field{
						Type: userType,
					},
					"token": &graphql.Field{
						Type: graphql.String,
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				email, emailOk := params.Args["email"].(string)
				password, passOk := params.Args["password"].(string)
				if !emailOk || email == "" {
					return nil, fmt.Errorf("email cannot be null")
				}
				if err := validateEmail(email); err != nil {
					return nil, err
				}
				if !passOk || password == "" {
					return nil, fmt.Errorf("password cannot be null")
				}
				if len(password) < 6 {
					return nil, fmt.Errorf("password should be at least 6 characters long")
				}
				email = sanitizeInput(email)
				password = sanitizeInput(password)
				// Check if the user exists and the password is correct
				for _, user := range users {
					if user.Email == email && user.Password == password {
						// Generate a JWT token for the authenticated user
						token, err := generateToken(&user)
						if err != nil {
							return nil, err
						}
						// Return the authenticated user and the token
						response := map[string]interface{}{
							"user":  user,
							"token": token,
						}
						return response, nil
					}
				}
				// Return an error if the email or password is incorrect
				return nil, fmt.Errorf("incorrect email or password")
			},
		},
	},
})

func main() {
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
	}() // Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a handle for your collection
	collection := client.Database("myDatabase").Collection("Credentials")

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

	// Create a new GraphQL schema with the query and mutation types
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GraphQL handler
	graphQLHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Handle the /graphql endpoint
	http.Handle("/signin", graphQLHandler)

	// Start the server
	fmt.Println("Server started on http://localhost:8080/signin")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
