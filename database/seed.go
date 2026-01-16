package database

import (
	"log"

	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func Seed() error {
	type userSeed struct {
		Name     string
		Username string
	}
	users := []userSeed{
		{Name: "Anyel", Username: "anyel"},
		{Name: "Alexis", Username: "alexis"},
	}
	for _, u := range users {
		var id int
		err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", u.Username).Scan(&id)
		if err == nil {
			log.Printf("El usuario %s ya existe, omitiendo...", u.Username)
			continue
		}
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = database.DB.Exec("INSERT INTO users (username, name, password) VALUES ($1, $2, $3)", u.Username, u.Name, string(hash))
		if err != nil {
			return err
		}
		log.Printf("Usuario creado: %s (%s) con contraseña: password", u.Username, u.Name)
	}
	log.Println("Sembrado de base de datos completado exitosamente")
	return nil
}
