package handler

import (
	"strconv"

	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

type MembersHandler struct {
	repo repository.Repository
}

func NewMembersHandler(repo repository.Repository) *MembersHandler {
	return &MembersHandler{repo: repo}
}

// GetIndividualsPage renders the individual members HTML page
func (h *MembersHandler) GetIndividualsPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Get pagination parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "12"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 12 // Default to 12 members per page
	}

	// Get member type filter (optional)
	memberType := c.Query("type", "all") // "all", "Apprentices", "Practicing", "Non-Practicing", "Retired"

	// Get members from repository
	var members interface{}
	var total int64

	if memberType == "all" || memberType == "" {
		members, total, err = h.repo.GetAllIndividualMembers(c.Context(), page, pageSize)
	} else {
		members, total, err = h.repo.GetIndividualMembersByType(c.Context(), memberType, page, pageSize)
	}

	if err != nil {
		// If error, render with empty data
		return c.Render("LACPA/members/individuals", fiber.Map{
			"Title":           "Individual Members",
			"Members":         nil,
			"CurrentPage":     1,
			"TotalPages":      1,
			"PageSize":        pageSize,
			"TotalCount":      0,
			"MemberType":      memberType,
			"ShowIndividuals": true,
			"ShowFirms":       false,
		})
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	// Render the members page template
	return c.Render("LACPA/members/individuals", fiber.Map{
		"Title":           "Individual Members",
		"Members":         members,
		"CurrentPage":     page,
		"TotalPages":      totalPages,
		"PageSize":        pageSize,
		"TotalCount":      total,
		"MemberType":      memberType,
		"ShowIndividuals": true,
		"ShowFirms":       false,
	})
}

// GetFirmsPage renders the firm members HTML page
func (h *MembersHandler) GetFirmsPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Get pagination parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "12"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 12 // Default to 12 firms per page
	}

	// Get firm type/size filter (optional)
	firmType := c.Query("type", "all") // "all", "Audit Firm", "Accounting Firm", "Consultancy"
	firmSize := c.Query("size", "all") // "all", "Big 4", "Large", "Medium", "Small"

	// Get firms from repository
	var firms interface{}
	var total int64

	// Priority: size filter > type filter > all
	if firmSize != "all" && firmSize != "" {
		firms, total, err = h.repo.GetFirmMembersBySize(c.Context(), firmSize, page, pageSize)
	} else if firmType != "all" && firmType != "" {
		firms, total, err = h.repo.GetFirmMembersByType(c.Context(), firmType, page, pageSize)
	} else {
		firms, total, err = h.repo.GetAllFirmMembers(c.Context(), page, pageSize)
	}

	if err != nil {
		// If error, render with empty data
		return c.Render("LACPA/members/firms", fiber.Map{
			"Title":           "Firm Members",
			"Firms":           nil,
			"CurrentPage":     1,
			"TotalPages":      1,
			"PageSize":        pageSize,
			"TotalCount":      0,
			"FirmType":        firmType,
			"FirmSize":        firmSize,
			"ShowIndividuals": false,
			"ShowFirms":       true,
		})
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	// Render the firms page template
	return c.Render("LACPA/members/firms", fiber.Map{
		"Title":           "Firm Members",
		"Firms":           firms,
		"CurrentPage":     page,
		"TotalPages":      totalPages,
		"PageSize":        pageSize,
		"TotalCount":      total,
		"FirmType":        firmType,
		"FirmSize":        firmSize,
		"ShowIndividuals": false,
		"ShowFirms":       true,
	})
}
