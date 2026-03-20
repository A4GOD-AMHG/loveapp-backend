// Tests unitarios para UserRepository — operaciones CRUD de usuarios.
package repository

import (
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
)

// TestFindByUsername_Encontrado verifica que FindByUsername retorne el usuario correcto.
func TestFindByUsername_Encontrado(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash123")

	repo := NewUserRepository()
	user, err := repo.FindByUsername("anyel")

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if user.Username != "anyel" {
		t.Errorf("username esperado 'anyel', se obtuvo '%s'", user.Username)
	}
	if user.Name != "Anyel" {
		t.Errorf("nombre esperado 'Anyel', se obtuvo '%s'", user.Name)
	}
	if user.Password != "hash123" {
		t.Error("la contraseña debe estar incluida en FindByUsername")
	}
}

// TestFindByUsername_NoExiste verifica que FindByUsername retorne error si el usuario no existe.
func TestFindByUsername_NoExiste(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewUserRepository()
	_, err := repo.FindByUsername("noexiste")

	if err == nil {
		t.Fatal("se esperaba error para usuario inexistente, se obtuvo nil")
	}
}

// TestFindByID_Encontrado verifica que FindByID retorne el usuario correcto sin contraseña.
func TestFindByID_Encontrado(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	id := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash456")

	repo := NewUserRepository()
	user, err := repo.FindByID(id)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if user.ID != id {
		t.Errorf("ID esperado %d, se obtuvo %d", id, user.ID)
	}
	if user.Username != "alexis" {
		t.Errorf("username esperado 'alexis', se obtuvo '%s'", user.Username)
	}
	// FindByID no debe retornar la contraseña
	if user.Password != "" {
		t.Error("FindByID no debe incluir la contraseña en la respuesta")
	}
}

// TestFindByID_NoExiste verifica que FindByID retorne error para un ID inexistente.
func TestFindByID_NoExiste(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewUserRepository()
	_, err := repo.FindByID(9999)

	if err == nil {
		t.Fatal("se esperaba error para ID inexistente, se obtuvo nil")
	}
}

// TestUpdatePassword verifica que UpdatePassword actualice el hash correctamente.
func TestUpdatePassword(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	id := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash_viejo")

	repo := NewUserRepository()
	err := repo.UpdatePassword(id, "hash_nuevo")

	if err != nil {
		t.Fatalf("error al actualizar contraseña: %v", err)
	}

	// Verificar que el nuevo hash esté guardado
	hash, err := repo.GetPasswordHash(id)
	if err != nil {
		t.Fatalf("error al obtener hash: %v", err)
	}
	if hash != "hash_nuevo" {
		t.Errorf("hash esperado 'hash_nuevo', se obtuvo '%s'", hash)
	}
}

// TestUpdatePassword_UsuarioInexistente verifica que UpdatePassword retorne error
// si el usuario no existe.
func TestUpdatePassword_UsuarioInexistente(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewUserRepository()
	err := repo.UpdatePassword(9999, "nuevo_hash")

	if err == nil {
		t.Fatal("se esperaba error para usuario inexistente, se obtuvo nil")
	}
}

// TestGetOtherUser verifica que GetOtherUser retorne el otro usuario del sistema.
func TestGetOtherUser(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	id1 := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash1")
	testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash2")

	repo := NewUserRepository()
	// Solicitar el "otro" usuario siendo anyel (id1)
	other, err := repo.GetOtherUser(uint(id1))

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if other.Username != "alexis" {
		t.Errorf("se esperaba 'alexis' como otro usuario, se obtuvo '%s'", other.Username)
	}
}

// TestGetOtherUser_SoloUnUsuario verifica que GetOtherUser retorne error si no hay otro usuario.
func TestGetOtherUser_SoloUnUsuario(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	id := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash1")

	repo := NewUserRepository()
	_, err := repo.GetOtherUser(uint(id))

	if err == nil {
		t.Fatal("se esperaba error cuando no hay otro usuario, se obtuvo nil")
	}
}
