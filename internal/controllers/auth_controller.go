package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Login maneja el inicio de sesión
// @Summary Inicio de sesión de usuario
// @Description Autenticar usuario con nombre de usuario y contraseña
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciales de inicio de sesión"
// @Success 200 {object} models.LoginResponse "Inicio de sesión exitoso"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "Credenciales inválidas"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /auth/login [post]
func (h *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "El usuario y la contraseña son obligatorios.")
		return
	}

	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		response.Unauthorized(w, "Las credenciales no son correctas. ¡Inténtalo de nuevo!")
		return
	}

	response.JSON(w, http.StatusOK, loginResponse)
}

// ChangePassword maneja el cambio de contraseña
// @Summary Cambiar contraseña de usuario
// @Description Cambiar la contraseña del usuario autenticado
// @Tags Autenticación
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Solicitud de cambio de contraseña"
// @Success 200 {object} models.ChangePasswordResponse "Contraseña cambiada exitosamente"
// @Failure 400 {object} models.ErrorResponse "Solicitud incorrecta"
// @Failure 401 {object} models.ErrorResponse "No autorizado"
// @Failure 403 {object} models.ErrorResponse "Prohibido"
// @Failure 500 {object} models.ErrorResponse "Error interno del servidor"
// @Router /auth/change-password [post]
func (h *AuthController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var req models.ChangePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	if req.NewPassword == "" {
		response.BadRequest(w, "Debes ingresar la nueva contraseña.")
		return
	}

	if len(req.NewPassword) < 6 {
		response.BadRequest(w, "La nueva contraseña debe tener al menos 6 caracteres.")
		return
	}

	_, err := h.authService.ChangePassword(user.ID, &req)
	if err != nil {
		response.InternalServerError(w, "Tuvimos un problema al cambiar tu contraseña. Inténtalo más tarde.")
		return
	}

	response.JSON(w, http.StatusOK, models.ChangePasswordResponse{
		Message: "¡Contraseña actualizada con éxito! ✨",
	})
}

// Logout maneja el cierre de sesión
// @Summary Cerrar sesión
// @Description Cerrar la sesión del usuario actual
// @Tags Autenticación
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Cierre de sesión exitoso"
// @Router /auth/logout [post]
func (h *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "¡Sesión cerrada con éxito! 👋",
	})
}
