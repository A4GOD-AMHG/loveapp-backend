// Paquete auth provee utilidades de autenticación basadas en JSON Web Tokens (JWT).
// Los tokens generados no tienen fecha de expiración y son válidos indefinidamente
// mientras la clave secreta no cambie.
package auth

import (
	"errors"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// Claims define la estructura de los datos (claims) embebidos en el token JWT.
// Incluye datos del usuario (UserID, Username) además de los claims estándar de JWT.
type Claims struct {
	UserID   int64  `json:"user_id"` // ID único del usuario autenticado
	Username string `json:"username"` // Nombre de usuario del titular del token
	jwt.RegisteredClaims               // Claims estándar JWT (IssuedAt, NotBefore, Subject, etc.)
}

// GenerateToken crea y firma un nuevo token JWT para el usuario dado.
// El token no tiene fecha de expiración (sin ExpiresAt), por lo que es válido indefinidamente.
// Utiliza el algoritmo HMAC-SHA256 (HS256) y la clave secreta configurada en AppConfig.
func GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Fecha de emisión del token
			NotBefore: jwt.NewNumericDate(time.Now()), // El token no es válido antes de este momento
			Subject:   user.Username,                  // Sujeto del token (identificador del usuario)
		},
	}

	// Crear el token con el método de firma HS256 y los claims definidos
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta de la configuración
	return token.SignedString(config.AppConfig.JWT.Secret)
}

// ValidateToken parsea y valida un token JWT firmado.
// Verifica que el método de firma sea HMAC y que la firma sea correcta con la clave secreta.
// Retorna los claims extraídos del token si es válido, o un error si no lo es.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el algoritmo de firma sea HMAC (previene ataques de algoritmo "none")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.AppConfig.JWT.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extraer y retornar los claims si el token es válido
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
