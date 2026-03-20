package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// DeviceController maneja el registro y eliminación de tokens push por dispositivo.
type DeviceController struct {
	service services.DeviceService
}

// NewDeviceController crea un controlador de dispositivos push.
func NewDeviceController(service services.DeviceService) *DeviceController {
	return &DeviceController{service: service}
}

// RegisterPushToken registra o actualiza el token push del dispositivo autenticado.
func (h *DeviceController) RegisterPushToken(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	var req models.RegisterPushTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos del dispositivo.")
		return
	}

	if err := h.service.RegisterPushToken(user.ID, &req); err != nil {
		response.BadRequest(w, "Los datos del token push no son válidos.")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Token push registrado correctamente"})
}

// DeletePushToken elimina la asociación de un token push con el usuario autenticado.
func (h *DeviceController) DeletePushToken(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
		return
	}

	var req models.DeletePushTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa el token push.")
		return
	}

	if err := h.service.DeletePushToken(user.ID, req.PushToken); err != nil {
		response.BadRequest(w, "El token push no es válido.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
