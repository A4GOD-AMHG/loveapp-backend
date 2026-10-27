// Paquete middleware contiene los middlewares HTTP de la aplicación,
// que se ejecutan antes o después de los manejadores de rutas.
package middleware

import (
	"net/http"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// UserLoader carga un usuario por su ID (implementado por services.AuthService).
type UserLoader interface {
	GetUserByID(userID int64) (*models.User, error)
}

// AuthMiddleware verifica que la solicitud incluya un token JWT válido.
// Acepta el token en el encabezado Authorization ("Bearer <token>")
// o como query string (?token=<token>), necesario para WebSocket.
// Si es válido, inyecta el usuario autenticado en el contexto.
func AuthMiddleware(users UserLoader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.URL.Query().Get("token")
			if authHeader := r.Header.Get("Authorization"); authHeader != "" {
				var found bool
				tokenString, found = strings.CutPrefix(authHeader, "Bearer ")
				if !found || tokenString == "" {
					response.Unauthorized(w, "Formato de token inválido")
					return
				}
			}
			if tokenString == "" {
				response.Unauthorized(w, "Token de autorización requerido")
				return
			}

			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				response.Unauthorized(w, "Token inválido")
				return
			}

			user, err := users.GetUserByID(claims.UserID)
			if err != nil {
				response.Unauthorized(w, "Usuario no encontrado")
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithUser(r.Context(), user)))
		})
	}
}

// CORSMiddleware agrega los encabezados CORS a todas las respuestas
// y responde directamente a las solicitudes de preflight (OPTIONS).
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
