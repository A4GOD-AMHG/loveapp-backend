package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error

	dbPath := config.AppConfig.GetDatabasePath()

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio de la base de datos: %w", err)
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("error al habilitar claves foráneas: %w", err)
	}

	log.Printf("Conexión a base de datos SQLite establecida exitosamente en: %s", dbPath)
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
