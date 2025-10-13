package services

import (
	"fmt"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}
	
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}
	
	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("error al generar token")
	}
	
	// Clear password from user object
	user.Password = ""
	
	return &models.LoginResponse{
		Message: "Inicio de sesión exitoso",
		Token:   token,
		User:    *user,
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID int64, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	// Get current password hash
	currentHash, err := s.userRepo.GetPasswordHash(userID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	
	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword))
	if err != nil {
		return nil, fmt.Errorf("contraseña actual incorrecta")
	}
	
	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al procesar nueva contraseña")
	}
	
	// Update password
	err = s.userRepo.UpdatePassword(userID, string(newHash))
	if err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña")
	}
	
	return &models.ChangePasswordResponse{
		Message: "Contraseña cambiada exitosamente",
	}, nil
}

// GetUserByID gets a user by ID
func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	
	return user, nil
}