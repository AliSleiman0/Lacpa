package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupAuthPageRoutes sets up routes for serving authentication-related HTML pages
// These are frontend pages (not API endpoints) that can be directly accessed and refreshed
func SetupAuthPageRoutes(app *fiber.App) {
	// Login page
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("../LACPA_Web/src/login.html")
	})

	// Signup page
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.SendFile("../LACPA_Web/src/signup.html")
	})

	// Forgot password page
	app.Get("/forgot-password", func(c *fiber.Ctx) error {
		return c.SendFile("../LACPA_Web/src/pages/login/forgot-password.html")
	})

	// Verify account page (OTP verification)
	app.Get("/verify-account", func(c *fiber.Ctx) error {
		return c.SendFile("../LACPA_Web/src/pages/login/verify-account.html")
	})

	// Reset password page
	app.Get("/reset-password", func(c *fiber.Ctx) error {
		return c.SendFile("../LACPA_Web/src/pages/login/reset-password.html")
	})

	// Dashboard route (protected page - basic version, can be enhanced later)
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		// TODO: You may want to add middleware to check authentication
		// For now, this serves a basic redirect or placeholder
		return c.SendString("Dashboard - Please create dashboard.html in LACPA_Web/src/")
	})
}
