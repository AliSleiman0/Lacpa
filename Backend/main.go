package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/AliSleiman0/Lacpa/config"
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/AliSleiman0/Lacpa/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
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

	// Initialize repositories
	database := mongoClient.Database(getEnv("MONGO_DATABASE", "lacpa"))
	repo := repository.NewMongoRepository(database)
	authRepo := repository.NewAuthRepository(database)

	// Initialize HTML template engine
	// Templates will be loaded from "./templates" directory
	engine := html.New("./templates", ".html")

	// Add template functions for pagination and helpers
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})
	engine.AddFunc("iterate", func(count int) []int {
		var i int
		var items []int
		for i = 0; i < count; i++ {
			items = append(items, i)
		}
		return items
	})

	// Enable development mode for template reloading (optional)
	if getEnv("APP_ENV", "development") == "development" {
		engine.Reload(true)
		engine.Debug(true)
	}

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		AppName: "Lacpa API",
		Views:   engine,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With, HX-Request, HX-Trigger, HX-Target, HX-Current-URL, HX-Boosted, HX-History-Restore-Request",
	}))

	// Add no-cache headers for all static files during development
	app.Use(func(c *fiber.Ctx) error {
		// Set no-cache headers for JS, CSS, and HTML files
		if getEnv("APP_ENV", "development") == "development" {
			c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")
		}
		return c.Next()
	})

	// Serve static files from LACPA_Web (parent directory)
	// Disable caching for development
	app.Static("/", "../LACPA_Web", fiber.Static{
		CacheDuration: 0,
		MaxAge:        0,
	})

	// Setup all API routes (includes health check and all endpoints)
	routes.SetupRoutes(app, repo)

	// Setup authentication routes
	authHandler := handler.NewAuthHandler(authRepo)
	routes.SetupAuthRoutes(app, authHandler)

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
