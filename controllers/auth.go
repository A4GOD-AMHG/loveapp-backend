package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		jsonResp(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	user, err := findUserByUsername(t.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(t.Password)) != nil {
		jsonResp(w, 401, map[string]string{"error": "invalid credentials"})
		return
	}
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(jwtSecret)
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": "could not create token"})
		return
	}
	jsonResp(w, 200, map[string]string{"token": s})
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	p := struct {
		Old string `json:"old"`
		New string `json:"new"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonResp(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	var hashed string
	err := db.QueryRow("SELECT password FROM users WHERE id = ?", user.ID).Scan(&hashed)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hashed), []byte(p.Old)) != nil {
		jsonResp(w, 403, map[string]string{"error": "old password incorrect"})
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(p.New), bcrypt.DefaultCost)
	if _, err := db.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHash), user.ID); err != nil {
		jsonResp(w, 500, map[string]string{"error": "could not update"})
		return
	}
	jsonResp(w, 200, map[string]string{"ok": "password changed"})
}
