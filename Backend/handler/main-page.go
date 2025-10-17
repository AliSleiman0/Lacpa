package handler

import (
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// MainPageHandler handles HTTP requests for the main landing page
type MainPageHandler struct {
	repo repository.Repository
}

// NewMainPageHandler creates a new main page handler
func NewMainPageHandler(repo repository.Repository) *MainPageHandler {
	return &MainPageHandler{
		repo: repo,
	}
}

// RenderPage handles GET /api/main/landing - renders the main landing page
func (h *MainPageHandler) RenderPage(c *fiber.Ctx) error {
	// Prepare data for the landing page
	data := fiber.Map{
		"Title":   "Main Landing Page",
		"BaseURL": "http://" + c.Get("Host"),
	}

	// Render with layout
	return c.Render("LACPA/main_LACPA/index", data, "layouts/main")
}
