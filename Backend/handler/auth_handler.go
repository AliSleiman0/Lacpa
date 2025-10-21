package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/AliSleiman0/Lacpa/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	authRepo *repository.AuthRepository
}

func NewAuthHandler(authRepo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
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

	// Create user
	user := &models.User{
		LACPAID:  lacpaID,
		FullName: req.FullName,
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
		Role:     "member", // Default role
	}

	if err := h.authRepo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create user",
			"success": false,
		})
	}

	// Generate OTP for email verification
	otp, err := utils.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate OTP",
			"success": false,
		})
	}

	// Set OTP with 10 minutes expiry
	otpExpiry := time.Now().Add(10 * time.Minute)
	if err := h.authRepo.SetOTP(user.Email, otp, otpExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to set OTP",
			"success": false,
		})
	}

	// Send OTP via email with beautiful HTML template (use the SAME OTP that's stored in DB)
	_, err = utils.SendOTPEmailWithCode(user.Email, user.FullName, otp)
	if err != nil {
		// Log error but don't fail registration - user is created and OTP is stored in DB
		fmt.Printf("Warning: Failed to send OTP email to %s: %v\n", user.Email, err)
		fmt.Printf("OTP for %s (email failed): %s\n", user.Email, otp)
	} else {
		fmt.Printf("✅ OTP email sent successfully to %s\n", user.Email)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Registration successful. Please verify your email with the OTP sent.",
		"data": fiber.Map{
			"email":    user.Email,
			"lacpa_id": user.LACPAID,
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
			"success": false,
		})
	}

	var req models.LoginRequest
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

	// Get user by LACPA ID
	user, err := h.authRepo.GetUserByLACPAID(req.LACPAID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid credentials",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve user",
			"success": false,
		})
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"success": false,
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Account is deactivated",
			"success": false,
		})
	}

	// Check if user is verified
	if !user.IsVerified {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Please verify your email first",
			"success": false,
		})
	}

	// Update last login
	if err := h.authRepo.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.LACPAID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate token",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": models.AuthResponse{
			Token: token,
			User:  user.ToResponse(),
		},
	})
}

// ForgotPassword initiates password reset process
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
			"success": false,
		})
	}

	var req models.ForgotPasswordRequest
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

	// Get user by email
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		// Don't reveal if email exists or not (security best practice)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "If the email exists, a reset OTP has been sent.",
		})
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate OTP",
			"success": false,
		})
	}

	// Set OTP with 10 minutes expiry
	otpExpiry := time.Now().Add(10 * time.Minute)
	if err := h.authRepo.SetOTP(user.Email, otp, otpExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to set OTP",
			"success": false,
		})
	}

	// Send password reset OTP via email with beautiful HTML template (use the SAME OTP that's stored in DB)
	_, err = utils.SendOTPEmailWithCode(user.Email, user.FullName, otp)
	if err != nil {
		fmt.Printf("Warning: Failed to send password reset OTP to %s: %v\n", user.Email, err)
		fmt.Printf("Password Reset OTP for %s (email failed): %s\n", user.Email, otp)
	} else {
		fmt.Printf("✅ Password reset OTP email sent successfully to %s\n", user.Email)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a reset OTP has been sent.",
	})
}

// VerifyOTP verifies the OTP and generates reset token
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	// Check Content-Type
	contentType := c.Get("Content-Type")
	if contentType != "application/json" && !strings.Contains(contentType, "application/json") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Content-Type must be application/json",
			"success": false,
		})
	}

	var req models.VerifyOTPRequest
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

	// Get user by email
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid OTP",
			"success": false,
		})
	}

	// Check if OTP exists and not expired
	if user.OTP == "" || time.Now().After(user.OTPExpiry) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "OTP expired or invalid",
			"success": false,
		})
	}

	// Verify OTP
	if user.OTP != req.OTP {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid OTP",
			"success": false,
		})
	}

	// Clear OTP
	if err := h.authRepo.ClearOTP(user.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to clear OTP",
			"success": false,
		})
	}

	// Verify user if not already verified
	if !user.IsVerified {
		if err := h.authRepo.VerifyUser(user.Email); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to verify user",
				"success": false,
			})
		}
	}

	// Generate reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate reset token",
			"success": false,
		})
	}

	// Set reset token with 15 minutes expiry
	resetTokenExpiry := time.Now().Add(15 * time.Minute)
	if err := h.authRepo.SetResetToken(user.Email, resetToken, resetTokenExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to set reset token",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "OTP verified successfully",
		"data": fiber.Map{
			"reset_token": resetToken,
		},
	})
}

// ResendOTP resends OTP to user
func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	var req models.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
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

	// Get user by email
	user, err := h.authRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		// Don't reveal if email exists or not
		return c.JSON(fiber.Map{
			"success": true,
			"message": "If the email exists, a new OTP has been sent.",
		})
	}

	// Generate new OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate OTP",
			"success": false,
		})
	}

	// Set OTP with 10 minutes expiry
	otpExpiry := time.Now().Add(10 * time.Minute)
	if err := h.authRepo.SetOTP(user.Email, otp, otpExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to set OTP",
			"success": false,
		})
	}

	// Resend OTP via email with beautiful HTML template (use the SAME OTP that's stored in DB)
	_, err = utils.SendOTPEmailWithCode(user.Email, user.FullName, otp)
	if err != nil {
		fmt.Printf("Warning: Failed to resend OTP to %s: %v\n", user.Email, err)
		fmt.Printf("Resend OTP for %s (email failed): %s\n", user.Email, otp)
	} else {
		fmt.Printf("✅ OTP resent successfully to %s\n", user.Email)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a new OTP has been sent.",
	})
}

// ResetPassword resets user password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
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

	// Get user by reset token
	user, err := h.authRepo.GetUserByResetToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid or expired reset token",
			"success": false,
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process password",
			"success": false,
		})
	}

	// Update password
	if err := h.authRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update password",
			"success": false,
		})
	}

	// Clear reset token
	if err := h.authRepo.ClearResetToken(user.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to clear reset token: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password reset successful. You can now login with your new password.",
	})
}

// GetProfile returns current user profile (requires authentication)
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user ID",
			"success": false,
		})
	}

	// Get user
	user, err := h.authRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user.ToResponse(),
	})
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a stateless JWT system, logout is handled client-side by removing the token
	// If you want to implement token blacklisting, you can add it here
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logout successful",
	})
}
