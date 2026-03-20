// Paquete models define las estructuras de datos utilizadas en toda la aplicación.
package models

// ErrorResponse representa una respuesta de error estándar de la API.
// @Description Respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Mensaje de error"`   // Tipo o categoría del error
	Message string `json:"message" example:"Error detallado"`  // Descripción detallada del error para el usuario
	Code    int    `json:"code" example:"400"`                 // Código de estado HTTP asociado al error
}

// SuccessResponse representa una respuesta genérica de éxito de la API.
// @Description Respuesta de éxito
type SuccessResponse struct {
	Message string      `json:"message" example:"Operación exitosa"` // Mensaje descriptivo del resultado exitoso
	Data    interface{} `json:"data,omitempty"`                      // Datos opcionales retornados junto al mensaje
}

// HealthResponse representa la respuesta del endpoint de verificación de estado del servidor.
// @Description Respuesta de verificación de estado
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`                    // Estado actual del servicio (ej. "ok")
	Message string `json:"message" example:"Servicio funcionando"` // Mensaje descriptivo del estado del servicio
}
