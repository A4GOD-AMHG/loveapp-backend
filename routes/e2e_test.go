package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/websocket"
	"github.com/A4GOD-AMHG/LoveApp-Backend/routes"
	gws "github.com/gorilla/websocket"
)

type client struct {
	t     *testing.T
	base  string
	token string
}

func (c *client) do(method, path string, body any) (int, map[string]any) {
	c.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, c.base+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func login(t *testing.T, base, user, pass string) *client {
	t.Helper()
	c := &client{t: t, base: base}
	status, body := c.do("POST", "/auth/login", map[string]string{"username": user, "password": pass})
	if status != http.StatusOK {
		t.Fatalf("login %s: status %d %v", user, status, body)
	}
	c.token = body["token"].(string)
	return c
}

func (c *client) dialWS() *gws.Conn {
	c.t.Helper()
	url := "ws" + strings.TrimPrefix(c.base, "http") + "/ws?token=" + c.token
	conn, _, err := gws.DefaultDialer.Dial(url, nil)
	if err != nil {
		c.t.Fatalf("ws dial: %v", err)
	}
	c.t.Cleanup(func() { conn.Close() })
	return conn
}

func nextEvent(t *testing.T, conn *gws.Conn, wantType string) map[string]any {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		var ev map[string]any
		if err := conn.ReadJSON(&ev); err != nil {
			t.Fatalf("esperando evento %s: %v", wantType, err)
		}
		if ev["type"] == wantType {
			return ev["payload"].(map[string]any)
		}
	}
}

// TestE2E_AnyelYAlexis recorre el flujo real de la app con ambos usuarios:
// login, WebSocket, mensajes (envío, entrega, lectura, edición), permisos,
// tareas compartidas, tokens push, CORS y cambio de contraseña.
func TestE2E_AnyelYAlexis(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
	}
	if err := database.Seed(); err != nil {
		t.Fatal(err)
	}

	hub := websocket.NewHub()
	go hub.Run()
	srv := httptest.NewServer(routes.SetupRoutes(hub))
	defer srv.Close()

	if status, _ := (&client{t: t, base: srv.URL}).do("POST", "/auth/login",
		map[string]string{"username": "anyel", "password": "mala"}); status != http.StatusUnauthorized {
		t.Fatalf("contraseña incorrecta debe dar 401, dio %d", status)
	}
	if status, _ := (&client{t: t, base: srv.URL}).do("GET", "/todos", nil); status != http.StatusUnauthorized {
		t.Fatalf("ruta protegida sin token debe dar 401, dio %d", status)
	}

	anyel := login(t, srv.URL, "anyel", "anyel")
	alexis := login(t, srv.URL, "alexis", "alexis")
	anyelWS, alexisWS := anyel.dialWS(), alexis.dialWS()
	time.Sleep(50 * time.Millisecond) // el hub registra las conexiones de forma asíncrona

	// Mensajes en tiempo real
	status, msg := anyel.do("POST", "/messages", map[string]string{"content": "  hola amor  "})
	if status != http.StatusCreated || msg["content"] != "hola amor" {
		t.Fatalf("enviar mensaje: %d %v", status, msg)
	}
	id := int(msg["id"].(float64))
	if got := nextEvent(t, alexisWS, "message_sent"); got["content"] != "hola amor" {
		t.Fatalf("Alexis no recibió el mensaje: %v", got)
	}
	nextEvent(t, anyelWS, "message_sent")

	if status, _ := anyel.do("POST", "/messages", map[string]string{"content": "   "}); status != http.StatusBadRequest {
		t.Fatalf("mensaje vacío debe dar 400, dio %d", status)
	}

	_, unread := alexis.do("GET", "/messages/unread-count", nil)
	if unread["unread_count"] != float64(1) {
		t.Fatalf("Alexis debe tener 1 no leído: %v", unread)
	}

	if status, _ := anyel.do("PATCH", "/messages/"+itoa(id)+"/read", nil); status != http.StatusForbidden {
		t.Fatalf("el remitente no puede marcar como leído: %d", status)
	}
	alexis.do("PATCH", "/messages/"+itoa(id)+"/delivered", nil)
	if got := nextEvent(t, anyelWS, "message_delivered"); got["status"] != "delivered" {
		t.Fatalf("evento delivered: %v", got)
	}
	alexis.do("PATCH", "/messages/"+itoa(id)+"/read", nil)
	if got := nextEvent(t, anyelWS, "message_read"); got["status"] != "read" {
		t.Fatalf("evento read: %v", got)
	}
	_, unread = alexis.do("GET", "/messages/unread-count", nil)
	if unread["unread_count"] != float64(0) {
		t.Fatalf("Alexis debe tener 0 no leídos: %v", unread)
	}

	if status, _ := alexis.do("PUT", "/messages/"+itoa(id), map[string]string{"content": "x"}); status != http.StatusForbidden {
		t.Fatalf("Alexis no puede editar mensajes de Anyel: %d", status)
	}
	if status, _ := anyel.do("PUT", "/messages/"+itoa(id), map[string]string{"content": "hola mi amor"}); status != http.StatusOK {
		t.Fatalf("editar: %d", status)
	}
	if got := nextEvent(t, alexisWS, "message_updated"); got["content"] != "hola mi amor" {
		t.Fatalf("message_updated debe incluir el contenido nuevo: %v", got)
	}

	req, _ := http.NewRequest("GET", srv.URL+"/messages/conversation?page=1&per_page=10", nil)
	req.Header.Set("Authorization", "Bearer "+alexis.token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var conversation []map[string]any
	_ = json.NewDecoder(res.Body).Decode(&conversation)
	res.Body.Close()
	if len(conversation) != 1 || conversation[0]["status"] != "read" {
		t.Fatalf("conversación: %v", conversation)
	}

	// Tareas compartidas: se completa solo cuando ambos la marcan
	status, created := anyel.do("POST", "/todos", map[string]string{"title": "Reservar cena"})
	if status != http.StatusCreated {
		t.Fatalf("crear tarea: %d %v", status, created)
	}
	todoID := itoa(int(created["todo"].(map[string]any)["id"].(float64)))
	if status, _ := alexis.do("DELETE", "/todos/"+todoID, nil); status != http.StatusForbidden {
		t.Fatalf("solo el creador borra: %d", status)
	}
	anyel.do("PATCH", "/todos/"+todoID, map[string]bool{"completed": true})
	_, done := alexis.do("PATCH", "/todos/"+todoID, map[string]bool{"completed": true})
	if done["todo"].(map[string]any)["is_completed"] != true {
		t.Fatalf("la tarea debe estar completa por ambos: %v", done)
	}
	if status, _ := alexis.do("PATCH", "/todos/"+todoID, map[string]bool{"completed": false}); status != http.StatusBadRequest {
		t.Fatalf("no se puede desmarcar una tarea completada por ambos: %d", status)
	}
	if status, _ := anyel.do("GET", "/todos/999", nil); status != http.StatusMethodNotAllowed {
		t.Fatalf("GET /todos/{id} no existe: %d", status)
	}
	if status, _ := anyel.do("DELETE", "/todos/999", nil); status != http.StatusNotFound {
		t.Fatalf("borrar tarea inexistente: %d", status)
	}

	// Notificaciones push: registro y baja del token del dispositivo
	if status, _ := alexis.do("POST", "/devices/push-token", map[string]string{
		"platform": "android", "push_token": "tok-1", "device_name": "Pixel",
	}); status != http.StatusOK {
		t.Fatalf("registrar push token: %d", status)
	}
	if status, _ := alexis.do("DELETE", "/devices/push-token", map[string]string{"push_token": "tok-1"}); status != http.StatusNoContent {
		t.Fatalf("borrar push token: %d", status)
	}

	// CORS: el preflight de PATCH debe responder aunque la ruta no declare OPTIONS
	pre, _ := http.NewRequest("OPTIONS", srv.URL+"/messages/1/read", nil)
	pre.Header.Set("Origin", "http://localhost")
	pre.Header.Set("Access-Control-Request-Method", "PATCH")
	preRes, err := http.DefaultClient.Do(pre)
	if err != nil {
		t.Fatal(err)
	}
	preRes.Body.Close()
	if preRes.StatusCode != http.StatusOK || !strings.Contains(preRes.Header.Get("Access-Control-Allow-Methods"), "PATCH") {
		t.Fatalf("preflight CORS: %d %q", preRes.StatusCode, preRes.Header.Get("Access-Control-Allow-Methods"))
	}

	// Cambio de contraseña
	if status, _ := anyel.do("POST", "/auth/change-password", map[string]string{"new_password": "nueva123"}); status != http.StatusOK {
		t.Fatalf("cambiar contraseña: %d", status)
	}
	login(t, srv.URL, "anyel", "nueva123")
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}
