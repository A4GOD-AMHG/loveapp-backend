// @title API Backend de LoveApp
// @version 1.0
// @description Esta es una API de gestión de tareas para parejas. Permite a dos usuarios (anyel y alexis) crear, gestionar y completar tareas juntos.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escriba "Bearer" seguido de un espacio y el token JWT.
package main

import (
	"log"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/docs"
	pkg_database "github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/A4GOD-AMHG/LoveApp-Backend/routes"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	if err := pkg_database.InitDB(); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer pkg_database.CloseDB()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Error al ejecutar las migraciones: %v", err)
	}

	if err := database.Seed(); err != nil {
		log.Fatalf("Error al sembrar la base de datos: %v", err)
	}

	docs.SwaggerInfo.Title = "API de LoveApp"
	docs.SwaggerInfo.Description = "Esta es la API para la aplicación LoveApp."
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http"}

	hub := websocket.NewHub()
	go hub.Run()

	r := routes.SetupRoutes(hub)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := config.AppConfig.GetServerPort()
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
