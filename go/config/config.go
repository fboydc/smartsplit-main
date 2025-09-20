package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// Config holds all configuration values for the application
type Config struct {
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          string
	PlaidProducts     string
	PlaidCountryCodes string
	PlaidRedirectURI  string
	AppPort           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	config := &Config{
		PlaidClientID:     os.Getenv("PLAID_CLIENT_ID"),
		PlaidSecret:       os.Getenv("PLAID_SECRET"),
		PlaidEnv:          os.Getenv("PLAID_ENV"),
		PlaidProducts:     os.Getenv("PLAID_PRODUCTS"),
		PlaidCountryCodes: os.Getenv("PLAID_COUNTRY_CODES"),
		PlaidRedirectURI:  os.Getenv("PLAID_REDIRECT_URI"),
		AppPort:           os.Getenv("APP_PORT"),
	}

	// Validate required fields
	if config.PlaidClientID == "" || config.PlaidSecret == "" {
		return nil, fmt.Errorf("PLAID_CLIENT_ID and PLAID_SECRET are required")
	}

	// Set defaults
	if config.PlaidProducts == "" {
		config.PlaidProducts = "transactions"
	}
	if config.PlaidCountryCodes == "" {
		config.PlaidCountryCodes = "US"
	}
	if config.PlaidEnv == "" {
		config.PlaidEnv = "sandbox"
	}
	if config.AppPort == "" {
		config.AppPort = "8000"
	}

	return config, nil
}

// ValidateConfig ensures all required configuration is present
func (c *Config) ValidateConfig() error {
	if c.PlaidClientID == "" {
		return fmt.Errorf("PLAID_CLIENT_ID is not set. Make sure to fill out the .env file")
	}
	if c.PlaidSecret == "" {
		return fmt.Errorf("PLAID_SECRET is not set. Make sure to fill out the .env file")
	}

	// Validate Plaid environment
	validEnvs := map[string]plaid.Environment{
		"sandbox":    plaid.Sandbox,
		"production": plaid.Production,
	}
	if _, exists := validEnvs[c.PlaidEnv]; !exists {
		return fmt.Errorf("invalid PLAID_ENV: %s. Must be 'sandbox' or 'production'", c.PlaidEnv)
	}

	return nil
}

// MustLoadConfig loads configuration and panics if it fails
func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := config.ValidateConfig(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	return config
}
