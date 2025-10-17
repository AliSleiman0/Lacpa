package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseData represents structured data for multi-format responses
type ResponseData struct {
	Data     interface{} `json:"data,omitempty"`
	Message  string      `json:"message,omitempty"`
	Success  bool        `json:"success"`
	Error    string      `json:"error,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// SendResponse sends an appropriate response based on request type
//
// ROLE: Unified Response Handler
// - Automatically detects preferred response format
// - Sends JSON, HTMX partial, or full HTML based on request
// - Centralizes response logic for consistent API behavior
//
// PARAMETERS:
//   - c: Fiber context
//   - data: Response data to send
//   - htmlTemplate: Template name for HTML responses
//   - htmxTemplate: Template name for HTMX partial responses (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
//
// USAGE EXAMPLE:
//
//	err := utils.SendResponse(c, user, "user/profile.html", "user/profile-partial.html")
func SendResponse(c *fiber.Ctx, data interface{}, htmlTemplate string, htmxTemplate ...string) error {
	switch GetPreferredContentType(c) {
	case "json":
		return c.JSON(ResponseData{
			Data:    data,
			Success: true,
		})

	case "htmx":
		template := htmlTemplate
		if len(htmxTemplate) > 0 && htmxTemplate[0] != "" {
			template = htmxTemplate[0]
		}
		return c.Render(template, data)

	default: // html
		return c.Render(htmlTemplate, data)
	}
}

// SendError sends an error response in the appropriate format
//
// ROLE: Unified Error Response Handler
// - Sends error responses in JSON, HTMX, or HTML format
// - Maintains consistent error structure across response types
// - Handles HTTP status codes appropriately
//
// PARAMETERS:
//   - c: Fiber context
//   - status: HTTP status code
//   - message: Error message
//   - errorTemplate: Template name for HTML error pages (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
func SendError(c *fiber.Ctx, status int, message string, errorTemplate ...string) error {
	if WantsJSON(c) {
		return c.Status(status).JSON(ResponseData{
			Success: false,
			Error:   message,
		})
	}

	// For HTML/HTMX requests
	template := "error.html"
	if len(errorTemplate) > 0 && errorTemplate[0] != "" {
		template = errorTemplate[0]
	}

	errorData := fiber.Map{
		"error":   message,
		"status":  status,
		"success": false,
	}

	return c.Status(status).Render(template, errorData)
}

// SendSuccess sends a success response with optional data
//
// ROLE: Unified Success Response Handler
// - Sends success responses in appropriate format
// - Handles both data and message-only success responses
// - Maintains consistent response structure
//
// PARAMETERS:
//   - c: Fiber context
//   - message: Success message
//   - data: Optional data to include (can be nil)
//   - template: Template name for HTML responses (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
func SendSuccess(c *fiber.Ctx, message string, data interface{}, template ...string) error {
	responseData := ResponseData{
		Success: true,
		Message: message,
		Data:    data,
	}

	if WantsJSON(c) {
		return c.JSON(responseData)
	}

	// For HTML/HTMX requests
	templateName := "success.html"
	if len(template) > 0 && template[0] != "" {
		templateName = template[0]
	}

	return c.Render(templateName, responseData)
}

// SendCreated sends a 201 Created response
//
// ROLE: Resource Creation Response Handler
// - Sends 201 status for newly created resources
// - Includes location header when appropriate
// - Maintains REST conventions
//
// PARAMETERS:
//   - c: Fiber context
//   - data: Created resource data
//   - location: Optional location header value
//   - template: Template name for HTML responses (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
func SendCreated(c *fiber.Ctx, data interface{}, location string, template ...string) error {
	c.Status(fiber.StatusCreated)

	if location != "" {
		c.Set("Location", location)
	}

	if WantsJSON(c) {
		return c.JSON(ResponseData{
			Success: true,
			Data:    data,
			Message: "Resource created successfully",
		})
	}

	// For HTML/HTMX requests
	templateName := "created.html"
	if len(template) > 0 && template[0] != "" {
		templateName = template[0]
	}

	return c.Render(templateName, fiber.Map{
		"data":    data,
		"success": true,
		"message": "Resource created successfully",
	})
}

// SendNoContent sends a 204 No Content response
//
// ROLE: No Content Response Handler
// - Sends 204 status for successful operations with no content
// - Used for DELETE operations and updates
// - Maintains REST conventions
//
// PARAMETERS:
//   - c: Fiber context
//
// RETURNS:
//   - error: Fiber error if response fails
func SendNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// SendNotFound sends a 404 Not Found response
//
// ROLE: Not Found Response Handler
// - Sends 404 status for missing resources
// - Provides consistent not found messaging
// - Handles both API and web requests
//
// PARAMETERS:
//   - c: Fiber context
//   - resource: Name of the resource that wasn't found (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
func SendNotFound(c *fiber.Ctx, resource ...string) error {
	resourceName := "Resource"
	if len(resource) > 0 && resource[0] != "" {
		resourceName = resource[0]
	}

	message := resourceName + " not found"
	return SendError(c, fiber.StatusNotFound, message, "404.html")
}

// SendBadRequest sends a 400 Bad Request response
//
// ROLE: Bad Request Response Handler
// - Sends 400 status for invalid requests
// - Provides validation error messaging
// - Maintains consistent error structure
//
// PARAMETERS:
//   - c: Fiber context
//   - message: Specific error message (optional)
//
// RETURNS:
//   - error: Fiber error if response fails
func SendBadRequest(c *fiber.Ctx, message ...string) error {
	errorMessage := "Bad request"
	if len(message) > 0 && message[0] != "" {
		errorMessage = message[0]
	}

	return SendError(c, fiber.StatusBadRequest, errorMessage, "400.html")
}

// SendInternalError sends a 500 Internal Server Error response
//
// ROLE: Internal Error Response Handler
// - Sends 500 status for server errors
// - Provides generic error messaging for security
// - Logs detailed errors internally while hiding from client
//
// PARAMETERS:
//   - c: Fiber context
//   - message: Optional custom error message
//
// RETURNS:
//   - error: Fiber error if response fails
func SendInternalError(c *fiber.Ctx, message ...string) error {
	errorMessage := "Internal server error"
	if len(message) > 0 && message[0] != "" {
		errorMessage = message[0]
	}

	return SendError(c, fiber.StatusInternalServerError, errorMessage, "500.html")
}
