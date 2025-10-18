package handler

import (
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/AliSleiman0/Lacpa/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouncilHandler struct {
	repo repository.Repository
}

func NewCouncilHandler(repo repository.Repository) *CouncilHandler {
	return &CouncilHandler{repo: repo}
}

// GetActiveCouncil retrieves the currently active council
func (h *CouncilHandler) GetActiveCouncil(c *fiber.Ctx) error {
	council, err := h.repo.GetActiveCouncil(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, "Active council retrieved successfully", council)
}

// GetCouncilByID retrieves a council by ID
func (h *CouncilHandler) GetCouncilByID(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	council, err := h.repo.GetCouncilByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, "Council retrieved successfully", council)
}

// GetAllCouncils retrieves all councils
func (h *CouncilHandler) GetAllCouncils(c *fiber.Ctx) error {
	councils, err := h.repo.GetAllCouncils(c.Context())
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Councils retrieved successfully", councils)
}

// CreateCouncil creates a new council
func (h *CouncilHandler) CreateCouncil(c *fiber.Ctx) error {
	var council models.Council
	if err := c.BodyParser(&council); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if council.Name == "" {
		return utils.SendBadRequest(c, "Council name is required")
	}
	if council.StartDate.IsZero() {
		return utils.SendBadRequest(c, "Start date is required")
	}

	if err := h.repo.CreateCouncil(c.Context(), &council); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendCreated(c, council, "", "")
}

// UpdateCouncil updates an existing council
func (h *CouncilHandler) UpdateCouncil(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	var council models.Council
	if err := c.BodyParser(&council); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	if err := h.repo.UpdateCouncil(c.Context(), id, &council); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council updated successfully", council)
}

// DeactivateCouncil deactivates a council
func (h *CouncilHandler) DeactivateCouncil(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	if err := h.repo.DeactivateCouncil(c.Context(), id); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council deactivated successfully", nil)
}

// GetCouncilComposition retrieves council composition
func (h *CouncilHandler) GetCouncilComposition(c *fiber.Ctx) error {
	councilID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	composition, err := h.repo.GetCouncilComposition(c.Context(), councilID)
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council composition retrieved successfully", composition)
}

// GetCouncilCompositionWithDetails retrieves council composition with member details
func (h *CouncilHandler) GetCouncilCompositionWithDetails(c *fiber.Ctx) error {
	councilID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	composition, err := h.repo.GetCouncilCompositionWithDetails(c.Context(), councilID)
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council composition with details retrieved successfully", composition)
}

// AssignCouncilPosition assigns a member to a council position
func (h *CouncilHandler) AssignCouncilPosition(c *fiber.Ctx) error {
	var position models.CouncilPosition
	if err := c.BodyParser(&position); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if position.MemberID.IsZero() {
		return utils.SendBadRequest(c, "Member ID is required")
	}
	if position.CouncilID.IsZero() {
		return utils.SendBadRequest(c, "Council ID is required")
	}
	if !position.Position.IsCouncilPosition() {
		return utils.SendBadRequest(c, "Invalid council position type")
	}
	if position.StartDate.IsZero() {
		position.StartDate = time.Now()
	}

	if err := h.repo.AssignCouncilPosition(c.Context(), &position); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendCreated(c, position, "", "")
}

// RemoveCouncilPosition removes a member from a council position
func (h *CouncilHandler) RemoveCouncilPosition(c *fiber.Ctx) error {
	positionID, err := primitive.ObjectIDFromHex(c.Params("positionId"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid position ID")
	}

	if err := h.repo.RemoveCouncilPosition(c.Context(), positionID); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council position removed successfully", nil)
}

// UpdateCouncilPosition updates a council position
func (h *CouncilHandler) UpdateCouncilPosition(c *fiber.Ctx) error {
	positionID, err := primitive.ObjectIDFromHex(c.Params("positionId"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid position ID")
	}

	var position models.CouncilPosition
	if err := c.BodyParser(&position); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	if err := h.repo.UpdateCouncilPosition(c.Context(), positionID, &position); err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Council position updated successfully", position)
}

// GetMemberCouncilHistory retrieves council history for a member
func (h *CouncilHandler) GetMemberCouncilHistory(c *fiber.Ctx) error {
	memberID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid member ID")
	}

	positions, err := h.repo.GetMemberCouncilHistory(c.Context(), memberID)
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Member council history retrieved successfully", positions)
}

// GetPositionByID retrieves a council position by ID
func (h *CouncilHandler) GetPositionByID(c *fiber.Ctx) error {
	positionID, err := primitive.ObjectIDFromHex(c.Params("positionId"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid position ID")
	}

	position, err := h.repo.GetPositionByID(c.Context(), positionID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, "Position retrieved successfully", position)
}

// ValidatePositionAvailability checks if a position slot is available
func (h *CouncilHandler) ValidatePositionAvailability(c *fiber.Ctx) error {
	councilID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	positionType := models.CouncilPositionType(c.Query("type"))
	if !positionType.IsCouncilPosition() {
		return utils.SendBadRequest(c, "Invalid position type")
	}

	available, err := h.repo.ValidatePositionAvailability(c.Context(), councilID, positionType)
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Position availability validated", fiber.Map{
		"available": available,
		"position":  positionType,
	})
}

// GetAvailablePositions returns available slots for all position types
func (h *CouncilHandler) GetAvailablePositions(c *fiber.Ctx) error {
	councilID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.SendBadRequest(c, "Invalid council ID")
	}

	available, err := h.repo.GetAvailablePositions(c.Context(), councilID)
	if err != nil {
		return utils.SendInternalError(c, err.Error())
	}

	return utils.SendSuccess(c, "Available positions retrieved successfully", available)
}
