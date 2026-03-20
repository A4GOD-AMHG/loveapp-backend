package controllers

import (
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// HealthController maneja endpoints simples de salud del servicio.
type HealthController struct{}

// NewHealthController crea una nueva instancia del controlador de salud.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check retorna el estado operativo basico del backend.
func (h *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "loveapp-backend",
	})
}
