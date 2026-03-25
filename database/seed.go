// Paquete database contiene las funciones de migración y sembrado de la base de datos.
package database

import (
	"log"

	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// Seed inserta los usuarios iniciales de la aplicación (anyel y alexis) si aún no existen.
// Si un usuario ya está registrado, se omite sin error.
// Las contraseñas iniciales son: anyel/anyel y alexis/alexis.
func Seed() error {
	// Definición de los usuarios iniciales que deben existir en el sistema
	type userSeed struct {
		Name     string
		Username string
		Password string
	}
	users := []userSeed{
		{Name: "Anyel", Username: "anyel", Password: "anyel"},
		{Name: "Alexis", Username: "alexis", Password: "alexis"},
	}

	for _, u := range users {
		var id int
		// Verificar si el usuario ya existe en la base de datos
		err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", u.Username).Scan(&id)
		if err == nil {
			// El usuario ya existe, saltar sin error
			log.Printf("El usuario %s ya existe, omitiendo...", u.Username)
			continue
		}

		// Generar hash seguro de la contraseña inicial del usuario
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Insertar el nuevo usuario en la base de datos
		_, err = database.DB.Exec("INSERT INTO users (username, name, password) VALUES ($1, $2, $3)", u.Username, u.Name, string(hash))
		if err != nil {
			return err
		}
		log.Printf("Usuario creado: %s (%s) con contraseña: %s", u.Username, u.Name, u.Password)
	}

	log.Println("Sembrado de base de datos completado exitosamente")
	return nil
}
