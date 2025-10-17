package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupMainPageRoutes configures main page routes
//
// ROLE: Main Page Route Configuration
// - Defines routes for main LACPA landing page
// - Maps HTTP methods and paths to handler functions
// - Groups related main page endpoints
// - Separates routing logic from business logic
//
// ROUTES:
//   - GET    /main/landing     - Get main landing page
//   - POST   /main/landing     - Handle landing page form submissions
//
// PARAMETERS:
//   - router: Fiber router group (typically /api or root)
//   - h: Handler instance containing main page business logic
func SetupMainPageRoutes(router fiber.Router, h *handler.MainPageHandler) {
	// Create main page routes group
	main := router.Group("main/landing")

	// Main page operations
	main.Get("/", h.RenderPage) // GET /main/landing - Render landing page

}
