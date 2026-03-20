// Paquete config gestiona la configuración global de la aplicación,
// leyendo valores desde variables de entorno con soporte para valores predeterminados.
package config

import (
	"os"
)

// Config agrupa todas las configuraciones de la aplicación.
type Config struct {
	Database DatabaseConfig // Configuración de la base de datos
	JWT      JWTConfig      // Configuración de autenticación JWT
	Server   ServerConfig   // Configuración del servidor HTTP
	Push     PushConfig     // Configuración para push notifications
}

// DatabaseConfig contiene la ruta al archivo de la base de datos SQLite.
type DatabaseConfig struct {
	Path string // Ruta al archivo .db de SQLite
}

// JWTConfig contiene el secreto utilizado para firmar y verificar tokens JWT.
type JWTConfig struct {
	Secret []byte // Clave secreta en bytes para firmar tokens JWT
}

// ServerConfig contiene la configuración del servidor HTTP.
type ServerConfig struct {
	Port string // Puerto en el que escucha el servidor (ej. "8080")
}

// PushConfig contiene la configuración del proveedor de notificaciones push.
type PushConfig struct {
	CredentialsFile string // Ruta al JSON de la service account de Firebase
}

// AppConfig es la instancia global de configuración, accesible desde todo el proyecto.
var AppConfig *Config

// LoadConfig carga la configuración desde las variables de entorno.
// Si una variable no está definida, se usa el valor predeterminado indicado.
func LoadConfig() error {
	AppConfig = &Config{
		Database: DatabaseConfig{
			// Ruta predeterminada al archivo SQLite si DB_PATH no está definido
			Path: getEnv("DB_PATH", "./data/loveapp.db"),
		},
		JWT: JWTConfig{
			// Secreto JWT: debe cambiarse en producción mediante la variable JWT_SECRET
			Secret: []byte(getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")),
		},
		Server: ServerConfig{
			// Puerto del servidor: predeterminado 8080 si SERVER_PORT no está definido
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Push: PushConfig{
			CredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", "loveapp-aa-firebase-adminsdk-fbsvc-ce92554680.json"),
		},
	}
	return nil
}

// GetDatabasePath retorna la ruta al archivo de la base de datos SQLite.
func (c *Config) GetDatabasePath() string {
	return c.Database.Path
}

// GetServerPort retorna el puerto del servidor con el prefijo ":" requerido por net/http.
// Si el puerto está vacío, retorna ":8080" como valor predeterminado.
// Si ya tiene el prefijo ":", lo retorna tal cual.
func (c *Config) GetServerPort() string {
	if c.Server.Port == "" {
		return ":8080"
	}
	if c.Server.Port[0] != ':' {
		return ":" + c.Server.Port
	}
	return c.Server.Port
}

// getEnv lee una variable de entorno por su clave.
// Si no existe o está vacía, retorna el valor de respaldo (fallback).
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
