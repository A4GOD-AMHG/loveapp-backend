package routes

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/controllers"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/middleware"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(hub *websocket.Hub) *mux.Router {
	r := mux.NewRouter()
	authController := controllers.NewAuthController()
	todoController := controllers.NewTodoController()

	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authController.Login).Methods("POST")

	protectedRoutes := r.NewRoute().Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	protectedRoutes.HandleFunc("/auth/change-password", authController.ChangePassword).Methods("POST")
	protectedRoutes.HandleFunc("/auth/logout", authController.Logout).Methods("POST")

	protectedRoutes.HandleFunc("/todos", todoController.CreateTodo).Methods("POST")
	protectedRoutes.HandleFunc("/todos", todoController.GetTodos).Methods("GET")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.DeleteTodo).Methods("DELETE")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.UpdateTodoStatus).Methods("PATCH")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.UpdateTodo).Methods("PUT")

	messageRepo := repository.NewMessageRepository()
	userRepo := repository.NewUserRepository()
	messageService := services.NewMessageService(messageRepo, userRepo, hub)
	messageController := controllers.NewMessageController(messageService, hub)

	protectedRoutes.HandleFunc("/ws", messageController.ServeWS)
	protectedRoutes.HandleFunc("/messages", messageController.SendMessage).Methods("POST")
	protectedRoutes.HandleFunc("/messages/{id}", messageController.EditMessage).Methods("PUT")
	protectedRoutes.HandleFunc("/messages/{id}/read", messageController.MarkAsRead).Methods("PATCH")
	protectedRoutes.HandleFunc("/messages/{id}/delivered", messageController.MarkAsDelivered).Methods("PATCH")
	protectedRoutes.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	protectedRoutes.HandleFunc("/messages/conversation", messageController.GetConversation).Methods("GET")

	return r
}
