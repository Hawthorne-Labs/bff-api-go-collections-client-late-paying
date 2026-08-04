// Package infrastructure provides configuration and core client.
package infrastructure

import (
	"net/url"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port            string
	CoreAPIBaseURL  string
	OTELServiceName string
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() Config {
	return Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		CoreAPIBaseURL:  getEnvOrDefault("CORE_API_BASE_URL", "http://core-api-go-collections-operations:8080"),
		OTELServiceName: getEnvOrDefault("OTEL_SERVICE_NAME", "bff-api-go-collections-client-late-paying"),
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// CoreClient is an HTTP client for core-api-go-collections-operations.
type CoreClient struct {
	BaseURL *url.URL
}

// NewCoreClient creates a new CoreClient from config.
func NewCoreClient(cfg Config) (*CoreClient, error) {
	baseURL, err := url.Parse(cfg.CoreAPIBaseURL)
	if err != nil {
		return nil, err
	}
	return &CoreClient{BaseURL: baseURL}, nil
}
