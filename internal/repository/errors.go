package repository

import "errors"

// ErrNotFound indica que el registro solicitado no existe.
var ErrNotFound = errors.New("record not found")
