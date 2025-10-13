package handlers

import (
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check if the service is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse "Service is healthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthResponse := models.HealthResponse{
		Status:  "ok",
		Message: "Servicio funcionando correctamente",
	}
	
	response.JSON(w, http.StatusOK, healthResponse)
}