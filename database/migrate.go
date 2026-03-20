// Paquete database contiene las funciones de migración y sembrado de la base de datos.
package database

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// Migrate ejecuta todas las sentencias DDL necesarias para crear las tablas,
// índices y triggers de la base de datos si aún no existen.
// Es seguro ejecutarla múltiples veces (usa IF NOT EXISTS).
func Migrate() error {
	stmts := []string{
		// Tabla de usuarios: almacena las cuentas de los dos usuarios de la app
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			name TEXT,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// Tabla de tareas: cada tarea tiene estado de completado independiente por usuario
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
		);`,

		// Índice para búsquedas de tareas por creador
		`CREATE INDEX IF NOT EXISTS idx_todos_creator_id ON todos(creator_id);`,
		// Índice para filtrar tareas por estado de completado
		`CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed_anyel, completed_alexis);`,
		// Índice para búsquedas de usuarios por nombre de usuario
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,

		// Trigger: actualiza automáticamente updated_at en usuarios al modificar un registro
		`CREATE TRIGGER IF NOT EXISTS update_users_updated_at
		AFTER UPDATE ON users
		FOR EACH ROW
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

		// Trigger: actualiza automáticamente updated_at en tareas al modificar un registro
		`CREATE TRIGGER IF NOT EXISTS update_todos_updated_at
		AFTER UPDATE ON todos
		FOR EACH ROW
		BEGIN
			UPDATE todos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

		// Tabla de mensajes: almacena la conversación entre los dos usuarios
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
		);`,

		// Índice para búsquedas de mensajes por remitente
		`CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);`,
		// Índice para búsquedas de mensajes por destinatario
		`CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);`,

		// Tabla de dispositivos push: almacena el token actual por usuario y dispositivo.
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
		);`,

		// Índices para búsquedas de tokens por usuario y token.
		`CREATE INDEX IF NOT EXISTS idx_device_push_tokens_user_id ON device_push_tokens(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_device_push_tokens_push_token ON device_push_tokens(push_token);`,

		// Trigger: actualiza automáticamente updated_at en mensajes al modificar un registro
		`CREATE TRIGGER IF NOT EXISTS update_messages_updated_at
		AFTER UPDATE ON messages
		FOR EACH ROW
		BEGIN
			UPDATE messages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

		// Trigger: actualiza automáticamente updated_at en tokens push al modificar un registro.
		`CREATE TRIGGER IF NOT EXISTS update_device_push_tokens_updated_at
		AFTER UPDATE ON device_push_tokens
		FOR EACH ROW
		BEGIN
			UPDATE device_push_tokens SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,
	}

	// Ejecutar cada sentencia DDL en orden
	for _, stmt := range stmts {
		if _, err := database.DB.Exec(stmt); err != nil {
			return err
		}
	}

	// Agregar columna name a users si no existe (migración incremental segura)
	_, _ = database.DB.Exec("ALTER TABLE users ADD COLUMN name TEXT")

	return nil
}
