package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// WantsHTMX checks if the request expects an HTMX response
//
// ROLE: Request Content Type Detection
// - Detects if the client is making an HTMX request
// - Checks for HTMX-specific headers that indicate partial page updates
// - Used to determine whether to return full HTML pages or partial fragments
//
// HTMX Headers Checked:
//   - HX-Request: Present if request is made by HTMX
//   - HX-Trigger: Contains the ID of the element that triggered the request
//   - HX-Target: Contains the ID of the target element
//
// PARAMETERS:
//   - c: Fiber context containing request headers
//
// RETURNS:
//   - bool: true if request wants HTMX response, false otherwise
//
// USAGE EXAMPLE:
//
//	if utils.WantsHTMX(c) {
//	    return c.Render("partial.html", data)
//	}
//	return c.Render("full-page.html", data)
func WantsHTMX(c *fiber.Ctx) bool {
	// Check for HTMX request header
	if c.Get("HX-Request") != "" {
		return true
	}

	// Check for other HTMX-specific headers
	if c.Get("HX-Trigger") != "" || c.Get("HX-Target") != "" {
		return true
	}

	return false
}

// WantsJSON checks if the request expects a JSON response
//
// ROLE: Request Content Type Detection
// - Determines if the client expects JSON response format
// - Checks Accept header and other indicators for JSON preference
// - Used for API endpoints that can return both HTML and JSON
//
// Detection Methods:
//   - Accept header contains "application/json"
//   - Content-Type header is "application/json"
//   - X-Requested-With header indicates AJAX request
//   - Request path contains /api/ prefix
//
// PARAMETERS:
//   - c: Fiber context containing request headers
//
// RETURNS:
//   - bool: true if request wants JSON response, false otherwise
//
// USAGE EXAMPLE:
//
//	if utils.WantsJSON(c) {
//	    return c.JSON(data)
//	}
//	return c.Render("template.html", data)
func WantsJSON(c *fiber.Ctx) bool {
	// Check Accept header for JSON preference
	accept := c.Get("Accept")
	if strings.Contains(strings.ToLower(accept), "application/json") {
		return true
	}

	// Check Content-Type header
	contentType := c.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		return true
	}

	// Check for AJAX request
	if strings.ToLower(c.Get("X-Requested-With")) == "xmlhttprequest" {
		return true
	}

	// Check if request path indicates API endpoint
	path := c.Path()
	if strings.HasPrefix(path, "/api/") {
		return true
	}

	return false
}

// WantsHTML checks if the request expects an HTML response
//
// ROLE: Request Content Type Detection
// - Determines if the client expects HTML response format
// - Useful for endpoints that serve multiple content types
// - Default fallback for web browser requests
//
// Detection Methods:
//   - Accept header contains "text/html"
//   - No specific JSON or HTMX indicators present
//   - Browser user-agent detected
//
// PARAMETERS:
//   - c: Fiber context containing request headers
//
// RETURNS:
//   - bool: true if request wants HTML response, false otherwise
func WantsHTML(c *fiber.Ctx) bool {
	// If specifically wants JSON or HTMX, not HTML
	if WantsJSON(c) {
		return false
	}

	// Check Accept header for HTML preference
	accept := c.Get("Accept")
	if strings.Contains(strings.ToLower(accept), "text/html") {
		return true
	}

	// Default to HTML for browser requests
	userAgent := strings.ToLower(c.Get("User-Agent"))
	return strings.Contains(userAgent, "mozilla") ||
		strings.Contains(userAgent, "chrome") ||
		strings.Contains(userAgent, "safari") ||
		strings.Contains(userAgent, "edge")
}

// GetPreferredContentType determines the best response format for the request
//
// ROLE: Content Negotiation
// - Analyzes request headers to determine preferred response format
// - Returns standardized content type strings
// - Provides priority-based content negotiation
//
// Priority Order:
//  1. JSON (for API requests)
//  2. HTMX (for partial page updates)
//  3. HTML (for full page requests)
//
// PARAMETERS:
//   - c: Fiber context containing request headers
//
// RETURNS:
//   - string: "json", "htmx", or "html"
//
// USAGE EXAMPLE:
//
//	switch utils.GetPreferredContentType(c) {
//	case "json":
//	    return c.JSON(data)
//	case "htmx":
//	    return c.Render("partial.html", data)
//	default:
//	    return c.Render("full-page.html", data)
//	}
func GetPreferredContentType(c *fiber.Ctx) string {
	if WantsJSON(c) {
		return "json"
	}

	if WantsHTMX(c) {
		return "htmx"
	}

	return "html"
}

// IsAPIRequest checks if the request is targeting an API endpoint
//
// ROLE: Request Classification
// - Identifies API vs web page requests
// - Used for middleware and routing decisions
// - Helps separate API and web concerns
//
// PARAMETERS:
//   - c: Fiber context containing request information
//
// RETURNS:
//   - bool: true if request is for API endpoint, false otherwise
func IsAPIRequest(c *fiber.Ctx) bool {
	path := c.Path()
	return strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/v1/") ||
		strings.HasPrefix(path, "/v2/")
}

// IsWebRequest checks if the request is targeting a web page
//
// ROLE: Request Classification
// - Identifies web page vs API requests
// - Used for middleware and routing decisions
// - Helps separate web and API concerns
//
// PARAMETERS:
//   - c: Fiber context containing request information
//
// RETURNS:
//   - bool: true if request is for web page, false otherwise
func IsWebRequest(c *fiber.Ctx) bool {
	return !IsAPIRequest(c)
}
