// Tests unitarios para el paquete config — carga de configuración y helpers.
package config

import (
	"os"
	"testing"
)

// TestLoadConfig_ValoresPredeterminados verifica que LoadConfig use los valores por defecto
// cuando no hay variables de entorno definidas.
func TestLoadConfig_ValoresPredeterminados(t *testing.T) {
	// Limpiar variables de entorno para asegurar valores predeterminados
	os.Unsetenv("DB_PATH")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SERVER_PORT")

	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig falló: %v", err)
	}

	if AppConfig.Database.Path != "./data/loveapp.db" {
		t.Errorf("DB path esperado './data/loveapp.db', se obtuvo '%s'", AppConfig.Database.Path)
	}
	if AppConfig.Server.Port != "8080" {
		t.Errorf("Puerto esperado '8080', se obtuvo '%s'", AppConfig.Server.Port)
	}
	if len(AppConfig.JWT.Secret) == 0 {
		t.Error("JWT Secret no debe estar vacío")
	}
}

// TestLoadConfig_DesdeVariablesDeEntorno verifica que LoadConfig lea correctamente
// las variables de entorno cuando están definidas.
func TestLoadConfig_DesdeVariablesDeEntorno(t *testing.T) {
	os.Setenv("DB_PATH", "/tmp/test.db")
	os.Setenv("JWT_SECRET", "mi-secreto-test")
	os.Setenv("SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("DB_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SERVER_PORT")
	}()

	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig falló: %v", err)
	}

	if AppConfig.Database.Path != "/tmp/test.db" {
		t.Errorf("DB path esperado '/tmp/test.db', se obtuvo '%s'", AppConfig.Database.Path)
	}
	if string(AppConfig.JWT.Secret) != "mi-secreto-test" {
		t.Errorf("JWT Secret esperado 'mi-secreto-test', se obtuvo '%s'", string(AppConfig.JWT.Secret))
	}
	if AppConfig.Server.Port != "9090" {
		t.Errorf("Puerto esperado '9090', se obtuvo '%s'", AppConfig.Server.Port)
	}
}

// TestGetServerPort_SinPrefijo verifica que se añada ":" al puerto si no lo tiene.
func TestGetServerPort_SinPrefijo(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Port: "8080"}}
	if puerto := cfg.GetServerPort(); puerto != ":8080" {
		t.Errorf("esperado ':8080', se obtuvo '%s'", puerto)
	}
}

// TestGetServerPort_ConPrefijo verifica que el puerto con ":" se retorne sin modificar.
func TestGetServerPort_ConPrefijo(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Port: ":8080"}}
	if puerto := cfg.GetServerPort(); puerto != ":8080" {
		t.Errorf("esperado ':8080', se obtuvo '%s'", puerto)
	}
}

// TestGetServerPort_Vacio verifica que un puerto vacío retorne ":8080" por defecto.
func TestGetServerPort_Vacio(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Port: ""}}
	if puerto := cfg.GetServerPort(); puerto != ":8080" {
		t.Errorf("esperado ':8080', se obtuvo '%s'", puerto)
	}
}

// TestGetDatabasePath verifica que GetDatabasePath retorne la ruta configurada.
func TestGetDatabasePath(t *testing.T) {
	cfg := &Config{Database: DatabaseConfig{Path: "./data/test.db"}}
	if ruta := cfg.GetDatabasePath(); ruta != "./data/test.db" {
		t.Errorf("ruta esperada './data/test.db', se obtuvo '%s'", ruta)
	}
}
