package auth

import (
	"context"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

// userKey es una clave privada para evitar colisiones en el contexto.
type userKey struct{}

// WithUser retorna un contexto que transporta al usuario autenticado.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// UserFrom extrae el usuario autenticado del contexto, si existe.
func UserFrom(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userKey{}).(*models.User)
	return user, ok && user != nil
}
