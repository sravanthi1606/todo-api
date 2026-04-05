package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGODB_URL := os.Getenv("MONGODB_URL")
	clientOptions := options.Client().ApplyURI(MONGODB_URL)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Successfully")

	collection = client.Database("golang_db").Collection("todos")

	// ✅ Gin setup
	app := gin.Default()

	// ✅ CORS middleware for Gin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	setupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// ✅ Correct for Gin
	app.Run(":" + port)
}
