package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	if !allowedCreators[user.Username] {
		jsonResp(w, 403, map[string]string{"error": "user not allowed to create todos"})
		return
	}
	in := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResp(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	res, err := db.Exec("INSERT INTO todos (title, description, creator_id) VALUES (?, ?, ?)", in.Title, in.Description, user.ID)
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": "db error"})
		return
	}
	id, _ := res.LastInsertId()
	jsonResp(w, 201, map[string]any{"id": id})
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var creatorID int64
	err := db.QueryRow("SELECT creator_id FROM todos WHERE id = ?", id).Scan(&creatorID)
	if err == sql.ErrNoRows {
		jsonResp(w, 404, map[string]string{"error": "not found"})
		return
	}
	if creatorID != user.ID {
		jsonResp(w, 403, map[string]string{"error": "only creator can delete"})
		return
	}
	db.Exec("DELETE FROM todos WHERE id = ?", id)
	jsonResp(w, 200, map[string]string{"ok": "deleted"})
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("state")
	rows, err := db.Query("SELECT t.id, t.title, t.description, t.creator_id, u.username, t.completed_anyel, t.completed_alexis, t.created_at FROM todos t JOIN users u ON u.id = t.creator_id")
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": "db error"})
		return
	}
	defer rows.Close()
	res := []Todo{}
	for rows.Next() {
		var t Todo
		var ca, cb int
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatorID, &t.CreatorUsername, &ca, &cb, &t.CreatedAt); err != nil {
			continue
		}
		t.CompletedAnyel = ca != 0
		t.CompletedAlexis = cb != 0
		isDone := t.CompletedAnyel && t.CompletedAlexis
		if q == "done" && !isDone {
			continue
		}
		if q == "pending" && isDone {
			continue
		}
		res = append(res, t)
	}
	jsonResp(w, 200, res)
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var ca, cb int
	err := db.QueryRow("SELECT completed_anyel, completed_alexis FROM todos WHERE id = ?", id).Scan(&ca, &cb)
	if err == sql.ErrNoRows {
		jsonResp(w, 404, map[string]string{"error": "not found"})
		return
	}
	if user.Username == "anyel" {
		ca = 1
	} else if user.Username == "alexis" {
		cb = 1
	}
	db.Exec("UPDATE todos SET completed_anyel = ?, completed_alexis = ? WHERE id = ?", ca, cb, id)
	jsonResp(w, 200, map[string]string{"ok": "marked"})
}
