package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/AliSleiman0/Lacpa/utils"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Missing authorization header",
			"success": false,
		})
	}

	// Check if it starts with "Bearer "
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid authorization header format",
			"success": false,
		})
	}

	tokenString := tokenParts[1]

	// Validate token
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid or expired token",
			"success": false,
		})
	}

	// Store user info in context for later use
	c.Locals("userID", claims.UserID)
	c.Locals("lacpaID", claims.LACPAID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "No role found in context",
				"success": false,
			})
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Insufficient permissions",
			"success": false,
		})
	}
}

// OptionalAuthMiddleware validates JWT token if present, but doesn't require it
func OptionalAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Next()
	}

	tokenString := tokenParts[1]
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Next()
	}

	// Store user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("lacpaID", claims.LACPAID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}
