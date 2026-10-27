// Paquete routes configura y registra todas las rutas HTTP de la aplicación.
package routes

import (
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/controllers"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/repository"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/services"
	"github.com/A4GOD-AMHG/LoveApp-Backend/middleware"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes construye las dependencias (composition root) y registra
// todas las rutas públicas y protegidas. CORS envuelve al router completo
// para que los preflight OPTIONS respondan aunque la ruta no declare ese método.
func SetupRoutes(hub *websocket.Hub) http.Handler {
	userRepo := repository.NewUserRepository()
	messageRepo := repository.NewMessageRepository()
	deviceRepo := repository.NewDevicePushTokenRepository()

	authService := services.NewAuthService(userRepo, messageRepo)
	messageService := services.NewMessageService(messageRepo, userRepo, deviceRepo, services.NewPushService(), hub)

	authController := controllers.NewAuthController(authService)
	healthController := controllers.NewHealthController()
	todoController := controllers.NewTodoController(services.NewTodoService(repository.NewTodoRepository()))
	messageController := controllers.NewMessageController(messageService, hub)
	deviceController := controllers.NewDeviceController(services.NewDeviceService(deviceRepo))

	r := mux.NewRouter()

	// Rutas públicas
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/health", healthController.Check).Methods("GET")
	r.HandleFunc("/auth/login", authController.Login).Methods("POST")

	// Rutas protegidas: requieren token JWT válido
	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("/auth/change-password", authController.ChangePassword).Methods("POST")
	protected.HandleFunc("/auth/logout", authController.Logout).Methods("POST")

	protected.HandleFunc("/todos", todoController.CreateTodo).Methods("POST")
	protected.HandleFunc("/todos", todoController.GetTodos).Methods("GET")
	protected.HandleFunc("/todos/{id}", todoController.DeleteTodo).Methods("DELETE")
	protected.HandleFunc("/todos/{id}", todoController.UpdateTodoStatus).Methods("PATCH")
	protected.HandleFunc("/todos/{id}", todoController.UpdateTodo).Methods("PUT")

	protected.HandleFunc("/ws", messageController.ServeWS)

	protected.HandleFunc("/messages", messageController.SendMessage).Methods("POST")
	protected.HandleFunc("/messages/conversation", messageController.GetConversation).Methods("GET")
	protected.HandleFunc("/messages/unread-count", messageController.GetUnreadCount).Methods("GET")
	protected.HandleFunc("/messages/{id}", messageController.EditMessage).Methods("PUT")
	protected.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	protected.HandleFunc("/messages/{id}/read", messageController.MarkAsRead).Methods("PATCH")
	protected.HandleFunc("/messages/{id}/delivered", messageController.MarkAsDelivered).Methods("PATCH")

	protected.HandleFunc("/devices/push-token", deviceController.RegisterPushToken).Methods("POST")
	protected.HandleFunc("/devices/push-token", deviceController.DeletePushToken).Methods("DELETE")

	return middleware.CORSMiddleware(r)
}
