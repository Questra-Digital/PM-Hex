package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmpassword" validate:"required,eqfield=Password"`
}

var users []*User

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.ID},
			"email": &graphql.Field{Type: graphql.String},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:    graphql.NewList(userType),
				Resolve: getAllUsers,
			},
		},
	},
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Get the secret key from the environment variable
			secretKey := os.Getenv("JWT_SECRET_KEY")
			if secretKey == "" {
				return nil, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
			}
			// Parse the secret key
			return []byte(secretKey), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		type ClaimKey string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), ClaimKey("claims"), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}

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
					// Call connectToDatabase function to connect to the database
					client, err := connectToDatabase()
					if err != nil {
						return nil, err
					}
					defer client.Disconnect(context.Background())

					// Create a new validator
					validate := validator.New()
					// Parse the user argument from the input parameters
					var user User
					err = json.Unmarshal([]byte(p.Args["user"].(string)), &user)
					if err != nil {
						return nil, err
					}

					// Validate the user input
					err = validate.Struct(user)
					if err != nil {
						return nil, err
					}
					// Check if the user already exists
					for _, u := range users {
						if u.Email == user.Email {
							return nil, fmt.Errorf("user already exists")
						}
					}
					// Generate a new UUID for the user
					user.ID = uuid.New().String()
					// Hash the user's password
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
					if err != nil {
						return nil, err
					}
					user.Password = string(hashedPassword)

					// Create the user
					err = createUser(user)
					if err != nil {
						return nil, err
					}

					// Add the user to the list of users
					users = append(users, &user)
					// Return the newly created user
					return &user, nil
				},
			},

			"login": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					// Find the user with the given email
					var user *User
					for _, u := range users {
						if u.Email == email {
							user = u
							break
						}
					}
					if user == nil {
						return nil, fmt.Errorf("user not found")
					}
					// Compare the password hash with the provided password
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
					if err != nil {
						return nil, err
					}
					// Generate a new JWT token
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"sub": user.ID,
						"exp": time.Now().Add(time.Hour * 24).Unix(),
					})
					// Get the secret key from the environment variable
					secretKey := os.Getenv("JWT_SECRET_KEY")
					if secretKey == "" {
						return nil, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
					}
					// Sign the token with the secret key
					tokenString, err := token.SignedString([]byte(secretKey))
					if err != nil {
						return nil, err
					}
					// Return the token
					return tokenString, nil
				},
			},
		},
	},
)

func main() {
	// Create a new GraphQL schema with the rootQuery and rootMutation
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rootQuery,
			Mutation: rootMutation,
		},
	)

	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	// Create a new HTTP server with the GraphQL endpoint
	http.Handle("/graphql", AuthMiddleware(&graphqlHandler{Schema: &schema}))
	// Start the HTTP server
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type graphqlHandler struct {
	Schema *graphql.Schema
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	// result := graphql.Do(graphql.Query(r.URL.Query().Get("query")).Schema(*h.Schema))
	result := graphql.Do(graphql.Params{
		Schema:        *h.Schema,
		RequestString: query,
	})

	// Check for errors in the result
	if len(result.Errors) > 0 {
		log.Printf("failed to execute query: %v", result.Errors)
		http.Error(w, "failed to execute query", http.StatusBadRequest)
		return
	}
	// Encode the result as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func connectToDatabase() (*mongo.Client, error) {
	// Get the MongoDB connection string from the environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable not set")
	}
	// Set the options for the MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * time.Second)
	// Create a new MongoDB client
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	return client, nil
}

func createUser(user User) error {
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	// Create a new user document
	userDoc := bson.M{
		"_id":             uuid.New().String(),
		"email":           user.Email,
		"hashed_password": hashedPassword,
		"created_at":      time.Now().UTC(),
		"updated_at":      time.Now().UTC(),
	}

	// Get the MongoDB client
	client, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	// Get the users collection
	usersCollection := client.Database("myDatabase").Collection("Credentials")
	// Insert the user document
	_, err = usersCollection.InsertOne(context.Background(), userDoc)
	if err != nil {
		return fmt.Errorf("failed to insert user document: %v", err)
	}
	return nil
}
func getAllUsers(p graphql.ResolveParams) (interface{}, error) {
	// Call connectToDatabase function to connect to the database
	client, err := connectToDatabase()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Get the users collection from the database
	collection := client.Database("myDatabase").Collection("Credentials")

	// Query the users collection to get all documents
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode the cursor into a slice of users
	var users []User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	// Return the list of users
	return users, nil
}
