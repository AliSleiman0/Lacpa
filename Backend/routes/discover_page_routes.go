package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupDiscoverPageRoutes(router fiber.Router, handler *handler.DiscoverPageHandler) {
	// Create discover page routes group
	discover := router.Group("/discover")

	// Discover page operations
	discover.Get("/", handler.HandleDiscoverPage) // GET /main/discover - Render discover page
}
