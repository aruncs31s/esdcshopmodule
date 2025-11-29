package handler

import (
	"net/http"
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "ready",
	})
}
