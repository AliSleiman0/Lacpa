package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes is the main route configuration function that sets up all API routes
//
// ROLE: Central Route Orchestrator
// - Acts as the single entry point for all route configuration
// - Creates and initializes all handler instances
// - Delegates route setup to specific route configuration functions
// - Maintains separation of concerns by grouping related routes
// - Makes it easy to add new route groups and maintain existing ones
//
// ARCHITECTURE:
//
//	This function follows the composition pattern where:
//	1. Handlers are created with repository dependencies
//	2. Route groups are configured by calling specific setup functions
//	3. Each route group is responsible for its own route definitions
//
// EXTENSIBILITY:
//
//	To add new routes (e.g., Orders):
//	1. Create routes/order_routes.go with SetupOrderRoutes function
//	2. Create handler/order_handler.go if needed
//	3. Add SetupOrderRoutes(api, orderHandler) call below
//
// PARAMETERS:
//
//	app: Main Fiber application instance
//	repo: Unified repository interface providing all data access methods
func SetupRoutes(app *fiber.App, repo repository.Repository) {
	// Web routes (HTML templates) - these serve HTML pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":   "Home",
			"BaseURL": "http://" + c.Get("Host"),
		}, "layouts/main")
	})

	app.Get("/items", func(c *fiber.Ctx) error {
		return c.Render("items", fiber.Map{
			"Title": "Items Management",
		}, "layouts/main")
	})

	// Create API route group - all API routes will be under /api prefix
	api := app.Group("/api")

	// Health check endpoint - simple endpoint to verify server is running
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Setup route groups - each function handles its own route configuration
	mainPageHandler := handler.NewMainPageHandler(repo)
	SetupMainPageRoutes(api, mainPageHandler) // Configures /api/main/* routes

	// Future route groups can be added here:
	// SetupItemRoutes(api, itemHandler)       // Would configure /api/items/* routes
	// SetupOrderRoutes(api, orderHandler)     // Would configure /api/orders/* routes
	// SetupProductRoutes(api, productHandler) // Would configure /api/products/* routes
}
