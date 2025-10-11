package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures all user-related API routes
//
// ROLE: User Route Configuration
// - Defines all HTTP routes for user operations (CRUD + special operations)
// - Maps HTTP methods and paths to handler functions
// - Groups related user endpoints under /users prefix
// - Includes additional user-specific routes (email, username lookup, activation)
//
// ROUTES:
//
//	POST   /users              - Create a new user
//	GET    /users              - Get all users
//	GET    /users/:id          - Get user by ID
//	GET    /users/email/:email - Get user by email
//	GET    /users/username/:username - Get user by username
//	PUT    /users/:id          - Update user by ID
//	PUT    /users/:id/activate - Activate user
//	PUT    /users/:id/deactivate - Deactivate user
//	DELETE /users/:id          - Delete user by ID
//
// PARAMETERS:
//
//	router: Fiber router group (typically /api)
//	uh: UserHandler instance containing user business logic
func SetupUserRoutes(router fiber.Router, uh *handler.UserHandler) {
	// Create user routes group under /users
	users := router.Group("/users")

	// Basic CRUD operations
	users.Post("/", uh.CreateUser)      // POST /api/users
	users.Get("/", uh.GetAllUsers)      // GET /api/users
	users.Get("/:id", uh.GetUser)       // GET /api/users/:id
	users.Put("/:id", uh.UpdateUser)    // PUT /api/users/:id
	users.Delete("/:id", uh.DeleteUser) // DELETE /api/users/:id

	// Lookup operations
	users.Get("/email/:email", uh.GetUserByEmail)          // GET /api/users/email/:email
	users.Get("/username/:username", uh.GetUserByUsername) // GET /api/users/username/:username

	// User activation operations
	users.Put("/:id/activate", uh.ActivateUser)     // PUT /api/users/:id/activate
	users.Put("/:id/deactivate", uh.DeactivateUser) // PUT /api/users/:id/deactivate
}
