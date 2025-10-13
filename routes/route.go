package routes

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/handlers"
	"github.com/A4GOD-AMHG/LoveApp-Backend/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(r *mux.Router) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	todoHandler := handlers.NewTodoHandler()
	healthHandler := handlers.NewHealthHandler()
	
	// Apply global middleware
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)
	
	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	// Health check endpoint
	r.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	
	// Authentication routes (public)
	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	
	// Protected routes
	protectedRoutes := r.NewRoute().Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)
	
	// Authentication protected routes
	protectedRoutes.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods("POST")
	
	// Todo routes
	protectedRoutes.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	protectedRoutes.HandleFunc("/todos", todoHandler.ListTodos).Methods("GET")
	protectedRoutes.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")
	protectedRoutes.HandleFunc("/todos/{id}/complete", todoHandler.CompleteTodo).Methods("POST")
}
