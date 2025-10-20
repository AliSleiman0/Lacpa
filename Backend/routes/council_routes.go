package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupCouncilRoutes configures all council-related routes
func SetupCouncilRoutes(app *fiber.App, repo repository.Repository) {
	councilHandler := handler.NewCouncilHandler(repo)

	// Council routes
	api := app.Group("/api/council")

	// Council CRUD operations
	api.Get("/", councilHandler.GetAllCouncils)          // Get all councils
	api.Get("/active", councilHandler.GetActiveCouncil)  // Get active council
	api.Get("/:id", councilHandler.GetCouncilByID)       // Get council by ID
	api.Post("/", councilHandler.CreateCouncil)          // Create new council
	api.Put("/:id", councilHandler.UpdateCouncil)        // Update council
	api.Delete("/:id", councilHandler.DeactivateCouncil) // Deactivate council

	// Council composition routes
	api.Get("/:id/composition", councilHandler.GetCouncilComposition)                    // Get composition (IDs only)
	api.Get("/:id/composition/details", councilHandler.GetCouncilCompositionWithDetails) // Get composition with member details
	api.Get("/:id/positions/available", councilHandler.GetAvailablePositions)            // Get available position slots
	api.Get("/:id/positions/validate", councilHandler.ValidatePositionAvailability)      // Validate position availability

	// Position assignment routes
	api.Post("/position", councilHandler.AssignCouncilPosition)               // Assign member to position
	api.Get("/position/:positionId", councilHandler.GetPositionByID)          // Get position by ID
	api.Put("/position/:positionId", councilHandler.UpdateCouncilPosition)    // Update position
	api.Delete("/position/:positionId", councilHandler.RemoveCouncilPosition) // Remove member from position

	// Member-specific routes
	api.Get("/member/:id/history", councilHandler.GetMemberCouncilHistory) // Get member's council history
}
