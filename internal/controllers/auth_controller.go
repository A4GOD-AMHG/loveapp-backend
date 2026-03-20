// Paquete controllers contiene los manejadores HTTP de la aplicación,
// responsables de recibir solicitudes, validar datos y delegar al servicio correspondiente.
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// AuthController maneja los endpoints relacionados con la autenticación de usuarios.
type AuthController struct {
	authService *services.AuthService // Servicio de autenticación con la lógica de negocio
}

// NewAuthController crea y retorna una nueva instancia de AuthController
// con su servicio de autenticación ya inicializado.
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

	// Decodificar el cuerpo JSON de la solicitud en la estructura de login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	// Validar que los campos obligatorios no estén vacíos
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "El usuario y la contraseña son obligatorios.")
		return
	}

	// Delegar la autenticación al servicio y obtener el token JWT
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
	// Obtener el usuario autenticado inyectado por el middleware de autenticación
	user := r.Context().Value("user").(*models.User)

	var req models.ChangePasswordRequest

	// Decodificar el cuerpo JSON con la nueva contraseña
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Por favor, revisa los datos ingresados.")
		return
	}

	// Validar que la nueva contraseña no esté vacía
	if req.NewPassword == "" {
		response.BadRequest(w, "Debes ingresar la nueva contraseña.")
		return
	}

	// Validar longitud mínima de la nueva contraseña
	if len(req.NewPassword) < 6 {
		response.BadRequest(w, "La nueva contraseña debe tener al menos 6 caracteres.")
		return
	}

	// Delegar el cambio de contraseña al servicio de autenticación
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
	// El cierre de sesión es manejado en el cliente eliminando el token JWT.
	// El servidor simplemente confirma la operación.
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "¡Sesión cerrada con éxito! 👋",
	})
}
