package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupMembersRoutes configures all members page routes
func SetupMembersRoutes(app *fiber.App, repo repository.Repository) {
	membersHandler := handler.NewMembersHandler(repo)
	councilHandler := handler.NewCouncilHandler(repo)

	// Members page routes - support both URL patterns
	app.Get("/members/individuals", membersHandler.GetIndividualsPage)
	app.Get("/membership", membersHandler.GetIndividualsPage)                   // Clean URL alias
	app.Get("/membership/firms", membersHandler.GetFirmsPage)                   // Firms page
	app.Get("/discover/board-of-directors", councilHandler.GetBoardMembersPage) // Board members page
}
