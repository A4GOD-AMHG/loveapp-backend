package controller

import (
	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/model"
)

var allowedCreators = map[string]bool{"anyel": true, "alexis": true}

func FindUserByUsername(username string) (*model.User, error) {
	u := &model.User{}
	err := database.Db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func FindUserByID(id int64) (*model.User, error) {
	u := &model.User{}
	err := database.Db.QueryRow("SELECT id, username FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return u, nil
}
