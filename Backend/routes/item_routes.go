package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupItemRoutes configures all item-related API routes
//
// ROLE: Item Route Configuration
// - Defines all HTTP routes for item operations (CRUD)
// - Maps HTTP methods and paths to handler functions
// - Groups related item endpoints under /items prefix
// - Separates routing logic from business logic
//
// ROUTES:
//
//	POST   /items     - Create a new item
//	GET    /items     - Get all items
//	GET    /items/:id - Get item by ID
//	PUT    /items/:id - Update item by ID
//	DELETE /items/:id - Delete item by ID
//
// PARAMETERS:
//
//	router: Fiber router group (typically /api)
//	h: Handler instance containing item business logic
func SetupItemRoutes(router fiber.Router, h *handler.Handler) {
	// Create item routes group under /items
	items := router.Group("/items")

	// Item CRUD operations
	items.Post("/", h.CreateItem)      // POST /api/items
	items.Get("/", h.GetAllItems)      // GET /api/items
	items.Get("/:id", h.GetItem)       // GET /api/items/:id
	items.Put("/:id", h.UpdateItem)    // PUT /api/items/:id
	items.Delete("/:id", h.DeleteItem) // DELETE /api/items/:id
}
