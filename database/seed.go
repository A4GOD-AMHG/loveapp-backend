package database

import (
	"log"

	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// Seed creates initial data in the database
func Seed() error {
	users := []string{"anyel", "alexis"}
	
	for _, username := range users {
		var id int
		err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
		if err == nil {
			log.Printf("User %s already exists, skipping...", username)
			continue
		}
		
		// Hash the default password
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		
		// Insert the user
		_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hash))
		if err != nil {
			return err
		}
		
		log.Printf("Created user: %s with password: password", username)
	}
	
	log.Println("Database seeding completed successfully")
	return nil
}
