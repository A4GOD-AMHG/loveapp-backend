package main

import (
	"net/http"
	"github.com/A4GOD-AMHG/LoveApp-Backend/controllers"
	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/login", controllers.loginHandler).Methods("POST")

	auth := r.NewRoute().Subrouter()
	auth.Use(authMiddleware)
	auth.HandleFunc("/change-password", changePasswordHandler).Methods("POST")
	auth.HandleFunc("/todos", createTodoHandler).Methods("POST")
	auth.HandleFunc("/todos", listTodosHandler).Methods("GET")
	auth.HandleFunc("/todos/{id}", deleteTodoHandler).Methods("DELETE")
	auth.HandleFunc("/todos/{id}/complete", completeTodoHandler).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}).Methods("GET")
}
