package handler

import (
	"sync"

	"github.com/AliSleiman0/Lacpa/utils"
	"github.com/gofiber/fiber/v2"
)

// In-memory OTP storage (for testing - use Redis or database in production)
var (
	otpStore = make(map[string]*utils.OTPEmail)
	otpMutex sync.RWMutex
)

// OTPHandler handles OTP-related requests
type OTPHandler struct{}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler() *OTPHandler {
	return &OTPHandler{}
}

// SendOTP sends an OTP to the specified email
// POST /api/otp/send
// Body: { "email": "user@example.com" }
func (h *OTPHandler) SendOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	if req.Email == "" {
		return utils.SendBadRequest(c, "Email is required")
	}

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return utils.SendBadRequest(c, "Invalid email format")
	}

	// Send OTP email
	otpData, err := utils.SendOTPEmail(req.Email)
	if err != nil {
		return utils.SendInternalError(c, "Failed to send OTP email")
	}

	// Store OTP in memory (use Redis or database in production)
	otpMutex.Lock()
	otpStore[req.Email] = otpData
	otpMutex.Unlock()

	return utils.SendSuccess(c, "OTP sent successfully", fiber.Map{
		"email":      otpData.Email,
		"expires_at": otpData.ExpiresAt,
	})
}

// VerifyOTP verifies the OTP provided by the user
// POST /api/otp/verify
// Body: { "email": "user@example.com", "otp": "123456" }
func (h *OTPHandler) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body")
	}

	if req.Email == "" || req.OTP == "" {
		return utils.SendBadRequest(c, "Email and OTP are required")
	}

	// Get stored OTP
	otpMutex.RLock()
	storedOTP, exists := otpStore[req.Email]
	otpMutex.RUnlock()

	if !exists {
		return utils.SendError(c, fiber.StatusNotFound, "No OTP found for this email")
	}

	// Verify OTP
	if !utils.VerifyOTP(storedOTP, req.OTP) {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid or expired OTP")
	}

	// Remove OTP from store after successful verification
	otpMutex.Lock()
	delete(otpStore, req.Email)
	otpMutex.Unlock()

	return utils.SendSuccess(c, "OTP verified successfully", fiber.Map{
		"email":    req.Email,
		"verified": true,
	})
}
