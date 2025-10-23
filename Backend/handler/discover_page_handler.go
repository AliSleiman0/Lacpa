package handler

import (
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

type DiscoverPageHandler struct {
	repo repository.DiscoverPageRepository
}

func NewDiscoverPageHandler(repo repository.DiscoverPageRepository) *DiscoverPageHandler {
	return &DiscoverPageHandler{
		repo: repo,
	}
}

func (h *DiscoverPageHandler) HandleDiscoverPage(c *fiber.Ctx) error {

	data := fiber.Map{
		"Title":   "Discover Page",
		"BaseURL": "http://" + c.Get("Host"),
	}

	return c.Render("LACPA/discover/index", data)
}
