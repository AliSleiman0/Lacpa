package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupMembersRoutes configures all members page routes
func SetupMembersRoutes(app *fiber.App, repo repository.Repository) {
	membersHandler := handler.NewMembersHandler(repo)

	// Members page routes
	app.Get("/members/individuals", membersHandler.GetIndividualsPage)
}
