package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up all authentication routes
func SetupAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler) {
	// Create auth group
	auth := app.Group("/api/auth")

	// Public routes (no authentication required)
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/resend-otp", authHandler.ResendOTP)
	auth.Post("/reset-password", authHandler.ResetPassword)

	// Protected routes (authentication required)
	auth.Get("/profile", middleware.AuthMiddleware, authHandler.GetProfile)
	auth.Post("/logout", middleware.AuthMiddleware, authHandler.Logout)
}
