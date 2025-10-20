package handler

import (
	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/AliSleiman0/Lacpa/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationHandler struct {
	repo repository.ApplicationRepository
}

func NewApplicationHandler(repo repository.ApplicationRepository) *ApplicationHandler {
	return &ApplicationHandler{repo: repo}
}

// GetApplyNowPage renders the application page with requirements
func (h *ApplicationHandler) GetApplyNowPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Get all requirements
	requirements, err := h.repo.GetAllRequirements(c.Context())
	if err != nil {
		// If error, render with empty data
		return c.Render("LACPA/membership/apply", fiber.Map{
			"Title":                  "Apply for Membership",
			"IndividualRequirements": []models.ApplicationRequirement{},
			"FirmRequirements":       []models.ApplicationRequirement{},
		})
	}

	// Separate requirements by type
	var individualReqs []models.ApplicationRequirement
	var firmReqs []models.ApplicationRequirement

	for _, req := range requirements {
		if req.ApplicationType == models.ApplicationTypeIndividual {
			individualReqs = append(individualReqs, req)
		} else if req.ApplicationType == models.ApplicationTypeFirm {
			firmReqs = append(firmReqs, req)
		}
	}

	// Render the application page template
	return c.Render("LACPA/membership/apply", fiber.Map{
		"Title":                  "Apply for Membership",
		"IndividualRequirements": individualReqs,
		"FirmRequirements":       firmReqs,
	})
}

// GetApplyFirmPage renders the firm application form page
func (h *ApplicationHandler) GetApplyFirmPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Render the firm application form template
	return c.Render("LACPA/membership/apply_firm", fiber.Map{
		"Title": "Apply to Firm - LACPA",
	})
}

// GetApplyIndividualPage renders the individual application form page
func (h *ApplicationHandler) GetApplyIndividualPage(c *fiber.Ctx) error {
	// Check if this is an HTMX request (fragment) or browser request (need full page)
	if c.Get("HX-Request") != "true" {
		// Browser request - serve index.html and let JavaScript load the content
		return c.SendFile("../LACPA_Web/src/index.html")
	}

	// Render the individual application form template (to be created)
	return c.Render("LACPA/membership/apply_individual", fiber.Map{
		"Title": "Apply as Individual - LACPA",
	})
}

// SubmitIndividualApplication handles individual membership application submission
func (h *ApplicationHandler) SubmitIndividualApplication(c *fiber.Ctx) error {
	var application models.IndividualApplication
	if err := c.BodyParser(&application); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if application.FirstName == "" || application.LastName == "" {
		return utils.SendBadRequest(c, "First name and last name are required")
	}
	if application.Email == "" || application.Phone == "" {
		return utils.SendBadRequest(c, "Email and phone are required")
	}

	if err := h.repo.CreateIndividualApplication(c.Context(), &application); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendCreated(c, application, "Application submitted successfully", "Your application has been received and is under review.")
}

// SubmitFirmApplication handles firm membership application submission
func (h *ApplicationHandler) SubmitFirmApplication(c *fiber.Ctx) error {
	var application models.FirmApplication
	if err := c.BodyParser(&application); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if application.FirmName == "" {
		return utils.SendBadRequest(c, "Firm name is required")
	}
	if application.Email == "" || application.Phone == "" {
		return utils.SendBadRequest(c, "Email and phone are required")
	}
	if application.RepresentativeName == "" {
		return utils.SendBadRequest(c, "Representative name is required")
	}

	if err := h.repo.CreateFirmApplication(c.Context(), &application); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendCreated(c, application, "Application submitted successfully", "Your application has been received and is under review.")
}

// GetAllIndividualApplications retrieves all individual applications
func (h *ApplicationHandler) GetAllIndividualApplications(c *fiber.Ctx) error {
	applications, err := h.repo.GetAllIndividualApplications(c.Context())
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Individual applications retrieved successfully", applications)
}

// GetAllFirmApplications retrieves all firm applications
func (h *ApplicationHandler) GetAllFirmApplications(c *fiber.Ctx) error {
	applications, err := h.repo.GetAllFirmApplications(c.Context())
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Firm applications retrieved successfully", applications)
}

// UpdateIndividualApplicationStatus updates the status of an individual application
func (h *ApplicationHandler) UpdateIndividualApplicationStatus(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid application ID")
	}

	type StatusUpdate struct {
		Status      models.ApplicationStatus `json:"status"`
		ReviewNotes string                   `json:"review_notes"`
		ReviewedBy  string                   `json:"reviewed_by"`
	}

	var update StatusUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	reviewedBy, err := primitive.ObjectIDFromHex(update.ReviewedBy)
	if err != nil {
		return utils.SendBadRequest(c, "Invalid reviewer ID")
	}

	if err := h.repo.UpdateIndividualApplicationStatus(c.Context(), id, update.Status, update.ReviewNotes, reviewedBy); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Application status updated successfully", nil)
}

// UpdateFirmApplicationStatus updates the status of a firm application
func (h *ApplicationHandler) UpdateFirmApplicationStatus(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid application ID")
	}

	type StatusUpdate struct {
		Status      models.ApplicationStatus `json:"status"`
		ReviewNotes string                   `json:"review_notes"`
		ReviewedBy  string                   `json:"reviewed_by"`
	}

	var update StatusUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	reviewedBy, err := primitive.ObjectIDFromHex(update.ReviewedBy)
	if err != nil {
		return utils.SendBadRequest(c, "Invalid reviewer ID")
	}

	if err := h.repo.UpdateFirmApplicationStatus(c.Context(), id, update.Status, update.ReviewNotes, reviewedBy); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Application status updated successfully", nil)
}
