// Tests unitarios para el paquete auth — generación y validación de tokens JWT.
package auth

import (
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

// setupConfig inicializa la configuración mínima necesaria para los tests de JWT.
func setupConfig() {
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret: []byte("clave-secreta-de-prueba"),
		},
	}
}

// TestGenerateToken verifica que GenerateToken produzca un token no vacío para un usuario válido.
func TestGenerateToken(t *testing.T) {
	setupConfig()

	user := &models.User{ID: 1, Username: "anyel"}
	token, err := GenerateToken(user)

	if err != nil {
		t.Fatalf("se esperaba nil, se obtuvo error: %v", err)
	}
	if token == "" {
		t.Fatal("se esperaba un token no vacío")
	}
}

// TestValidateToken_TokenValido verifica que un token generado pueda ser validado correctamente
// y que los claims contengan los valores originales del usuario.
func TestValidateToken_TokenValido(t *testing.T) {
	setupConfig()

	user := &models.User{ID: 42, Username: "alexis"}
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("error al generar token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("error al validar token válido: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID esperado 42, se obtuvo %d", claims.UserID)
	}
	if claims.Username != "alexis" {
		t.Errorf("Username esperado 'alexis', se obtuvo '%s'", claims.Username)
	}
}

// TestValidateToken_TokenInvalido verifica que un token malformado o falso sea rechazado.
func TestValidateToken_TokenInvalido(t *testing.T) {
	setupConfig()

	_, err := ValidateToken("esto.no.es.un.token.valido")
	if err == nil {
		t.Fatal("se esperaba error para token inválido, se obtuvo nil")
	}
}

// TestValidateToken_TokenConFirmaDistinta verifica que un token firmado con otra clave sea rechazado.
func TestValidateToken_TokenConFirmaDistinta(t *testing.T) {
	// Generar token con clave correcta
	setupConfig()
	user := &models.User{ID: 1, Username: "anyel"}
	token, _ := GenerateToken(user)

	// Cambiar la clave secreta para simular una firma inválida
	config.AppConfig.JWT.Secret = []byte("clave-diferente-invalida")

	_, err := ValidateToken(token)
	if err == nil {
		t.Fatal("se esperaba error para token con firma inválida, se obtuvo nil")
	}
}

// TestGenerateToken_IDsDistintosProducenTokensDistintos verifica que usuarios diferentes
// generen tokens distintos.
func TestGenerateToken_IDsDistintosProducenTokensDistintos(t *testing.T) {
	setupConfig()

	user1 := &models.User{ID: 1, Username: "anyel"}
	user2 := &models.User{ID: 2, Username: "alexis"}

	token1, _ := GenerateToken(user1)
	token2, _ := GenerateToken(user2)

	if token1 == token2 {
		t.Error("tokens para usuarios distintos no deberían ser iguales")
	}
}
