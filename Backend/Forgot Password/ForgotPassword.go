package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var jwtKey = []byte("secret_key")
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Define the GraphQL schema
	fields := graphql.Fields{
		"forgotPassword": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Get the email parameter
				email, ok := params.Args["email"].(string)
				if !ok {
					return nil, nil
				}

				// Check if the email exists in the database
				collection := client.Database("myDatabase").Collection("Credentials")
				filter := bson.M{"email": email}
				var user User
				err := collection.FindOne(context.Background(), filter).Decode(&user)
				if err == mongo.ErrNoDocuments {
					// Email not found
					return nil, nil
				}
				if err != nil {
					log.Fatal(err)
				}

				// Generate a reset token and save it to the database
				resetToken := generateResetToken()
				resetTokenExpiration := time.Now().Add(time.Hour * 1) // Token expires in 24 hours
				update := bson.M{
					"$set": bson.M{
						"reset_token":            resetToken,
						"reset_token_expiration": resetTokenExpiration,
					},
				}
				_, err = collection.UpdateOne(context.Background(), filter, update)
				if err != nil {
					log.Fatal(err)
				}

				// TODO: Send an email with a link to reset the password

				return "An email has been sent to your email address with instructions to reset your password.", nil
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var schema, _ = graphql.NewSchema(graphql.SchemaConfig{})
				// Create a new ResponseWriter for testing purposes
				w := httptest.NewRecorder()
				// Call the authMiddleware function to authenticate the request
				authMiddleware(handler.New(&handler.Config{
					Schema:   &schema,
					Pretty:   true,
					GraphiQL: true,
				})).ServeHTTP(w, params.Info.RootValue.(*http.Request))

				// Check the status code of the response to see if authentication was successful
				if w.Code != http.StatusOK {
					return nil, fmt.Errorf("authentication failed")
				}

				// Call the resolver function to fetch the users from the database
				collection := client.Database("myDatabase").Collection("Credentials")
				cursor, err := collection.Find(context.Background(), bson.M{})
				if err != nil {
					return nil, err
				}
				defer cursor.Close(context.Background())

				var users []User
				for cursor.Next(context.Background()) {
					var user User
					if err := cursor.Decode(&user); err != nil {
						return nil, err
					}
					users = append(users, user)
				}
				if err := cursor.Err(); err != nil {
					return nil, err
				}

				return users, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "Mutation", Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Get the email and password parameters
				email, ok := params.Args["email"].(string)
				if !ok {
					return nil, nil
				}
				password, ok := params.Args["password"].(string)
				if !ok {
					return nil, nil
				}

				// Hash the password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
				}

				// Check if the user already exists
				collection := client.Database("myDatabase").Collection("Credentials")
				filter := bson.M{"email": email}
				var user User
				err = collection.FindOne(context.Background(), filter).Decode(&user)
				if err == nil {
					// User already exists
					return nil, nil
				}
				if err != mongo.ErrNoDocuments {
					log.Fatal(err)
				}

				// Insert the new user into the database
				newUser := User{
					Email:    email,
					Password: string(hashedPassword),
				}
				result, err := collection.InsertOne(context.Background(), newUser)
				if err != nil {
					log.Fatal(err)
				}

				// Generate a JWT token
				claims := &jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hours
					Subject:   fmt.Sprintf("%v", result.InsertedID),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString(jwtKey)
				if err != nil {
					log.Fatal(err)
				}

				// Set a cookie with the token
				http.SetCookie(params.Context.Value("response").(http.ResponseWriter), &http.Cookie{
					Name:     "token",
					Value:    tokenString,
					HttpOnly: true,
				})

				return "User created successfully", nil
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
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Get the email and password parameters
				email, ok := params.Args["email"].(string)
				if !ok {
					return nil, nil
				}
				password, ok := params.Args["password"].(string)
				if !ok {
					return nil, nil
				}

				// Check if the email exists in the database
				collection := client.Database("myDatabase").Collection("Credentials")
				filter := bson.M{"email": email}
				var user User
				err := collection.FindOne(context.Background(), filter).Decode(&user)
				if err == mongo.ErrNoDocuments {
					// Email not found
					return nil, nil
				}
				if err != nil {
					return nil, err // return the error
				}

				// Compare the password with the hash
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err != nil {
					return nil, nil
				}

				// Create a JWT token
				expirationTime := time.Now().Add(time.Minute * 5)
				claims := &jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString(jwtKey)
				if err != nil {
					log.Fatal(err)
				}

				return Token{Token: tokenString}, nil
			},
		},
	}}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/forgotpassword", withResponseWriter(h))

	// Start the server
	log.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Generates a random reset token
func generateResetToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", b)
}

type contextKey string

func withResponseWriter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey("response"), w)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Authorization header is present in the request
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check if the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
