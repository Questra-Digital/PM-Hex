package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
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

func generateToken(user *User) (map[string]interface{}, error) {
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
		return nil, err
	}

	return map[string]interface{}{
		"token": tokenString,
	}, nil
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
				email, emailOk := params.Args["email"].(string)
				password, passOk := params.Args["password"].(string)
				if !emailOk || email == "" {
					return nil, fmt.Errorf("email cannot be null")
				}

				if !passOk || password == "" {
					return nil, fmt.Errorf("password cannot be null")
				}
				if len(password) < 6 {
					return nil, fmt.Errorf("password should be at least 6 characters long")
				}

				// Check if the user exists and the password is correct
				for _, user := range users {
					if user.Email == email && user.Password == password {
						// Generate a JWT token for the authenticated user
						token, err := generateToken(&user)
						if err != nil {
							return nil, err
						}

						// Return the authenticated user and the token
						return map[string]interface{}{
							"user":  user,
							"token": token,
						}, nil
					}
				}

				// Return an error if the email or password is incorrect
				return nil, fmt.Errorf("incorrect email or password")
			},
		},
	},
})

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Parse the JWT token using the secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			http.Error(w, "Invalid Authorization token", http.StatusUnauthorized)
			return
		}

		// Call the next middleware or handler function
		next.ServeHTTP(w, r)
	})
}

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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST"},
	})
	handlerWithCORS := c.Handler(h)

	http.Handle("/graphql", handlerWithCORS)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the user exists and the password is correct
		for _, u := range users {
			if u.Email == user.Email && u.Password == user.Password {
				// Generate a JWT token for the authenticated user
				token, err := generateToken(&u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Return the authenticated user and the token
				response := map[string]interface{}{
					"user":  u,
					"token": token,
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		// Return an error if the email or password is incorrect
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
	})

	log.Fatal(s.ListenAndServe())
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

	// Wrap the GraphQL schema handler with the authentication middleware
	authHandler := authMiddleware(h)

	// Register the authentication-wrapped GraphQL schema handler with the /graphql endpoint
	http.Handle("/graphql", authHandler)

	fmt.Println("Data inserted successfully")

	fmt.Println("Server is running on port 8080")
	s.ListenAndServe()
}
