package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/infrastructure"
)

// AuditHandler proxies audit requests to the core audit service.
type AuditHandler struct {
	core *infrastructure.CoreClient
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(core *infrastructure.CoreClient) *AuditHandler {
	return &AuditHandler{core: core}
}

// Recent proxies GET /audit/recent to core.
func (h *AuditHandler) Recent(c *gin.Context) {
	tenantID := c.DefaultQuery("tenant_id", "default")
	resp, err := h.core.ForwardRequest(http.MethodGet, "/api/v1/audit/recent?tenant_id="+tenantID, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "DEPENDENT_SERVICE_FAILED", "message": "No fue posible consultar auditoria."}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// ByEntity proxies GET /audit/by-entity to core.
func (h *AuditHandler) ByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	tenantID := c.DefaultQuery("tenant_id", "default")
	resp, err := h.core.ForwardRequest(http.MethodGet, "/api/v1/audit/by-entity?entity_type="+entityType+"&entity_id="+entityID+"&tenant_id="+tenantID, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "DEPENDENT_SERVICE_FAILED", "message": "No fue posible consultar auditoria."}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// Integrity proxies GET /audit/integrity to core.
func (h *AuditHandler) Integrity(c *gin.Context) {
	tenantID := c.DefaultQuery("tenant_id", "default")
	resp, err := h.core.ForwardRequest(http.MethodGet, "/api/v1/audit/integrity?tenant_id="+tenantID, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "DEPENDENT_SERVICE_FAILED", "message": "No fue posible consultar integridad."}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}
