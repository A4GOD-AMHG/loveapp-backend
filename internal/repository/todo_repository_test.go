// Tests unitarios para TodoRepository — operaciones CRUD de tareas.
package repository

import (
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
)

// TestTodoCreate verifica que Create persista una tarea y retorne su ID generado.
func TestTodoCreate(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")

	repo := NewTodoRepository()
	todo := &models.Todo{
		Title:       "Comprar comida",
		Description: "Frutas y verduras",
		CreatorID:   creatorID,
	}

	creado, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("error al crear todo: %v", err)
	}
	if creado.ID == 0 {
		t.Error("el ID del todo creado no debe ser 0")
	}
	if creado.Title != "Comprar comida" {
		t.Errorf("título esperado 'Comprar comida', se obtuvo '%s'", creado.Title)
	}
	if creado.CreatorUsername != "anyel" {
		t.Errorf("creator_username esperado 'anyel', se obtuvo '%s'", creado.CreatorUsername)
	}
}

// TestTodoFindByID verifica que FindByID retorne la tarea correcta por su ID.
func TestTodoFindByID(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Lavar ropa", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	encontrado, err := repo.FindByID(creado.ID)
	if err != nil {
		t.Fatalf("error al buscar todo: %v", err)
	}
	if encontrado.ID != creado.ID {
		t.Errorf("ID esperado %d, se obtuvo %d", creado.ID, encontrado.ID)
	}
}

// TestTodoFindByID_NoExiste verifica que FindByID retorne error para un ID inexistente.
func TestTodoFindByID_NoExiste(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewTodoRepository()
	_, err := repo.FindByID(99999)

	if err == nil {
		t.Fatal("se esperaba error para todo inexistente, se obtuvo nil")
	}
}

// TestTodoUpdate verifica que Update modifique el título y la descripción correctamente.
func TestTodoUpdate(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Título original", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	creado.Title = "Título actualizado"
	creado.Description = "Descripción nueva"
	actualizado, err := repo.Update(creado)

	if err != nil {
		t.Fatalf("error al actualizar todo: %v", err)
	}
	if actualizado.Title != "Título actualizado" {
		t.Errorf("título esperado 'Título actualizado', se obtuvo '%s'", actualizado.Title)
	}
}

// TestTodoUpdateCompletion_Anyel verifica que UpdateCompletion marque correctamente para anyel.
func TestTodoUpdateCompletion_Anyel(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Tarea", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	actualizado, err := repo.UpdateCompletion(creado.ID, "anyel", true)
	if err != nil {
		t.Fatalf("error al marcar completado: %v", err)
	}
	if !actualizado.CompletedAnyel {
		t.Error("CompletedAnyel debería ser true después de marcar")
	}
	if actualizado.IsCompleted {
		t.Error("IsCompleted debe ser false si solo anyel completó la tarea")
	}
}

// TestTodoUpdateCompletion_AmbosCumplen verifica que IsCompleted sea true
// cuando ambos usuarios completan la tarea.
func TestTodoUpdateCompletion_AmbosCumplen(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Tarea conjunta", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	repo.UpdateCompletion(creado.ID, "anyel", true)
	actualizado, err := repo.UpdateCompletion(creado.ID, "alexis", true)

	if err != nil {
		t.Fatalf("error al marcar completado: %v", err)
	}
	if !actualizado.IsCompleted {
		t.Error("IsCompleted debería ser true cuando ambos completaron")
	}
}

// TestTodoUpdateCompletion_UsuarioInvalido verifica que UpdateCompletion rechace
// un username que no sea "anyel" ni "alexis".
func TestTodoUpdateCompletion_UsuarioInvalido(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewTodoRepository()
	_, err := repo.UpdateCompletion(1, "hackerXD", true)

	if err == nil {
		t.Fatal("se esperaba error para username inválido, se obtuvo nil")
	}
}

// TestTodoDelete verifica que Delete elimine la tarea correctamente.
func TestTodoDelete(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Eliminar esto", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	err := repo.Delete(creado.ID)
	if err != nil {
		t.Fatalf("error al eliminar todo: %v", err)
	}

	// Verificar que ya no existe
	_, err = repo.FindByID(creado.ID)
	if err == nil {
		t.Fatal("el todo debería haber sido eliminado, pero aún existe")
	}
}

// TestTodoDelete_NoExiste verifica que Delete retorne error para un ID inexistente.
func TestTodoDelete_NoExiste(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	repo := NewTodoRepository()
	err := repo.Delete(99999)

	if err == nil {
		t.Fatal("se esperaba error al eliminar todo inexistente, se obtuvo nil")
	}
}

// TestGetCreatorID verifica que GetCreatorID retorne el ID del creador correcto.
func TestGetCreatorID(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	todo := &models.Todo{Title: "Mi tarea", CreatorID: creatorID}
	creado, _ := repo.Create(todo)

	obtenerID, err := repo.GetCreatorID(creado.ID)
	if err != nil {
		t.Fatalf("error al obtener creator_id: %v", err)
	}
	if obtenerID != creatorID {
		t.Errorf("creator_id esperado %d, se obtuvo %d", creatorID, obtenerID)
	}
}

// TestGetTodos_SinFiltros verifica que GetTodos retorne todas las tareas sin filtros.
func TestGetTodos_SinFiltros(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	repo.Create(&models.Todo{Title: "Tarea 1", CreatorID: creatorID})
	repo.Create(&models.Todo{Title: "Tarea 2", CreatorID: creatorID})

	todos, total, err := repo.GetTodos(models.TodoStatusAll, nil, "anyel", "desc", "", 10, 0)
	if err != nil {
		t.Fatalf("error al obtener todos: %v", err)
	}
	if total != 2 {
		t.Errorf("total esperado 2, se obtuvo %d", total)
	}
	if len(todos) != 2 {
		t.Errorf("cantidad de todos esperada 2, se obtuvo %d", len(todos))
	}
}

// TestGetTodos_FiltroCompletados verifica el filtro de tareas completadas por ambos.
func TestGetTodos_FiltroCompletados(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	creatorID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewTodoRepository()

	tarea1, _ := repo.Create(&models.Todo{Title: "Completa", CreatorID: creatorID})
	repo.Create(&models.Todo{Title: "Incompleta", CreatorID: creatorID})

	// Completar tarea1 por ambos usuarios
	repo.UpdateCompletion(tarea1.ID, "anyel", true)
	repo.UpdateCompletion(tarea1.ID, "alexis", true)

	todos, total, err := repo.GetTodos(models.TodoStatusCompleted, nil, "anyel", "desc", "", 10, 0)
	if err != nil {
		t.Fatalf("error al filtrar completados: %v", err)
	}
	if total != 1 {
		t.Errorf("total esperado 1 tarea completada, se obtuvo %d", total)
	}
	if len(todos) != 1 || todos[0].Title != "Completa" {
		t.Error("se esperaba solo la tarea completada en los resultados")
	}
}
