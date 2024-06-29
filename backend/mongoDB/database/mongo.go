package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	BooksCollection *mongo.Collection
	AuthorCollection *mongo.Collection
)

func Connect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Get value from .env
    MONGO_URI := os.Getenv("MONGO_URI")

    // Connect to the database.
    clientOptions := options.Client().ApplyURI(MONGO_URI)
    localClient, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
	Client = localClient
	
	BooksCollection = Client.Database("database").Collection("books")
	AuthorCollection = Client.Database("database").Collection("authors")

	// Check the connection
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to mongoDB!!!")
	}
}

func Close() error {
	return Client.Disconnect(context.Background())
}