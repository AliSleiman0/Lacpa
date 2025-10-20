package routes

import (
	"github.com/AliSleiman0/Lacpa/handler"
	"github.com/AliSleiman0/Lacpa/repository"
	"github.com/gofiber/fiber/v2"
)

// SetupApplicationRoutes configures all application-related routes
func SetupApplicationRoutes(app *fiber.App, repo repository.Repository) {
	// Create application handler (repo already includes ApplicationRepository)
	appHandler := handler.NewApplicationHandler(repo)

	// Public routes - viewing requirements and submitting applications
	app.Get("/membership/apply-now", appHandler.GetApplyNowPage)
	app.Get("/membership/apply/firm", appHandler.GetApplyFirmPage)
	app.Get("/membership/apply/individual", appHandler.GetApplyIndividualPage)

	// API routes for application submission
	api := app.Group("/api/applications")
	api.Post("/individual", appHandler.SubmitIndividualApplication)
	api.Post("/firm", appHandler.SubmitFirmApplication)

	// Admin routes for managing applications (to be protected with auth middleware later)
	api.Get("/individual", appHandler.GetAllIndividualApplications)
	api.Get("/firm", appHandler.GetAllFirmApplications)
	api.Put("/individual/:id/status", appHandler.UpdateIndividualApplicationStatus)
	api.Put("/firm/:id/status", appHandler.UpdateFirmApplicationStatus)
}
