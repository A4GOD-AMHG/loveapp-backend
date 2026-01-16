package database

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

func Migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			name TEXT,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

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

		`CREATE INDEX IF NOT EXISTS idx_todos_creator_id ON todos(creator_id);`,
		`CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed_anyel, completed_alexis);`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,

		`CREATE TRIGGER IF NOT EXISTS update_users_updated_at
		AFTER UPDATE ON users
		FOR EACH ROW
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

		`CREATE TRIGGER IF NOT EXISTS update_todos_updated_at
		AFTER UPDATE ON todos
		FOR EACH ROW
		BEGIN
			UPDATE todos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

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

		`CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);`,

		`CREATE TRIGGER IF NOT EXISTS update_messages_updated_at
		AFTER UPDATE ON messages
		FOR EACH ROW
		BEGIN
			UPDATE messages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,
	}

	for _, stmt := range stmts {
		if _, err := database.DB.Exec(stmt); err != nil {
			return err
		}
	}

	_, _ = database.DB.Exec("ALTER TABLE users ADD COLUMN name TEXT")

	return nil
}
