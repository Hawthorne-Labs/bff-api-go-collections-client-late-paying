package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/infrastructure/fieldcrypto"
)

// CryptoSessionHandler handles the local P-256 ECDH crypto-session handshake.
type CryptoSessionHandler struct {
	mgr *fieldcrypto.SessionManager
}

// NewCryptoSessionHandler creates a new handler.
func NewCryptoSessionHandler(mgr *fieldcrypto.SessionManager) *CryptoSessionHandler {
	return &CryptoSessionHandler{mgr: mgr}
}

type handshakeRequest struct {
	ClientPublicKey string `json:"clientPublicKey"`
}

// Handshake handles POST /api/v1/collections/crypto-session
func (h *CryptoSessionHandler) Handshake(c *gin.Context) {
	var req handshakeRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil || req.ClientPublicKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": 90100, "message": "Solicitud de cifrado invalida."}})
		return
	}

	sub := c.GetString("user_id")
	if sub == "" {
		sub = "user"
	}
	scope := c.GetString("scope")
	if scope == "" {
		scope = "collections:read"
	}

	result, err := h.mgr.Handshake(req.ClientPublicKey, sub, scope)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": 90101, "message": "Material de cifrado invalido."}})
		return
	}

	c.JSON(http.StatusOK, result)
}
