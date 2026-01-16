package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Token de autorización requerido")
			return
		}
		
		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Formato de token inválido")
			return
		}
		
		tokenString := parts[1]
		
		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Token inválido")
			return
		}
		
		// Get user from database
		authService := services.NewAuthService()
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			response.Unauthorized(w, "Usuario no encontrado")
			return
		}
		
		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "username", user.Username)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// You can add logging logic here
		next.ServeHTTP(w, r)
	})
}
