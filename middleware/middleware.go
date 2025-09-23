package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/controller"
	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/util"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		var tok string
		fmt.Sscanf(header, "Bearer %s", &tok)
		if tok == "" {
			util.JsonResp(w, 401, map[string]string{"error": "invalid token"})
			return
		}
		token, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return config.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			util.JsonResp(w, 401, map[string]string{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		sub := int64(claims["sub"].(float64))
		user, err := controller.FindUserByID(sub)
		if err != nil {
			util.JsonResp(w, 401, map[string]string{"error": "invalid token user"})
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
