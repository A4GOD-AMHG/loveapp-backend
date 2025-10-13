// Package main LoveApp Backend API
// @title LoveApp Backend API
// @version 1.0
// @description This is a todo management API for couples. It allows two users (anyel and alexis) to create, manage, and complete todos together.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@loveapp.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	_ "github.com/A4GOD-AMHG/LoveApp-Backend/docs"
	pkgdatabase "github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/routes"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize configuration
	config.InitConfig()
	log.Println("Configuration initialized")

	// Initialize database
	if err := pkgdatabase.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := pkgdatabase.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Seed database
	if err := database.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeding completed")

	// Initialize router
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// Start server
	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📚 Swagger documentation available at http://localhost%s/swagger/index.html", addr)
	log.Printf("🏥 Health check available at http://localhost%s/health", addr)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(addr, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-c
	log.Println("🛑 Server shutting down gracefully...")
}
