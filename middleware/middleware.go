// Paquete middleware contiene los middlewares HTTP de la aplicación,
// que se ejecutan antes o después de los manejadores de rutas.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// AuthMiddleware verifica que la solicitud entrante incluya un token JWT válido.
// Acepta el token en el encabezado Authorization (formato "Bearer <token>")
// o como parámetro de query string (?token=<token>).
// Si el token es válido, carga el usuario autenticado en el contexto de la solicitud
// para que esté disponible en los manejadores protegidos.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Intentar leer el token desde el encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// El encabezado debe tener el formato "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "Formato de token inválido")
				return
			}
			tokenString = parts[1]
		} else {
			// Fallback: intentar leer el token desde la query string (?token=...)
			tokenString = r.URL.Query().Get("token")
			if tokenString == "" {
				response.Unauthorized(w, "Token de autorización requerido")
				return
			}
		}

		// Validar el token JWT y extraer los claims (user_id, username)
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Token inválido")
			return
		}

		// Cargar el usuario completo desde la base de datos usando el ID del token
		authService := services.NewAuthService()
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			response.Unauthorized(w, "Usuario no encontrado")
			return
		}

		// Inyectar el usuario y sus datos en el contexto de la solicitud
		// para que los manejadores downstream puedan accederlos
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "username", user.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware agrega los encabezados de Cross-Origin Resource Sharing (CORS)
// a todas las respuestas, permitiendo que el frontend acceda a la API desde cualquier origen.
// Responde directamente a las solicitudes de preflight (OPTIONS) con HTTP 200.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir solicitudes desde cualquier origen
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Métodos HTTP permitidos en solicitudes cross-origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Encabezados permitidos en solicitudes cross-origin
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Responder inmediatamente a solicitudes de preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware es un middleware de registro de solicitudes HTTP.
// Actualmente delega directamente al siguiente handler (placeholder para logging futuro).
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
