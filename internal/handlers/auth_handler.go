package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Datos de entrada inválidos")
		return
	}
	
	// Validate required fields
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "Usuario y contraseña son requeridos")
		return
	}
	
	// Authenticate user
	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, loginResponse)
}

// ChangePassword handles password change
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Password change request"
// @Success 200 {object} models.ChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		response.Unauthorized(w, "Usuario no autenticado")
		return
	}
	
	var req models.ChangePasswordRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Datos de entrada inválidos")
		return
	}
	
	// Validate required fields
	if req.OldPassword == "" || req.NewPassword == "" {
		response.BadRequest(w, "Contraseña actual y nueva contraseña son requeridas")
		return
	}
	
	// Validate new password length
	if len(req.NewPassword) < 6 {
		response.BadRequest(w, "La nueva contraseña debe tener al menos 6 caracteres")
		return
	}
	
	// Change password
	changeResponse, err := h.authService.ChangePassword(user.ID, &req)
	if err != nil {
		if err.Error() == "contraseña actual incorrecta" {
			response.Forbidden(w, err.Error())
			return
		}
		response.InternalServerError(w, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, changeResponse)
}