package models

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Error message"`    // Error message
	Message string `json:"message" example:"Detailed error"` // Detailed error message
	Code    int    `json:"code" example:"400"`               // Error code
}

// SuccessResponse represents a generic success response
// @Description Success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operación exitosa"` // Success message
	Data    interface{} `json:"data,omitempty"`                      // Optional data
}

// HealthResponse represents the health check response
// @Description Health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`                    // Service status
	Message string `json:"message" example:"Servicio funcionando"` // Status message
}
