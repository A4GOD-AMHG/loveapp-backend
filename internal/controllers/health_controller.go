package controllers

import (
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck maneja las solicitudes de verificación de estado
// @Summary Verificación de estado
// @Description Verificar si el servicio está funcionando
// @Tags Estado
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse "Servicio saludable"
// @Router /health [get]
func (h *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthResponse := models.HealthResponse{
		Status:  "ok",
		Message: "Servicio funcionando correctamente",
	}

	response.JSON(w, http.StatusOK, healthResponse)
}