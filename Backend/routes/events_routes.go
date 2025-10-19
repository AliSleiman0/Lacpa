package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupEventsRoutes configures all events page routes
func SetupEventsRoutes(app *fiber.App, repo repository.Repository) {
	eventsHandler := handler.NewEventsHandler(repo)

	// Events page route
	app.Get("/events", eventsHandler.GetEventsPage)
}
