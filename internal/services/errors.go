package services

import "errors"

// Errores de dominio: los controladores los distinguen con errors.Is
// en lugar de comparar textos.
var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrExpired          = errors.New("edit window expired")
	ErrInvalidInput     = errors.New("invalid input")
	ErrNoReceiver       = errors.New("other user not found")
	ErrAlreadyCompleted = errors.New("todo already completed by both")
)
