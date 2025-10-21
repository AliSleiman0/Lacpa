package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/AliSleiman0/Lacpa/utils"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	authRepo *repository.AuthRepository
}

func NewAdminHandler(authRepo *repository.AuthRepository) *AdminHandler {
	return &AdminHandler{
		authRepo: authRepo,
	}
}

// CreateAdmin allows admins to create other admin accounts
func (h *AdminHandler) CreateAdmin(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
			"success": false,
		})
	}

	// Get requesting user's role from context
	requestingUserRole := c.Locals("role").(string)
	if requestingUserRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Only admins can create admin accounts",
			"success": false,
		})
	}

	var req models.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body. Please check your JSON format.",
			"success": false,
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Check if email already exists
	existingUser, err := h.authRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "Email already registered",
			"success": false,
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process password",
			"success": false,
		})
	}

	// Generate LACPA ID (format: LACPA-YYYY-XXXXX)
	lacpaID := fmt.Sprintf("LACPA-%d-%05d", time.Now().Year(), time.Now().Unix()%100000)

	// Create admin user
	user := &models.User{
		LACPAID:    lacpaID,
		FullName:   req.FullName,
		Email:      strings.ToLower(strings.TrimSpace(req.Email)),
		Password:   hashedPassword,
		Role:       "admin", // Set as admin
		IsVerified: true,    // Auto-verify admin accounts
		IsActive:   true,
	}

	if err := h.authRepo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create admin user",
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Admin account created successfully",
		"data": fiber.Map{
			"email":    user.Email,
			"lacpa_id": user.LACPAID,
			"role":     user.Role,
		},
	})
}

// UpdateUserRole allows admins to change user roles
func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
			"success": false,
		})
	}

	type UpdateRoleRequest struct {
		Email string `json:"email" validate:"required,email"`
		Role  string `json:"role" validate:"required"`
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body. Please check your JSON format.",
			"success": false,
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Validate role value
	validRoles := map[string]bool{
		"admin":  true,
		"member": true,
		"guest":  true,
	}
	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid role. Allowed roles: admin, member, guest",
			"success": false,
		})
	}

	// Get user
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"success": false,
		})
	}

	// Update role
	user.Role = req.Role
	if err := h.authRepo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update user role",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("User role updated to %s", req.Role),
		"data": fiber.Map{
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// ListUsers allows admins to list all users
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	// This would need a new repository method
	// For now, return a simple response
	return c.JSON(fiber.Map{
		"success": true,
		"message": "List users endpoint - to be implemented",
	})
}

// DeactivateUser allows admins to deactivate user accounts
func (h *AdminHandler) DeactivateUser(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Email parameter is required",
			"success": false,
		})
	}

	// Get user
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"success": false,
		})
	}

	// Deactivate
	user.IsActive = false
	if err := h.authRepo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to deactivate user",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User account deactivated successfully",
	})
}

// ActivateUser allows admins to activate user accounts
func (h *AdminHandler) ActivateUser(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Email parameter is required",
			"success": false,
		})
	}

	// Get user
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"success": false,
		})
	}

	// Activate
	user.IsActive = true
	if err := h.authRepo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to activate user",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User account activated successfully",
	})
}
