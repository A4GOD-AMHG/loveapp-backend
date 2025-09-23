package database

import "golang.org/x/crypto/bcrypt"

func seed() error {
	users := []string{"anyel", "alexis"}
	for _, u := range users {
		var id int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", u).Scan(&id)
		if err == nil {
			continue
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", u, string(hash))
		if err != nil {
			return err
		}
	}
	return nil
}
