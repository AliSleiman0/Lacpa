package handler

import (
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo repository.Repository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// SetupUserRoutes configures user API routes
func (h *UserHandler) SetupUserRoutes(router fiber.Router) {
	// User routes
	users := router.Group("/users")
	users.Post("/", h.CreateUser)
	users.Get("/", h.GetAllUsers)
	users.Get("/:id", h.GetUser)
	users.Get("/email/:email", h.GetUserByEmail)
	users.Get("/username/:username", h.GetUserByUsername)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
	users.Put("/:id/activate", h.ActivateUser)
	users.Put("/:id/deactivate", h.DeactivateUser)
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true // Default to active

	if err := h.repo.CreateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser handles GET /api/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.repo.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// GetUserByEmail handles GET /api/users/email/:email
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	user, err := h.repo.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// GetUserByUsername handles GET /api/users/username/:username
func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.repo.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// GetAllUsers handles GET /api/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.repo.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// UpdateUser handles PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user.UpdatedAt = time.Now()

	if err := h.repo.UpdateUser(c.Context(), id, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.repo.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ActivateUser handles PUT /api/users/:id/activate
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.repo.ActivateUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User activated successfully",
	})
}

// DeactivateUser handles PUT /api/users/:id/deactivate
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.repo.DeactivateUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deactivate user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}
