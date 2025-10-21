package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes sets up all admin-only routes
func SetupAdminRoutes(app *fiber.App, adminHandler *handler.AdminHandler) {
	// Create admin group with authentication and admin role required
	admin := app.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.RoleMiddleware("admin"))

	// Admin user management
	admin.Post("/create-admin", adminHandler.CreateAdmin)
	admin.Post("/update-role", adminHandler.UpdateUserRole)
	admin.Get("/users", adminHandler.ListUsers)
	admin.Post("/deactivate-user", adminHandler.DeactivateUser)
	admin.Post("/activate-user", adminHandler.ActivateUser)
}
