package utils

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	// Server settings
	Port        string
	Host        string
	Environment string

	// Database settings
	DatabaseURL  string
	DatabaseName string

	// Security settings
	JWTSecret      string
	AllowedOrigins string

	// Application settings
	AppName    string
	AppVersion string
	LogLevel   string

	// Feature flags
	EnableRateLimit bool
	EnableCORS      bool
	EnableLogger    bool
}

// LoadConfig loads configuration from environment variables
//
// ROLE: Configuration Management
// - Loads all configuration from environment variables
// - Provides default values for missing configurations
// - Validates required configuration values
//
// RETURNS:
//   - *Config: Configuration struct with loaded values
//   - error: Error if required configuration is missing
func LoadConfig() (*Config, error) {
	config := &Config{
		// Server defaults
		Port:        GetEnv("PORT", "3000"),
		Host:        GetEnv("HOST", "localhost"),
		Environment: GetEnv("ENVIRONMENT", "development"),

		// Database defaults
		DatabaseURL:  GetEnv("DATABASE_URL", ""),
		DatabaseName: GetEnv("DATABASE_NAME", "lacpa_db"),

		// Security defaults
		JWTSecret:      GetEnv("JWT_SECRET", "your-secret-key"),
		AllowedOrigins: GetEnv("ALLOWED_ORIGINS", "*"),

		// Application defaults
		AppName:    GetEnv("APP_NAME", "LACPA Backend"),
		AppVersion: GetEnv("APP_VERSION", "1.0.0"),
		LogLevel:   GetEnv("LOG_LEVEL", "info"),

		// Feature flags
		EnableRateLimit: GetEnvBool("ENABLE_RATE_LIMIT", true),
		EnableCORS:      GetEnvBool("ENABLE_CORS", true),
		EnableLogger:    GetEnvBool("ENABLE_LOGGER", true),
	}

	return config, nil
}

// GetEnv gets an environment variable with a default value
//
// ROLE: Environment Variable Helper
// - Safely retrieves environment variables
// - Provides default values for missing variables
// - Centralizes environment variable access
//
// PARAMETERS:
//   - key: Environment variable name
//   - defaultValue: Default value if variable is not set
//
// RETURNS:
//   - string: Environment variable value or default
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvInt gets an integer environment variable with a default value
//
// PARAMETERS:
//   - key: Environment variable name
//   - defaultValue: Default value if variable is not set or invalid
//
// RETURNS:
//   - int: Environment variable value or default
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetEnvBool gets a boolean environment variable with a default value
//
// PARAMETERS:
//   - key: Environment variable name
//   - defaultValue: Default value if variable is not set or invalid
//
// RETURNS:
//   - bool: Environment variable value or default
func GetEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(key))
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

// GetEnvSlice gets a comma-separated environment variable as a slice
//
// PARAMETERS:
//   - key: Environment variable name
//   - defaultValue: Default slice if variable is not set
//
// RETURNS:
//   - []string: Environment variable values as slice
func GetEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Split by comma and trim whitespace
	var result []string
	parts := strings.Split(value, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// IsDevelopment checks if the application is running in development mode
//
// ROLE: Environment Detection
// - Determines if application is in development mode
// - Used for conditional behavior (debug logging, etc.)
//
// PARAMETERS:
//   - config: Configuration struct
//
// RETURNS:
//   - bool: true if in development mode
func IsDevelopment(config *Config) bool {
	return strings.ToLower(config.Environment) == "development" ||
		strings.ToLower(config.Environment) == "dev"
}

// IsProduction checks if the application is running in production mode
//
// ROLE: Environment Detection
// - Determines if application is in production mode
// - Used for conditional behavior (security, performance)
//
// PARAMETERS:
//   - config: Configuration struct
//
// RETURNS:
//   - bool: true if in production mode
func IsProduction(config *Config) bool {
	return strings.ToLower(config.Environment) == "production" ||
		strings.ToLower(config.Environment) == "prod"
}

// IsStaging checks if the application is running in staging mode
//
// ROLE: Environment Detection
// - Determines if application is in staging mode
// - Used for conditional behavior (testing features)
//
// PARAMETERS:
//   - config: Configuration struct
//
// RETURNS:
//   - bool: true if in staging mode
func IsStaging(config *Config) bool {
	return strings.ToLower(config.Environment) == "staging" ||
		strings.ToLower(config.Environment) == "stage"
}

// GetDatabaseURL constructs database URL from components
//
// ROLE: Database URL Construction
// - Builds database connection string
// - Handles different database URL formats
// - Provides fallback construction from components
//
// PARAMETERS:
//   - config: Configuration struct
//
// RETURNS:
//   - string: Complete database URL
func GetDatabaseURL(config *Config) string {
	if config.DatabaseURL != "" {
		return config.DatabaseURL
	}

	// Construct from individual components if needed
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	dbname := config.DatabaseName
	sslmode := GetEnv("DB_SSLMODE", "disable")

	if password != "" {
		return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
	}

	return "postgres://" + user + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
}

// ValidateConfig validates required configuration values
//
// ROLE: Configuration Validation
// - Ensures required configuration is present
// - Validates configuration value formats
// - Returns descriptive error messages
//
// PARAMETERS:
//   - config: Configuration struct to validate
//
// RETURNS:
//   - error: Validation error if configuration is invalid
func ValidateConfig(config *Config) error {
	ve := NewValidationErrors()

	// Validate required fields
	ValidateRequired(ve, "app_name", config.AppName)
	ValidateRequired(ve, "environment", config.Environment)

	// Validate port is numeric
	if _, err := strconv.Atoi(config.Port); err != nil {
		ve.AddError("port", "Port must be a valid number", config.Port)
	}

	// Validate environment values
	validEnvs := []string{"development", "dev", "staging", "stage", "production", "prod"}
	isValidEnv := false
	for _, validEnv := range validEnvs {
		if strings.ToLower(config.Environment) == validEnv {
			isValidEnv = true
			break
		}
	}
	if !isValidEnv {
		ve.AddError("environment", "Environment must be one of: "+strings.Join(validEnvs, ", "), config.Environment)
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}
