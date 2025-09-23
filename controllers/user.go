package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var allowedCreators = map[string]bool{"anyel": true, "alexis": true}

func jsonResp(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func findUserByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func findUserByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow("SELECT id, username FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return u, nil
}
