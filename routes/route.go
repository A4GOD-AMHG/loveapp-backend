// Paquete routes configura y registra todas las rutas HTTP de la aplicación.
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

// SetupRoutes inicializa el enrutador Gorilla Mux, registra los middlewares globales
// y define todas las rutas públicas y protegidas de la aplicación.
// Retorna el enrutador configurado listo para ser usado por el servidor HTTP.
func SetupRoutes(hub *websocket.Hub) *mux.Router {
	r := mux.NewRouter()

	// Instanciar los controladores
	authController := controllers.NewAuthController()
	healthController := controllers.NewHealthController()
	todoController := controllers.NewTodoController()

	// Aplicar middlewares globales a todas las rutas
	r.Use(middleware.CORSMiddleware)    // Soporte para solicitudes cross-origin
	r.Use(middleware.LoggingMiddleware) // Middleware de logging (extensible)

	// Ruta de documentación Swagger UI (pública)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/health", healthController.Check).Methods("GET")

	// Rutas públicas de autenticación (sin middleware de auth)
	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authController.Login).Methods("POST")

	// Subrouter protegido: todas las rutas aquí requieren token JWT válido
	protectedRoutes := r.NewRoute().Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	// Rutas de autenticación protegidas
	protectedRoutes.HandleFunc("/auth/change-password", authController.ChangePassword).Methods("POST")
	protectedRoutes.HandleFunc("/auth/logout", authController.Logout).Methods("POST")

	// Rutas de gestión de tareas (todos)
	protectedRoutes.HandleFunc("/todos", todoController.CreateTodo).Methods("POST")
	protectedRoutes.HandleFunc("/todos", todoController.GetTodos).Methods("GET")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.DeleteTodo).Methods("DELETE")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.UpdateTodoStatus).Methods("PATCH")
	protectedRoutes.HandleFunc("/todos/{id}", todoController.UpdateTodo).Methods("PUT")

	// Inicializar dependencias del sistema de mensajería
	messageRepo := repository.NewMessageRepository()
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDevicePushTokenRepository()
	pushService := services.NewPushService()
	messageService := services.NewMessageService(messageRepo, userRepo, deviceRepo, pushService, hub)
	messageController := controllers.NewMessageController(messageService, hub)
	deviceService := services.NewDeviceService(deviceRepo)
	deviceController := controllers.NewDeviceController(deviceService)

	// Ruta de WebSocket para comunicación en tiempo real
	protectedRoutes.HandleFunc("/ws", messageController.ServeWS)

	// Rutas de gestión de mensajes
	protectedRoutes.HandleFunc("/messages", messageController.SendMessage).Methods("POST")
	protectedRoutes.HandleFunc("/messages/{id}", messageController.EditMessage).Methods("PUT")
	protectedRoutes.HandleFunc("/messages/{id}/read", messageController.MarkAsRead).Methods("PATCH")
	protectedRoutes.HandleFunc("/messages/{id}/delivered", messageController.MarkAsDelivered).Methods("PATCH")
	protectedRoutes.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	protectedRoutes.HandleFunc("/messages/conversation", messageController.GetConversation).Methods("GET")
	protectedRoutes.HandleFunc("/messages/unread-count", messageController.GetUnreadCount).Methods("GET")
	protectedRoutes.HandleFunc("/devices/push-token", deviceController.RegisterPushToken).Methods("POST")
	protectedRoutes.HandleFunc("/devices/push-token", deviceController.DeletePushToken).Methods("DELETE")

	return r
}
