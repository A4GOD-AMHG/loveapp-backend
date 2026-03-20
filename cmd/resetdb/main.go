package main

import (
	"log"
	"os"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	appdb "github.com/A4GOD-AMHG/LoveApp-Backend/database"
	pkgdb "github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	dbPath := config.AppConfig.GetDatabasePath()
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error al eliminar la base de datos %s: %v", dbPath, err)
	}

	if err := pkgdb.InitDB(); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer pkgdb.CloseDB()

	if err := appdb.Migrate(); err != nil {
		log.Fatalf("Error al ejecutar las migraciones: %v", err)
	}

	if err := appdb.Seed(); err != nil {
		log.Fatalf("Error al sembrar la base de datos: %v", err)
	}

	log.Printf("Base de datos reiniciada correctamente en %s", dbPath)
	log.Printf("Usuarios sembrados: anyel / password, alexis / password")
}
