package response

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

// JSON sends a JSON response with the given status code and data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Success sends a success response
func Success(w http.ResponseWriter, message string, data interface{}) {
	response := models.SuccessResponse{
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, response)
}

// Created sends a created response
func Created(w http.ResponseWriter, message string, data interface{}) {
	response := models.SuccessResponse{
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusCreated, response)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, message string, err string) {
	response := models.ErrorResponse{
		Error:   err,
		Message: message,
		Code:    statusCode,
	}
	JSON(w, statusCode, response)
}

// BadRequest sends a bad request error response
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message, "Solicitud inválida")
}

// Unauthorized sends an unauthorized error response
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, "No autorizado")
}

// Forbidden sends a forbidden error response
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message, "Acceso denegado")
}

// NotFound sends a not found error response
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, "No encontrado")
}

// InternalServerError sends an internal server error response
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message, "Error interno del servidor")
}
