// @title API Backend de LoveApp
// @version 1.0
// @description Esta es una API de gestión de tareas para parejas. Permite a dos usuarios (anyel y alexis) crear, gestionar y completar tareas juntos.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escriba "Bearer" seguido de un espacio y el token JWT.

// Paquete principal — punto de entrada de la aplicación LoveApp Backend.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/docs"
	pkg_database "github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/A4GOD-AMHG/LoveApp-Backend/routes"
	httpSwagger "github.com/swaggo/http-swagger"
)

// main es la función de arranque del servidor.
// Carga la configuración, inicializa la base de datos, ejecuta las migraciones
// y el sembrado inicial, y luego levanta el servidor HTTP en el puerto configurado.
func main() {
	// Cargar variables de entorno y configuración de la aplicación
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Inicializar la conexión a la base de datos SQLite
	if err := pkg_database.InitDB(); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer pkg_database.CloseDB()

	// Ejecutar migraciones para crear tablas e índices si no existen
	if err := database.Migrate(); err != nil {
		log.Fatalf("Error al ejecutar las migraciones: %v", err)
	}

	// Sembrar usuarios iniciales (anyel y alexis) si no existen
	if err := database.Seed(); err != nil {
		log.Fatalf("Error al sembrar la base de datos: %v", err)
	}

	// Permite reutilizar el arranque para reiniciar DB sin levantar el servidor HTTP.
	if os.Getenv("LOVEAPP_RESET_ONLY") == "1" {
		log.Printf("Base de datos migrada y sembrada; finalizando por LOVEAPP_RESET_ONLY=1")
		return
	}

	// Configurar metadatos de Swagger para la documentación interactiva
	docs.SwaggerInfo.Title = "API de LoveApp"
	docs.SwaggerInfo.Description = "Esta es la API para la aplicación LoveApp."
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http"}

	// Crear e iniciar el hub de WebSocket en una goroutine separada
	hub := websocket.NewHub()
	go hub.Run()

	// Configurar el enrutador con todas las rutas de la aplicación
	r := routes.SetupRoutes(hub)

	// Registrar el handler de Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Obtener el puerto del servidor desde la configuración
	port := config.AppConfig.GetServerPort()
	log.Printf("Servidor iniciado en el puerto %s", port)

	// Iniciar el servidor HTTP y escuchar conexiones entrantes
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
