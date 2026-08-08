// Package middleware provides tracing and CORS middleware.
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/domain"
)

// Tracing adds trace_id and correlation_id to requests and tracks duration.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = "trace-" + time.Now().Format("20060102150405")
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-Id", traceID)

		corrID := c.GetHeader("X-Correlation-Id")
		if corrID != "" {
			c.Header("X-Correlation-Id", corrID)
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log request duration
		c.Set("duration_ms", duration.Milliseconds())
	}
}

// CORS adds Cross-Origin Resource Sharing headers.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-Id, X-Tenant-Id, Idempotency-Key, Crypto-Session-Id, Crypto-Request-Id, Crypto-Version, Crypto-Tenant-Id")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RequireScope checks that the user has the required scope.
func RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, exists := c.Get("cognito_context")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		authCtx, ok := ctx.(*domain.CognitoContext)
		if !ok {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Simple scope check — in production, validate against Cognito groups
		if authCtx.Role == "admin" || authCtx.Scope == requiredScope {
			c.Next()
			return
		}

		c.JSON(403, gin.H{"error": "forbidden", "required_scope": requiredScope})
		c.Abort()
	}
}
