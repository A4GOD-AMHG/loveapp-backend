// Paquete config gestiona la configuración global de la aplicación,
// leyendo valores desde variables de entorno con soporte para valores predeterminados.
package config

import (
	"bufio"
	"os"
	"strings"
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
	Port string // Puerto en el que escucha el servidor (ej. "4418")
}

// PushConfig contiene la configuración del proveedor de notificaciones push.
type PushConfig struct {
	CredentialsFile         string // Ruta opcional al JSON de la service account de Firebase
	Type                    string // Tipo de credencial de Firebase
	ProjectID               string // Project ID de Firebase
	PrivateKeyID            string // ID de la llave privada
	PrivateKey              string // Llave privada PEM
	ClientEmail             string // Email de la service account
	ClientID                string // ID del cliente OAuth
	AuthURI                 string // URI de autenticación
	TokenURI                string // URI para obtener access tokens
	AuthProviderX509CertURL string // URL del cert provider
	ClientX509CertURL       string // URL del cert del cliente
	UniverseDomain          string // Dominio de Google APIs
}

// AppConfig es la instancia global de configuración, accesible desde todo el proyecto.
var AppConfig *Config

// LoadConfig carga la configuración desde las variables de entorno.
// Si una variable no está definida, se usa el valor predeterminado indicado.
func LoadConfig() error {
	loadDotEnv(".env")

	AppConfig = &Config{
		Database: DatabaseConfig{
			// Ruta predeterminada al archivo SQLite si DB_PATH no está definido
			Path: cleanEnvValue(getEnv("DB_PATH", "./data/loveapp.db")),
		},
		JWT: JWTConfig{
			// Secreto JWT: debe cambiarse en producción mediante la variable JWT_SECRET
			Secret: []byte(cleanEnvValue(getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"))),
		},
		Server: ServerConfig{
			// Puerto del servidor: predeterminado 4418 si SERVER_PORT no está definido
			Port: cleanEnvValue(getEnv("SERVER_PORT", "4418")),
		},
		Push: PushConfig{
			CredentialsFile:         cleanEnvValue(getEnv("FIREBASE_CREDENTIALS_FILE", "loveapp-aa-firebase-adminsdk-fbsvc-ce92554680.json")),
			Type:                    cleanEnvValue(getEnv("FIREBASE_TYPE", "")),
			ProjectID:               cleanEnvValue(getEnv("FIREBASE_PROJECT_ID", "")),
			PrivateKeyID:            cleanEnvValue(getEnv("FIREBASE_PRIVATE_KEY_ID", "")),
			PrivateKey:              strings.ReplaceAll(cleanEnvValue(getEnv("FIREBASE_PRIVATE_KEY", "")), `\n`, "\n"),
			ClientEmail:             cleanEnvValue(getEnv("FIREBASE_CLIENT_EMAIL", "")),
			ClientID:                cleanEnvValue(getEnv("FIREBASE_CLIENT_ID", "")),
			AuthURI:                 cleanEnvValue(getEnv("FIREBASE_AUTH_URI", "")),
			TokenURI:                cleanEnvValue(getEnv("FIREBASE_TOKEN_URI", "")),
			AuthProviderX509CertURL: cleanEnvValue(getEnv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL", "")),
			ClientX509CertURL:       cleanEnvValue(getEnv("FIREBASE_CLIENT_X509_CERT_URL", "")),
			UniverseDomain:          cleanEnvValue(getEnv("FIREBASE_UNIVERSE_DOMAIN", "")),
		},
	}
	return nil
}

// GetDatabasePath retorna la ruta al archivo de la base de datos SQLite.
func (c *Config) GetDatabasePath() string {
	return c.Database.Path
}

// GetServerPort retorna el puerto del servidor con el prefijo ":" requerido por net/http.
// Si el puerto está vacío, retorna ":4418" como valor predeterminado.
// Si ya tiene el prefijo ":", lo retorna tal cual.
func (c *Config) GetServerPort() string {
	if c.Server.Port == "" {
		return ":4418"
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

// loadDotEnv carga variables simples KEY=VALUE desde un archivo .env si existe.
// No sobreescribe variables ya presentes en el entorno del proceso.
func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if len(value) >= 2 {
			value = cleanEnvValue(value)
		}

		_ = os.Setenv(key, value)
	}
}

func cleanEnvValue(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			value = value[1 : len(value)-1]
		}
	}
	return value
}
