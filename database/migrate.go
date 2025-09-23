package database

func Migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		creator_id INTEGER NOT NULL,
		completed_anyel INTEGER NOT NULL DEFAULT 0,
		completed_alexis INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	for _, s := range stmts {
		if _, err := Db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
