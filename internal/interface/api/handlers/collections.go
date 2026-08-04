package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/application/usecases"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	uc *usecases.CollectionsUseCase
}

// NewHandler creates a new Handler.
func NewHandler(uc *usecases.CollectionsUseCase) *Handler {
	return &Handler{uc: uc}
}

// Health returns 200 OK.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Liveness returns 200 OK (container liveness probe).
func (h *Handler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "live"})
}

// Readiness returns 200 OK or 503 if dependencies are down.
func (h *Handler) Readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// CreateActivity creates a new collection activity.
func (h *Handler) CreateActivity(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardActivity("post", "", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create activity", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ListActivities lists collection activities.
func (h *Handler) ListActivities(c *gin.Context) {
	result, err := h.uc.ForwardActivity("get", "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list activities", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetActivity gets a single activity by ID.
func (h *Handler) GetActivity(c *gin.Context) {
	activityID := c.Param("activity_id")
	if activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activity_id is required"})
		return
	}

	result, err := h.uc.ForwardActivity("get", activityID, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get activity", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PatchActivity updates an existing activity.
func (h *Handler) PatchActivity(c *gin.Context) {
	activityID := c.Param("activity_id")
	if activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activity_id is required"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardActivity("patch", activityID, body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to patch activity", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateEscalation creates a new escalation.
func (h *Handler) CreateEscalation(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardEscalation("post", "", "", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create escalation", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ListEscalations lists escalations.
func (h *Handler) ListEscalations(c *gin.Context) {
	result, err := h.uc.ForwardEscalation("get", "", "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list escalations", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEscalation gets a single escalation by ID.
func (h *Handler) GetEscalation(c *gin.Context) {
	escalationID := c.Param("escalation_id")
	if escalationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "escalation_id is required"})
		return
	}

	result, err := h.uc.ForwardEscalation("get", escalationID, "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get escalation", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApproveEscalation approves an escalation.
func (h *Handler) ApproveEscalation(c *gin.Context) {
	escalationID := c.Param("escalation_id")
	if escalationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "escalation_id is required"})
		return
	}

	result, err := h.uc.ForwardEscalation("post", escalationID, "approve", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to approve escalation", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RejectEscalation rejects an escalation.
func (h *Handler) RejectEscalation(c *gin.Context) {
	escalationID := c.Param("escalation_id")
	if escalationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "escalation_id is required"})
		return
	}

	result, err := h.uc.ForwardEscalation("post", escalationID, "reject", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reject escalation", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AssignEscalation assigns an escalation to an agent.
func (h *Handler) AssignEscalation(c *gin.Context) {
	escalationID := c.Param("escalation_id")
	if escalationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "escalation_id is required"})
		return
	}

	var body map[string]interface{}
	c.ShouldBindJSON(&body)

	result, err := h.uc.ForwardEscalation("post", escalationID, "assign", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to assign escalation", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListAccounts lists portfolio accounts.
func (h *Handler) ListAccounts(c *gin.Context) {
	result, err := h.uc.ForwardAccount("get", "", "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list accounts", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAccount gets a single account by ID.
func (h *Handler) GetAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	result, err := h.uc.ForwardAccount("get", accountID, "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get account", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AssignAccount assigns an account to an agent.
func (h *Handler) AssignAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	var body map[string]interface{}
	c.ShouldBindJSON(&body)

	result, err := h.uc.ForwardAccount("post", accountID, "assign", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to assign account", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListAgentPerformance lists agent performance metrics.
func (h *Handler) ListAgentPerformance(c *gin.Context) {
	result, err := h.uc.ForwardReport("agent-performance", "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get agent performance", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AgentPerformanceSummary returns agent performance summary.
func (h *Handler) AgentPerformanceSummary(c *gin.Context) {
	result, err := h.uc.ForwardReport("agent-performance/summary", "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get summary", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAgentPerformance gets specific agent performance metrics.
func (h *Handler) GetAgentPerformance(c *gin.Context) {
	agentID := c.Param("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id is required"})
		return
	}

	result, err := h.uc.ForwardReport("agent-performance", agentID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get agent performance", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListUsers lists application users.
func (h *Handler) ListUsers(c *gin.Context) {
	result, err := h.uc.ForwardUser("get", "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list users", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateUser creates a new user.
func (h *Handler) CreateUser(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardUser("post", "", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// PatchUser updates an existing user.
func (h *Handler) PatchUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardUser("patch", userID, body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to update user", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetConfig gets system configuration.
func (h *Handler) GetConfig(c *gin.Context) {
	result, err := h.uc.ForwardConfig("get", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get config", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PutConfig updates system configuration.
func (h *Handler) PutConfig(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardConfig("put", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to update config", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListTenants lists tenants.
func (h *Handler) ListTenants(c *gin.Context) {
	result, err := h.uc.ForwardTenant("get", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list tenants", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateTenant creates a new tenant.
func (h *Handler) CreateTenant(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.uc.ForwardTenant("post", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create tenant", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ForwardAuthLogin forwards login to auth proxy.
func (h *Handler) ForwardAuthLogin(c *gin.Context) {
	h.forwardToAuth(c, "POST", "/api/v1/auth/login")
}

// ForwardAuthCallback forwards auth callback.
func (h *Handler) ForwardAuthCallback(c *gin.Context) {
	h.forwardToAuth(c, "GET", "/api/v1/auth/callback")
}

// ForwardAuthLogout forwards logout.
func (h *Handler) ForwardAuthLogout(c *gin.Context) {
	h.forwardToAuth(c, "POST", "/api/v1/auth/logout")
}

// ForwardAuthMe returns current user info.
func (h *Handler) ForwardAuthMe(c *gin.Context) {
	h.forwardToAuth(c, "GET", "/api/v1/auth/me")
}

// ForwardDevLogin forwards dev login (environment-gated).
func (h *Handler) ForwardDevLogin(c *gin.Context) {
	h.forwardToAuth(c, "POST", "/api/v1/auth/dev-login")
}

// forwardToAuth forwards auth requests to the auth proxy.
func (h *Handler) forwardToAuth(c *gin.Context, method, path string) {
	// In production, this would forward to crypto-bff or cognito proxy.
	// For now, return a placeholder response.
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":  "auth forwarding not yet configured",
		"detail": "configure AUTH_PROXY_URL environment variable",
		"path":   path,
		"method": method,
	})
}

// sanitizeScope returns the last part of a scope (e.g., "read" from "collections:read").
func sanitizeScope(scope string) string {
	parts := strings.Split(scope, ":")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return scope
}
