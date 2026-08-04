// Package api provides the HTTP router for bff-api-go-collections-client-late-paying.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/application/usecases"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/interface/api/handlers"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/interface/api/middleware"
)

// NewRouter creates and configures the Gin router.
func NewRouter(uc *usecases.CollectionsUseCase) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Tracing())
	r.Use(middleware.CORS())

	// Health endpoints
	h := handlers.NewHandler(uc)
	r.GET("/health", h.Health)
	r.GET("/health/live", h.Liveness)
	r.GET("/health/ready", h.Readiness)

	// Collections routes
	collections := r.Group("/api/v1/collections")
	{
		activities := collections.Group("/activities")
		{
			activities.POST("", h.CreateActivity)
			activities.GET("", h.ListActivities)
			activities.GET("/:activity_id", h.GetActivity)
			activities.PATCH("/:activity_id", h.PatchActivity)
		}

		escalations := collections.Group("/escalations")
		{
			escalations.POST("", h.CreateEscalation)
			escalations.GET("", h.ListEscalations)
			escalations.GET("/:escalation_id", h.GetEscalation)
			escalations.POST("/:escalation_id/approve", h.ApproveEscalation)
			escalations.POST("/:escalation_id/reject", h.RejectEscalation)
			escalations.POST("/:escalation_id/assign", h.AssignEscalation)
		}

		portfolio := collections.Group("/portfolio/accounts")
		{
			portfolio.GET("", h.ListAccounts)
			portfolio.GET("/:account_id", h.GetAccount)
			portfolio.POST("/:account_id/assign", h.AssignAccount)
		}

		reports := collections.Group("/reports")
		{
			reports.GET("/agent-performance", h.ListAgentPerformance)
			reports.GET("/agent-performance/summary", h.AgentPerformanceSummary)
			reports.GET("/agent-performance/:agent_id", h.GetAgentPerformance)
		}
	}

	// Admin routes
	admin := r.Group("/api/v1/admin")
	{
		users := admin.Group("/users")
		{
			users.GET("", h.ListUsers)
			users.POST("", h.CreateUser)
			users.PATCH("/:user_id", h.PatchUser)
		}

		config := admin.Group("/config")
		{
			config.GET("", h.GetConfig)
			config.PUT("", h.PutConfig)
		}

		tenants := admin.Group("/tenants")
		{
			tenants.GET("", h.ListTenants)
			tenants.POST("", h.CreateTenant)
		}
	}

	// Auth proxy routes (forward to crypto-bff or cognito proxy)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", h.ForwardAuthLogin)
		auth.GET("/callback", h.ForwardAuthCallback)
		auth.POST("/logout", h.ForwardAuthLogout)
		auth.GET("/me", h.ForwardAuthMe)
		auth.POST("/dev-login", h.ForwardDevLogin)
	}

	return r
}

// CORS adds CORS headers.
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
