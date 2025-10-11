package handler

import (
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for item operations
//
// ROLE: Item Business Logic Handler
// - Contains all business logic for item operations (CRUD)
// - Handles HTTP request/response processing for items
// - Validates request data and formats responses
// - Interfaces with the repository layer for data persistence
// - Focuses purely on business logic, route configuration moved to routes package
type Handler struct {
	repo repository.Repository
}

// NewHandler creates a new item handler instance
//
// PARAMETERS:
//
//	repo: Unified repository interface providing access to all data operations
//
// RETURNS:
//
//	*Handler: Handler instance ready to process item-related HTTP requests
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// CreateItem handles POST /api/items
func (h *Handler) CreateItem(c *fiber.Ctx) error {
	var item models.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := h.repo.CreateItem(c.Context(), &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// GetItem handles GET /api/items/:id
func (h *Handler) GetItem(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.repo.GetItem(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	return c.JSON(item)
}

// GetAllItems handles GET /api/items
func (h *Handler) GetAllItems(c *fiber.Ctx) error {
	items, err := h.repo.GetAllItems(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch items",
		})
	}

	return c.JSON(items)
}

// UpdateItem handles PUT /api/items/:id
func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	item.UpdatedAt = time.Now()

	if err := h.repo.UpdateItem(c.Context(), id, &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item updated successfully",
	})
}

// DeleteItem handles DELETE /api/items/:id
func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.repo.DeleteItem(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item deleted successfully",
	})
}
