// Paquete testhelpers provee utilidades compartidas para los tests de la aplicación.
// Incluye la configuración de una base de datos SQLite en memoria para tests aislados.
package testhelpers

import (
	"database/sql"
	"testing"

	pkg_db "github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB inicializa una base de datos SQLite en memoria para tests.
// Crea todas las tablas necesarias y retorna una función de limpieza para usar con defer.
// Ejemplo de uso:
//
//	cleanup := testhelpers.SetupTestDB(t)
//	defer cleanup()
func SetupTestDB(t *testing.T) func() {
	t.Helper()

	var err error
	// Usar SQLite en memoria con cache compartida para que todos los repositorios
	// accedan a la misma instancia durante el test
	pkg_db.DB, err = sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error al abrir DB en memoria: %v", err)
	}

	// Limitar a 1 conexión para evitar conflictos con SQLite en modo in-memory
	pkg_db.DB.SetMaxOpenConns(1)

	// Habilitar claves foráneas
	if _, err = pkg_db.DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("error al habilitar claves foráneas: %v", err)
	}

	// Crear esquema de tablas necesarias para los tests
	schema := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			name TEXT,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			creator_id INTEGER NOT NULL,
			completed_anyel INTEGER NOT NULL DEFAULT 0,
			completed_alexis INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender_id INTEGER NOT NULL,
			receiver_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'sent',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS device_push_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			platform TEXT NOT NULL,
			push_token TEXT NOT NULL UNIQUE,
			device_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, platform, device_name),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, stmt := range schema {
		if _, err := pkg_db.DB.Exec(stmt); err != nil {
			t.Fatalf("error al crear esquema: %v\nSQL: %s", err, stmt)
		}
	}

	// Función de limpieza: cerrar la conexión al finalizar el test
	return func() {
		pkg_db.DB.Close()
	}
}

// InsertTestUser inserta un usuario de prueba en la base de datos y retorna su ID.
func InsertTestUser(t *testing.T, username, name, password string) int64 {
	t.Helper()
	res, err := pkg_db.DB.Exec(
		"INSERT INTO users (username, name, password) VALUES (?, ?, ?)",
		username, name, password,
	)
	if err != nil {
		t.Fatalf("error al insertar usuario de prueba '%s': %v", username, err)
	}
	id, _ := res.LastInsertId()
	return id
}
