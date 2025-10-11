package handler

import (
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests
type Handler struct {
	repo repository.Repository
}

// NewHandler creates a new handler
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// SetupRoutes configures all API routes
func (h *Handler) SetupRoutes(router fiber.Router) {
	// Item routes
	items := router.Group("/items")
	items.Post("/", h.CreateItem)
	items.Get("/", h.GetAllItems)
	items.Get("/:id", h.GetItem)
	items.Put("/:id", h.UpdateItem)
	items.Delete("/:id", h.DeleteItem)
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
