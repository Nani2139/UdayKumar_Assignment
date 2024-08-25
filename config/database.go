package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is the global MongoDB database connection object.
// It's used to interact with the MongoDB instance.
var DB *mongo.Database

// ConnectDB initializes the MongoDB client and connects to the specified database.
// It sets the global DB variable to the connected database instance.
func ConnectDB() {
	// Define client options, including the URI to the MongoDB server and a connection timeout.
	clientOptions := options.Client().ApplyURI("mongodb+srv://udayiiitl039:nani39@cluster0.vfzi9kb.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetConnectTimeout(10 * time.Second)

	// Create a new MongoDB client using the defined options.
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Establish a context with a timeout for the connection process.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server.
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB server: %v", err)
	}

	// Verify the connection by pinging the MongoDB server.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB server: %v", err)
	}

	// Assign the connected database to the global DB variable.
	DB = client.Database("Cluster0")
	fmt.Println("Connected to MongoDB!")
}
