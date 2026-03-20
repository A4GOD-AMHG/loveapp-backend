// Paquete services implementa la lógica de negocio de la aplicación,
// actuando como intermediario entre los controladores y los repositorios.
package services

import (
	"fmt"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

// AuthService encapsula la lógica de negocio relacionada con la autenticación de usuarios.
type AuthService struct {
	userRepo    *repository.UserRepository // Repositorio de usuarios para acceso a datos
	messageRepo repository.MessageRepository
}

// NewAuthService crea y retorna una nueva instancia de AuthService
// con su repositorio de usuarios ya inicializado.
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:    repository.NewUserRepository(),
		messageRepo: repository.NewMessageRepository(),
	}
}

// Login autentica a un usuario verificando sus credenciales y genera un token JWT.
// Retorna la respuesta de login con el token y los datos del usuario si las credenciales son válidas.
// Retorna error con mensaje genérico "credenciales inválidas" para no revelar si el usuario existe.
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Buscar al usuario por nombre de usuario en la base de datos
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Verificar que la contraseña proporcionada coincida con el hash almacenado
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Generar el token JWT para la sesión del usuario
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("error al generar token")
	}

	// Limpiar el hash de la contraseña antes de incluir el usuario en la respuesta
	user.Password = ""

	unreadCount, err := s.messageRepo.CountUnreadByReceiver(uint(user.ID))
	if err != nil {
		return nil, fmt.Errorf("error al calcular mensajes no leídos")
	}

	return &models.LoginResponse{
		Message:     "Inicio de sesión exitoso",
		Token:       token,
		User:        *user,
		UnreadCount: unreadCount,
	}, nil
}

// ChangePassword actualiza la contraseña del usuario especificado.
// Genera un nuevo hash bcrypt para la nueva contraseña antes de almacenarla.
// Retorna error si el usuario no existe o si falla el proceso de hashing o actualización.
func (s *AuthService) ChangePassword(userID int64, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	// Verificar que el usuario exista en el sistema
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Generar el hash bcrypt de la nueva contraseña con el costo predeterminado
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al procesar nueva contraseña")
	}

	// Persistir el nuevo hash en la base de datos
	err = s.userRepo.UpdatePassword(userID, string(newHash))
	if err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña")
	}

	return &models.ChangePasswordResponse{
		Message: "Contraseña cambiada exitosamente",
	}, nil
}

// GetUserByID obtiene los datos de un usuario por su ID sin incluir la contraseña.
// Se usa principalmente en el middleware de autenticación para cargar el usuario en el contexto.
// Retorna error si el usuario no existe.
func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	return user, nil
}
