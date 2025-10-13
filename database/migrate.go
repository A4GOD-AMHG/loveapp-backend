package database

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
)

// Migrate runs database migrations
func Migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
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

		// Create trigger for updating updated_at timestamp for users
		`CREATE TRIGGER IF NOT EXISTS update_users_updated_at 
		AFTER UPDATE ON users
		FOR EACH ROW
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,

		// Create trigger for updating updated_at timestamp for todos
		`CREATE TRIGGER IF NOT EXISTS update_todos_updated_at 
		AFTER UPDATE ON todos
		FOR EACH ROW
		BEGIN
			UPDATE todos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,
	}

	for _, stmt := range stmts {
		if _, err := database.DB.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
