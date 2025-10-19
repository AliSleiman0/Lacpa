package handler

import (
	"math"
	"strconv"

	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

type EventsHandler struct {
	repo repository.Repository
}

func NewEventsHandler(repo repository.Repository) *EventsHandler {
	return &EventsHandler{repo: repo}
}

// GetEventsPage renders the events HTML page with pagination
func (h *EventsHandler) GetEventsPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Get page and pageSize from query params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "8"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 8 // Default to 8 events per page
	}

	// Get total count of all events
	totalCount, err := h.repo.CountEventsByCategory(c.Context(), nil)
	if err != nil {
		totalCount = 0
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Ensure page doesn't exceed total pages
	if page > totalPages {
		page = totalPages
	}

	// Get events grouped by category with pagination
	eventsGrouped, err := h.repo.GetEventsGroupedByCategoryPaginated(c.Context(), page, pageSize)
	if err != nil {
		// If error, render with empty data
		return c.Render("LACPA/events/index", fiber.Map{
			"Title":       "Events & News",
			"Events":      nil,
			"CurrentPage": 1,
			"TotalPages":  1,
			"PageSize":    pageSize,
			"TotalCount":  0,
		})
	}

	// Render the events page template with grouped events and pagination info
	return c.Render("LACPA/events/index", fiber.Map{
		"Title":       "Events & News",
		"Events":      eventsGrouped,
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"PageSize":    pageSize,
		"TotalCount":  totalCount,
	})
}
