package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds all configuration for the paperwork service
type Config struct {
	// Server configuration
	Port        string `json:"port"`
	Environment string `json:"environment"`

	// Supabase configuration
	SupabaseURL string `json:"supabase_url"`
	SupabaseKey string `json:"supabase_key"`

	// Template paths
	TemplatesPath   string `json:"templates_path"`
	FontsPath       string `json:"fonts_path"`
	BackgroundsPath string `json:"backgrounds_path"`
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		SupabaseURL:     getEnvRequired("SUPABASE_URL"),
		SupabaseKey:     getEnvRequired("SUPABASE_KEY"),
		TemplatesPath:   getEnv("TEMPLATES_PATH", "./templates"),
		FontsPath:       getEnv("FONTS_PATH", "./templates/fonts"),
		BackgroundsPath: getEnv("BACKGROUNDS_PATH", "./templates/backgrounds"),
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvRequired gets a required environment variable
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// getEnvInt gets an environment variable as an integer
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets an environment variable as a boolean
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}