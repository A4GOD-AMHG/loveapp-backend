// Tests unitarios para el paquete response — helpers de respuestas HTTP JSON.
package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
)

// TestJSON verifica que la función JSON establezca el Content-Type correcto
// y serialice los datos correctamente.
func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"clave": "valor"}

	JSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type esperado 'application/json', se obtuvo '%s'", ct)
	}

	var resultado map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resultado); err != nil {
		t.Fatalf("error al deserializar respuesta: %v", err)
	}
	if resultado["clave"] != "valor" {
		t.Errorf("valor esperado 'valor', se obtuvo '%s'", resultado["clave"])
	}
}

// TestBadRequest verifica que BadRequest envíe un HTTP 400 con el mensaje correcto.
func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "datos inválidos")

	if w.Code != http.StatusBadRequest {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusBadRequest, w.Code)
	}

	var resp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message != "datos inválidos" {
		t.Errorf("mensaje esperado 'datos inválidos', se obtuvo '%s'", resp.Message)
	}
	if resp.Code != http.StatusBadRequest {
		t.Errorf("código de error esperado %d, se obtuvo %d", http.StatusBadRequest, resp.Code)
	}
}

// TestUnauthorized verifica que Unauthorized envíe un HTTP 401 con el mensaje correcto.
func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	Unauthorized(w, "no autorizado")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusUnauthorized, w.Code)
	}

	var resp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message != "no autorizado" {
		t.Errorf("mensaje esperado 'no autorizado', se obtuvo '%s'", resp.Message)
	}
}

// TestForbidden verifica que Forbidden envíe un HTTP 403 con el mensaje correcto.
func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	Forbidden(w, "acceso denegado")

	if w.Code != http.StatusForbidden {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusForbidden, w.Code)
	}
}

// TestNotFound verifica que NotFound envíe un HTTP 404 con el mensaje correcto.
func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	NotFound(w, "no encontrado")

	if w.Code != http.StatusNotFound {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusNotFound, w.Code)
	}
}

// TestInternalServerError verifica que InternalServerError envíe un HTTP 500.
func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerError(w, "error interno")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusInternalServerError, w.Code)
	}
}

// TestSuccess verifica que Success envíe HTTP 200 con el mensaje y datos correctos.
func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, "operación exitosa", nil)

	if w.Code != http.StatusOK {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusOK, w.Code)
	}

	var resp models.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message != "operación exitosa" {
		t.Errorf("mensaje esperado 'operación exitosa', se obtuvo '%s'", resp.Message)
	}
}

// TestCreated verifica que Created envíe HTTP 201 con el mensaje correcto.
func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, "recurso creado", nil)

	if w.Code != http.StatusCreated {
		t.Errorf("código esperado %d, se obtuvo %d", http.StatusCreated, w.Code)
	}
}
