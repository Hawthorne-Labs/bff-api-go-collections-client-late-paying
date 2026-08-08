package infrastructure

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all environment-dependent configuration for the client-late-paying BFF.
type Config struct {
	CoreAPIURL          string
	OIDCProvider        string
	CryptoBFFURL        string
	CryptoEnabled       bool
	CryptoBFFTimeoutSec float64
	ServerHost          string
	ServerPort          int
	LogLevel            string
	OtelServiceName     string
	OtelExporterURL     string
	RedisURL            string
	SessionBackend      string
	RateLimitRequests   int
	RateLimitWindowSec  int
	MaxRequestBodyBytes int
	RequestTimeoutSec   float64
	CORSEnabled         bool
	CryptoSessionSecret string
	CryptoSessionIssuer string
	CryptoSessionTTL    int
}

// LoadConfig reads configuration from environment variables with safe defaults.
func LoadConfig() Config {
	return Config{
		CoreAPIURL:          envOr("CORE_API_URL", "http://localhost:8080"),
		OIDCProvider:        envOr("OIDC_PROVIDER", "keycloak"),
		CryptoBFFURL:        envOr("CRYPTO_BFF_URL", "http://crypto-bff:9000"),
		CryptoEnabled:       isTrueEnv("CRYPTO_ENABLED"),
		CryptoBFFTimeoutSec: parseFloatEnv("CRYPTO_BFF_TIMEOUT_SECONDS", 3.0),
		ServerHost:          envOr("SERVER_HOST", "0.0.0.0"),
		ServerPort:          parseIntEnv("SERVER_PORT", 8080),
		LogLevel:            envOr("LOG_LEVEL", "info"),
		OtelServiceName:     envOr("OTEL_SERVICE_NAME", "bff-api-go-collections-client-late-paying"),
		OtelExporterURL:     envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"),
		RedisURL:            envOr("REDIS_URL", "redis://127.0.0.1:6379/0"),
		SessionBackend:      envOr("SESSION_BACKEND", "redis"),
		RateLimitRequests:   parseIntEnv("RATE_LIMIT_REQUESTS", 60),
		RateLimitWindowSec:  parseIntEnv("RATE_LIMIT_WINDOW_SEC", 60),
		MaxRequestBodyBytes: parseIntEnv("MAX_REQUEST_BODY_BYTES", 65536),
		RequestTimeoutSec:   parseFloatEnv("REQUEST_TIMEOUT_SECONDS", 30.0),
		CORSEnabled:         !isFalseEnv("BFF_DISABLE_CORS"),
		CryptoSessionSecret: envOr("CRYPTO_SESSION_SECRET", "dev-crypto-secret-change-in-production"),
		CryptoSessionIssuer: envOr("CRYPTO_SESSION_ISSUER", "bff-api-go-collections-client-late-paying"),
		CryptoSessionTTL:    parseIntEnv("CRYPTO_SESSION_TTL", 900),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseIntEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var n int
	_, err := fmt.Sscanf(v, "%d", &n)
	if err != nil {
		return fallback
	}
	if n <= 0 {
		return fallback
	}
	return n
}

func parseFloatEnv(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var f float64
	_, err := fmt.Sscanf(v, "%f", &f)
	if err != nil {
		return fallback
	}
	if f <= 0 {
		return fallback
	}
	return f
}

func isTrueEnv(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	v = strings.ToLower(v)
	return v == "true" || v == "1" || v == "yes"
}

func isFalseEnv(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	v = strings.ToLower(v)
	return v == "false" || v == "0" || v == "no"
}
