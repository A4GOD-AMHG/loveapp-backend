// Paquete response provee funciones helper para construir y enviar respuestas HTTP
// con formato JSON estándar en toda la aplicación.
package response

import (
	"encoding/json"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

// JSON serializa el dato proporcionado como JSON y lo escribe en la respuesta HTTP
// con el código de estado indicado. Establece el Content-Type a "application/json".
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Success envía una respuesta HTTP 200 OK con un mensaje de éxito y datos opcionales.
func Success(w http.ResponseWriter, message string, data interface{}) {
	response := models.SuccessResponse{
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, response)
}

// Created envía una respuesta HTTP 201 Created con un mensaje de éxito y datos opcionales.
func Created(w http.ResponseWriter, message string, data interface{}) {
	response := models.SuccessResponse{
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusCreated, response)
}

// Error construye y envía una respuesta de error con el código de estado,
// mensaje para el usuario y texto del error indicados.
func Error(w http.ResponseWriter, statusCode int, message string, err string) {
	response := models.ErrorResponse{
		Error:   err,
		Message: message,
		Code:    statusCode,
	}
	JSON(w, statusCode, response)
}

// BadRequest envía una respuesta HTTP 400 Bad Request con el mensaje proporcionado.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message, "Solicitud inválida")
}

// Unauthorized envía una respuesta HTTP 401 Unauthorized con el mensaje proporcionado.
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, "No autorizado")
}

// Forbidden envía una respuesta HTTP 403 Forbidden con el mensaje proporcionado.
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message, "Acceso denegado")
}

// NotFound envía una respuesta HTTP 404 Not Found con el mensaje proporcionado.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, "No encontrado")
}

// InternalServerError envía una respuesta HTTP 500 Internal Server Error con el mensaje proporcionado.
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message, "Error interno del servidor")
}
