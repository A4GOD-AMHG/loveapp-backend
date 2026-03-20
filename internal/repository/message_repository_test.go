// Tests unitarios para messageRepository — operaciones CRUD de mensajes.
package repository

import (
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
)

// TestMessageCreate verifica que Create inserte un mensaje y retorne un ID válido.
func TestMessageCreate(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	msg := &models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Hola amor",
		Status:     "sent",
	}

	id, err := repo.Create(msg)
	if err != nil {
		t.Fatalf("error al crear mensaje: %v", err)
	}
	if id == 0 {
		t.Error("el ID del mensaje creado no debe ser 0")
	}
}

// TestMessageFindByID verifica que FindByID retorne el mensaje correcto con datos de sender y receiver.
func TestMessageFindByID(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	msg := &models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Te quiero mucho",
		Status:     "sent",
	}

	id, _ := repo.Create(msg)
	encontrado, err := repo.FindByID(id)

	if err != nil {
		t.Fatalf("error al buscar mensaje: %v", err)
	}
	if encontrado == nil {
		t.Fatal("se esperaba un mensaje, se obtuvo nil")
	}
	if encontrado.Content != "Te quiero mucho" {
		t.Errorf("contenido esperado 'Te quiero mucho', se obtuvo '%s'", encontrado.Content)
	}
	if encontrado.Sender.Username != "anyel" {
		t.Errorf("sender username esperado 'anyel', se obtuvo '%s'", encontrado.Sender.Username)
	}
	if encontrado.Receiver.Username != "alexis" {
		t.Errorf("receiver username esperado 'alexis', se obtuvo '%s'", encontrado.Receiver.Username)
	}
}

// TestMessageFindByID_NoExiste verifica que FindByID retorne nil para un ID inexistente.
func TestMessageFindByID_NoExiste(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewMessageRepository()
	encontrado, err := repo.FindByID(99999)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if encontrado != nil {
		t.Error("se esperaba nil para mensaje inexistente")
	}
}

// TestMessageUpdateContent verifica que UpdateContent modifique el contenido correctamente.
func TestMessageUpdateContent(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	msg := &models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Contenido original",
		Status:     "sent",
	}
	id, _ := repo.Create(msg)

	err := repo.UpdateContent(id, "Contenido editado")
	if err != nil {
		t.Fatalf("error al actualizar contenido: %v", err)
	}

	actualizado, _ := repo.FindByID(id)
	if actualizado.Content != "Contenido editado" {
		t.Errorf("contenido esperado 'Contenido editado', se obtuvo '%s'", actualizado.Content)
	}
}

// TestMessageUpdateStatus verifica que UpdateStatus cambie el estado del mensaje.
func TestMessageUpdateStatus(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	msg := &models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Mensaje de prueba",
		Status:     "sent",
	}
	id, _ := repo.Create(msg)

	err := repo.UpdateStatus(id, "read")
	if err != nil {
		t.Fatalf("error al actualizar estado: %v", err)
	}

	actualizado, _ := repo.FindByID(id)
	if actualizado.Status != "read" {
		t.Errorf("estado esperado 'read', se obtuvo '%s'", actualizado.Status)
	}
}

// TestMessageDelete verifica que Delete elimine el mensaje correctamente.
func TestMessageDelete(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	msg := &models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Mensaje a eliminar",
		Status:     "sent",
	}
	id, _ := repo.Create(msg)

	err := repo.Delete(id)
	if err != nil {
		t.Fatalf("error al eliminar mensaje: %v", err)
	}

	// Verificar que ya no existe
	eliminado, _ := repo.FindByID(id)
	if eliminado != nil {
		t.Error("el mensaje debería haberse eliminado, pero aún existe")
	}
}

// TestGetConversacion verifica que GetConversation retorne los mensajes entre dos usuarios.
func TestGetConversacion(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()

	// Insertar 3 mensajes entre los dos usuarios
	for i := 0; i < 3; i++ {
		repo.Create(&models.Message{
			SenderID:   uint(senderID),
			ReceiverID: uint(receiverID),
			Content:    "Mensaje de prueba",
			Status:     "sent",
		})
	}

	mensajes, err := repo.GetConversation(uint(senderID), uint(receiverID), 1, 10)
	if err != nil {
		t.Fatalf("error al obtener conversación: %v", err)
	}
	if len(mensajes) != 3 {
		t.Errorf("se esperaban 3 mensajes, se obtuvieron %d", len(mensajes))
	}
}

// TestGetConversacion_Paginacion verifica que la paginación funcione correctamente.
func TestGetConversacion_Paginacion(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()

	// Insertar 5 mensajes
	for i := 0; i < 5; i++ {
		repo.Create(&models.Message{
			SenderID:   uint(senderID),
			ReceiverID: uint(receiverID),
			Content:    "Mensaje",
			Status:     "sent",
		})
	}

	// Pedir solo 2 por página
	pagina1, err := repo.GetConversation(uint(senderID), uint(receiverID), 1, 2)
	if err != nil {
		t.Fatalf("error al obtener página 1: %v", err)
	}
	if len(pagina1) != 2 {
		t.Errorf("se esperaban 2 mensajes en página 1, se obtuvieron %d", len(pagina1))
	}

	pagina2, err := repo.GetConversation(uint(senderID), uint(receiverID), 2, 2)
	if err != nil {
		t.Fatalf("error al obtener página 2: %v", err)
	}
	if len(pagina2) != 2 {
		t.Errorf("se esperaban 2 mensajes en página 2, se obtuvieron %d", len(pagina2))
	}
}

// TestGetConversacion_SinMensajes verifica que GetConversation retorne una lista vacía
// cuando no hay mensajes entre los usuarios.
func TestGetConversacion_SinMensajes(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()
	mensajes, err := repo.GetConversation(uint(senderID), uint(receiverID), 1, 10)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(mensajes) != 0 {
		t.Errorf("se esperaba lista vacía, se obtuvieron %d mensajes", len(mensajes))
	}
}

func TestCountUnreadByReceiver(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	senderID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	receiverID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")

	repo := NewMessageRepository()

	_, _ = repo.Create(&models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Mensaje 1",
		Status:     "sent",
	})
	_, _ = repo.Create(&models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Mensaje 2",
		Status:     "delivered",
	})
	_, _ = repo.Create(&models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    "Mensaje 3",
		Status:     "read",
	})

	count, err := repo.CountUnreadByReceiver(uint(receiverID))
	if err != nil {
		t.Fatalf("error al contar mensajes no leídos: %v", err)
	}

	if count != 2 {
		t.Fatalf("se esperaban 2 mensajes no leídos, se obtuvieron %d", count)
	}
}
