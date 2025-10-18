package handler

import (
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

type MembersHandler struct {
	repo repository.Repository
}

func NewMembersHandler(repo repository.Repository) *MembersHandler {
	return &MembersHandler{repo: repo}
}

// GetIndividualsPage renders the members HTML page
func (h *MembersHandler) GetIndividualsPage(c *fiber.Ctx) error {
	// For now, just render the template
	// In the future, you can fetch member data and pass it to the template
	return c.Render("LACPA/members/individuals", fiber.Map{
		"Title": "Individual Members",
	})
}
