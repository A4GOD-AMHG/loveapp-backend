package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/auth"
	"github.com/A4GOD-AMHG/LoveApp-Backend/pkg/response"
	"github.com/gorilla/mux"
)

// currentUser obtiene el usuario autenticado o responde 401.
func currentUser(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		response.Unauthorized(w, "Tu sesión ha expirado, por favor inicia sesión de nuevo.")
	}
	return user, ok
}

// pathID parsea el parámetro {id} de la URL o responde 400 con invalidMsg.
func pathID(w http.ResponseWriter, r *http.Request, invalidMsg string) (int64, bool) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		response.BadRequest(w, invalidMsg)
		return 0, false
	}
	return id, true
}

// decodeJSON decodifica el cuerpo en dst o responde 400 con invalidMsg.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any, invalidMsg string) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.BadRequest(w, invalidMsg)
		return false
	}
	return true
}
