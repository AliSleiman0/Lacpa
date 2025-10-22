package admin

import (
	"context"
	"fmt"
	"html/template"
	"time"

	adminModel "github.com/AliSleiman0/Lacpa/models/admin"
	adminRepo "github.com/AliSleiman0/Lacpa/repository/admin"

	"github.com/gofiber/fiber/v2"
)

type AdminHeroSlideHandler struct {
	repo *adminRepo.HeroSlideRepository
}

func NewAdminHeroSlideHandler(repo *adminRepo.HeroSlideRepository) *AdminHeroSlideHandler {
	return &AdminHeroSlideHandler{
		repo: repo,
	}
}

// GetAllSlides handles GET /api/admin/slides
func (h *AdminHeroSlideHandler) GetAllSlides(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slides, err := h.repo.GetAllSlides(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch slides",
		})
	}

	return c.JSON(slides)
}

// GetSlideTabs handles GET /api/admin/slides/tabs
// Returns HTML for slide tabs (for HTMX)
func (h *AdminHeroSlideHandler) GetSlideTabs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slides, err := h.repo.GetAllSlides(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			`<div class="text-red-500">Failed to load slides</div>`)
	}

	// Build HTML for tabs
	html := ""
	for i, slide := range slides {
		activeClass := ""
		if i == 0 {
			activeClass = "text-white border-b-2 border-blue-500"
		} else {
			activeClass = "text-gray-500"
		}

		html += `<button 
			class="slide-tab px-4 py-2 text-sm font-medium hover:text-white transition-colors ` + activeClass + `"
			hx-get="/api/admin/slides/` + slide.ID.Hex() + `/render"
			hx-target="#section-content"
			hx-swap="innerHTML">
			Slide ` + fmt.Sprintf("%d", i+1) + `
		</button>`
	}

	// Add the "Add Slide" button
	html += `
		<button id="add-slide-btn" class="ml-auto px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-2">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			Add Slide
		</button>`

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}

// GetSlideByID handles GET /api/admin/slides/:id
func (h *AdminHeroSlideHandler) GetSlideByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	slide, err := h.repo.GetSlideByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Slide not found",
		})
	}

	return c.JSON(slide)
}

// RenderSlide handles GET /api/admin/slides/:id/render
// Returns HTML fragment for HTMX
func (h *AdminHeroSlideHandler) RenderSlide(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")

	// Handle special case for "first" - get the first slide
	var slide *adminModel.HeroSlide
	var err error

	if id == "first" {
		slides, err := h.repo.GetAllSlides(ctx)
		if err != nil || len(slides) == 0 {
			return c.Status(fiber.StatusNotFound).SendString(
				`<div class="text-center text-gray-400 py-12">
					<i class="fas fa-inbox text-4xl mb-4"></i>
					<p>No slides found. Click "Add Slide" to create one.</p>
				</div>`)
		}
		slide = slides[0]
	} else {
		slide, err = h.repo.GetSlideByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("<div class='text-red-500'>Slide not found</div>")
		}
	}

	// Parse template
	tmpl, err := template.ParseFiles("templates/Admin_Dashboard/hero_section/slide.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			"<div class='text-red-500'>Template error: " + err.Error() + "</div>")
	}

	// Set content type
	c.Set("Content-Type", "text/html; charset=utf-8")

	// Execute template
	return tmpl.Execute(c.Response().BodyWriter(), slide)
}

// CreateSlide handles POST /api/admin/slides
func (h *AdminHeroSlideHandler) CreateSlide(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req adminModel.CreateSlideRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get the next order index
	count, err := h.repo.GetSlideCount(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get slide count",
		})
	}

	slide := &adminModel.HeroSlide{
		Title:             req.Title,
		Description:       req.Description,
		ImgSrc:            req.ImgSrc,
		ButtonTitle:       req.ButtonTitle,
		ButtonLink:        req.ButtonLink,
		IsActive:          req.IsActive,
		ImageActive:       req.ImageActive,
		ButtonActive:      req.ButtonActive,
		TitleActive:       req.TitleActive,
		DescriptionActive: req.DescriptionActive,
		OrderIndex:        int(count) + 1,
	}

	if err := h.repo.CreateSlide(ctx, slide); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create slide",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(slide)
}

// UpdateSlide handles PATCH /api/admin/slides/:id
func (h *AdminHeroSlideHandler) UpdateSlide(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")

	// Check if slide exists
	existingSlide, err := h.repo.GetSlideByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Slide not found",
		})
	}

	var req adminModel.UpdateSlideRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update only provided fields
	if req.Title != "" {
		existingSlide.Title = req.Title
	}
	if req.Description != "" {
		existingSlide.Description = req.Description
	}
	if req.ImgSrc != "" {
		existingSlide.ImgSrc = req.ImgSrc
	}
	if req.ButtonTitle != "" {
		existingSlide.ButtonTitle = req.ButtonTitle
	}
	if req.ButtonLink != "" {
		existingSlide.ButtonLink = req.ButtonLink
	}
	if req.IsActive != nil {
		existingSlide.IsActive = *req.IsActive
	}
	if req.ImageActive != nil {
		existingSlide.ImageActive = *req.ImageActive
	}
	if req.ButtonActive != nil {
		existingSlide.ButtonActive = *req.ButtonActive
	}
	if req.TitleActive != nil {
		existingSlide.TitleActive = *req.TitleActive
	}
	if req.DescriptionActive != nil {
		existingSlide.DescriptionActive = *req.DescriptionActive
	}
	if req.OrderIndex != nil {
		existingSlide.OrderIndex = *req.OrderIndex
	}

	if err := h.repo.UpdateSlide(ctx, id, existingSlide); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update slide",
		})
	}

	return c.JSON(existingSlide)
}

// DeleteSlide handles DELETE /api/admin/slides/:id
func (h *AdminHeroSlideHandler) DeleteSlide(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")

	if err := h.repo.DeleteSlide(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete slide",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
