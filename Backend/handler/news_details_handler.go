package handler

import (
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

type NewsDetailsPageHandler struct {
	repo repository.NewsDetailsPageRepository
}

func NewNewsDetailsPageHandler(repo repository.NewsDetailsPageRepository) *NewsDetailsPageHandler {
	return &NewsDetailsPageHandler{
		repo: repo,
	}
}

func (h *NewsDetailsPageHandler) HandleNewsDetailsPage(c *fiber.Ctx) error {

	data := fiber.Map{
		"Title":   "News Details",
		"BaseURL": "http://" + c.Get("Host"),
	}

	return c.Render("LACPA/sub_pages/news_details_page", data)
}
