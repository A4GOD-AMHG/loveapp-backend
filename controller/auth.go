package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/model"
	"github.com/A4GOD-AMHG/LoveApp-Backend/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	t := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		util.JsonResp(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	user, err := FindUserByUsername(t.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(t.Password)) != nil {
		util.JsonResp(w, 401, map[string]string{"error": "invalid credentials"})
		return
	}
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(config.JwtSecret)
	if err != nil {
		util.JsonResp(w, 500, map[string]string{"error": "could not create token"})
		return
	}
	util.JsonResp(w, 200, map[string]string{"token": s})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*model.User)
	p := struct {
		Old string `json:"old"`
		New string `json:"new"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		util.JsonResp(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	var hashed string
	err := database.Db.QueryRow("SELECT password FROM users WHERE id = ?", user.ID).Scan(&hashed)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hashed), []byte(p.Old)) != nil {
		util.JsonResp(w, 403, map[string]string{"error": "old password incorrect"})
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(p.New), bcrypt.DefaultCost)
	if _, err := database.Db.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHash), user.ID); err != nil {
		util.JsonResp(w, 500, map[string]string{"error": "could not update"})
		return
	}
	util.JsonResp(w, 200, map[string]string{"ok": "password changed"})
}
