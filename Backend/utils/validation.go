package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Field+": "+err.Message)
	}

	return strings.Join(messages, ", ")
}

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (ve *ValidationErrors) AddError(field, message, value string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// GetQueryParam gets a query parameter with default value
//
// ROLE: Query Parameter Helper
// - Safely extracts query parameters with defaults
// - Provides type conversion helpers
// - Handles missing parameters gracefully
//
// PARAMETERS:
//   - c: Fiber context
//   - key: Parameter name
//   - defaultValue: Default value if parameter is missing
//
// RETURNS:
//   - string: Parameter value or default
func GetQueryParam(c *fiber.Ctx, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt gets an integer query parameter with default value
//
// PARAMETERS:
//   - c: Fiber context
//   - key: Parameter name
//   - defaultValue: Default value if parameter is missing or invalid
//
// RETURNS:
//   - int: Parameter value or default
func GetQueryParamInt(c *fiber.Ctx, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetQueryParamBool gets a boolean query parameter with default value
//
// PARAMETERS:
//   - c: Fiber context
//   - key: Parameter name
//   - defaultValue: Default value if parameter is missing or invalid
//
// RETURNS:
//   - bool: Parameter value or default
func GetQueryParamBool(c *fiber.Ctx, key string, defaultValue bool) bool {
	value := strings.ToLower(c.Query(key))
	if value == "" {
		return defaultValue
	}

	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// ValidateRequired validates that a field is not empty
//
// ROLE: Field Validation Helper
// - Validates required fields
// - Provides consistent validation messages
// - Supports string and numeric validations
//
// PARAMETERS:
//   - ve: ValidationErrors instance to add errors to
//   - field: Field name
//   - value: Field value to validate
//
// RETURNS:
//   - bool: true if valid, false if invalid
func ValidateRequired(ve *ValidationErrors, field, value string) bool {
	if strings.TrimSpace(value) == "" {
		ve.AddError(field, "This field is required", value)
		return false
	}
	return true
}

// ValidateEmail validates email format
//
// PARAMETERS:
//   - ve: ValidationErrors instance to add errors to
//   - field: Field name
//   - email: Email value to validate
//
// RETURNS:
//   - bool: true if valid, false if invalid
func ValidateEmail(ve *ValidationErrors, field, email string) bool {
	if email == "" {
		return true // Don't validate empty emails here, use ValidateRequired separately
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		ve.AddError(field, "Please enter a valid email address", email)
		return false
	}

	return true
}

// ValidateMinLength validates minimum string length
//
// PARAMETERS:
//   - ve: ValidationErrors instance to add errors to
//   - field: Field name
//   - value: String value to validate
//   - minLength: Minimum required length
//
// RETURNS:
//   - bool: true if valid, false if invalid
func ValidateMinLength(ve *ValidationErrors, field, value string, minLength int) bool {
	if len(value) < minLength {
		ve.AddError(field, "Must be at least "+strconv.Itoa(minLength)+" characters long", value)
		return false
	}
	return true
}

// ValidateMaxLength validates maximum string length
//
// PARAMETERS:
//   - ve: ValidationErrors instance to add errors to
//   - field: Field name
//   - value: String value to validate
//   - maxLength: Maximum allowed length
//
// RETURNS:
//   - bool: true if valid, false if invalid
func ValidateMaxLength(ve *ValidationErrors, field, value string, maxLength int) bool {
	if len(value) > maxLength {
		ve.AddError(field, "Must be no more than "+strconv.Itoa(maxLength)+" characters long", value)
		return false
	}
	return true
}

// ValidatePasswordStrength validates password strength
//
// PARAMETERS:
//   - ve: ValidationErrors instance to add errors to
//   - field: Field name
//   - password: Password to validate
//
// RETURNS:
//   - bool: true if valid, false if invalid
func ValidatePasswordStrength(ve *ValidationErrors, field, password string) bool {
	if len(password) < 8 {
		ve.AddError(field, "Password must be at least 8 characters long", "")
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		ve.AddError(field, "Password must contain at least one uppercase letter", "")
		return false
	}

	if !hasLower {
		ve.AddError(field, "Password must contain at least one lowercase letter", "")
		return false
	}

	if !hasDigit {
		ve.AddError(field, "Password must contain at least one digit", "")
		return false
	}

	if !hasSpecial {
		ve.AddError(field, "Password must contain at least one special character", "")
		return false
	}

	return true
}

// SanitizeString removes potentially dangerous characters from string input
//
// ROLE: Input Sanitization
// - Removes HTML tags and scripts
// - Trims whitespace
// - Prevents basic XSS attacks
//
// PARAMETERS:
//   - input: String to sanitize
//
// RETURNS:
//   - string: Sanitized string
func SanitizeString(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	sanitized := htmlRegex.ReplaceAllString(input, "")

	// Remove script content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized = scriptRegex.ReplaceAllString(sanitized, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// Paginate calculates pagination values
//
// ROLE: Pagination Helper
// - Calculates offset and limit for database queries
// - Provides pagination metadata
// - Handles edge cases and invalid inputs
//
// PARAMETERS:
//   - page: Current page number (1-based)
//   - pageSize: Number of items per page
//   - totalItems: Total number of items
//
// RETURNS:
//   - offset: Database offset value
//   - limit: Database limit value
//   - metadata: Pagination metadata
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
}

func Paginate(page, pageSize, totalItems int) (offset, limit int, meta PaginationMeta) {
	// Ensure valid page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	// Calculate values
	offset = (page - 1) * pageSize
	limit = pageSize
	totalPages := (totalItems + pageSize - 1) / pageSize // Ceiling division

	// Build metadata
	meta = PaginationMeta{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	// Set next/prev page numbers
	if meta.HasNext {
		nextPage := page + 1
		meta.NextPage = &nextPage
	}
	if meta.HasPrev {
		prevPage := page - 1
		meta.PrevPage = &prevPage
	}

	return offset, limit, meta
}
