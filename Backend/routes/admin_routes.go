package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	adminHandler "github.com/AliSleiman0/Lacpa/handler/admin"
	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes sets up all admin-only routes
func SetupAdminRoutes(app *fiber.App, adminUserHandler *handler.AdminHandler, heroSlideHandler *adminHandler.AdminHeroSlideHandler) {
	// Create admin group with authentication and admin role required
	admin := app.Group("/api/admin")
	//admin.Use(middleware.AuthMiddleware)
	//	admin.Use(middleware.RoleMiddleware("admin"))

	// Admin user management
	admin.Post("/create-admin", adminUserHandler.CreateAdmin)
	admin.Post("/update-role", adminUserHandler.UpdateUserRole)
	admin.Get("/users", adminUserHandler.ListUsers)
	admin.Post("/deactivate-user", adminUserHandler.DeactivateUser)
	admin.Post("/activate-user", adminUserHandler.ActivateUser)

	// Hero Slides Management
	admin.Get("/slides", heroSlideHandler.GetAllSlides)
	admin.Get("/slides/tabs", heroSlideHandler.GetSlideTabs) // Returns HTML for slide tabs
	admin.Get("/slides/:id", heroSlideHandler.GetSlideByID)
	admin.Get("/slides/:id/render", heroSlideHandler.RenderSlide) // Returns HTML fragment
	admin.Post("/slides", heroSlideHandler.CreateSlide)
	admin.Patch("/slides/:id", heroSlideHandler.UpdateSlide)
	admin.Delete("/slides/:id", heroSlideHandler.DeleteSlide)
}
