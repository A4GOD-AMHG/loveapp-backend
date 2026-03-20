// Paquete database gestiona la conexión a la base de datos SQLite de la aplicación.
// Expone una instancia global DB utilizada por todos los repositorios.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	_ "github.com/mattn/go-sqlite3" // Driver SQLite3 registrado como efecto secundario
)

// DB es la instancia global de conexión a la base de datos SQLite.
// Todos los repositorios acceden a la base de datos a través de esta variable.
var DB *sql.DB

// InitDB inicializa la conexión a la base de datos SQLite.
// Crea el directorio de la base de datos si no existe, abre la conexión,
// verifica su disponibilidad con un ping, configura el pool de conexiones
// y habilita el soporte de claves foráneas (PRAGMA foreign_keys = ON).
// Debe llamarse una sola vez al inicio de la aplicación.
func InitDB() error {
	var err error

	// Obtener la ruta del archivo de base de datos desde la configuración
	dbPath := config.AppConfig.GetDatabasePath()

	// Crear el directorio padre si no existe (ej. ./data/)
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio de la base de datos: %w", err)
	}

	// Abrir la conexión al archivo SQLite (lo crea si no existe)
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Verificar que la conexión esté activa y funcional
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	// SQLite no soporta concurrencia con múltiples escritores; limitar a 1 conexión activa
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	// Habilitar el cumplimiento de claves foráneas (deshabilitado por defecto en SQLite)
	if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("error al habilitar claves foráneas: %w", err)
	}

	log.Printf("Conexión a base de datos SQLite establecida exitosamente en: %s", dbPath)
	return nil
}

// CloseDB cierra la conexión a la base de datos de forma segura.
// Debe llamarse mediante defer en la función main para liberar recursos al cerrar la app.
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
