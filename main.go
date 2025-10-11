package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/AliSleiman0/Lacpa/config"
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := config.ConnectMongoDB(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize repository
	database := mongoClient.Database(getEnv("MONGO_DATABASE", "lacpa"))
	repo := repository.NewMongoRepository(database)

	// Initialize handler
	h := handler.NewHandler(repo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Lacpa API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Serve static files from LACPA_Web
	app.Static("/", "./LACPA_Web")

	// API routes
	api := app.Group("/api")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Setup API endpoints
	h.SetupRoutes(api)

	// Start server
	port := getEnv("PORT", "3000")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
