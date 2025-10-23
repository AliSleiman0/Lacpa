package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupNewsDetailsRoutes(router fiber.Router, handler *handler.NewsDetailsPageHandler) {
	// Create news page routes group
	news := router.Group("/news")

	// News page operations
	news.Get("/", handler.HandleNewsDetailsPage) // GET /main/news - Render discover page
}
