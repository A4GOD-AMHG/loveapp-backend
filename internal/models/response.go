package models

// ErrorResponse representa una respuesta de error
// @Description Respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Mensaje de error"`         // Mensaje de error
	Message string `json:"message" example:"Error detallado"`        // Mensaje de error detallado
	Code    int    `json:"code" example:"400"`                       // Código de error
}

// SuccessResponse representa una respuesta genérica de éxito
// @Description Respuesta de éxito
type SuccessResponse struct {
	Message string      `json:"message" example:"Operación exitosa"` // Mensaje de éxito
	Data    interface{} `json:"data,omitempty"`                      // Datos opcionales
}

// HealthResponse representa la respuesta de verificación de estado
// @Description Respuesta de verificación de estado
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`                          // Estado del servicio
	Message string `json:"message" example:"Servicio funcionando"`       // Mensaje de estado
}
