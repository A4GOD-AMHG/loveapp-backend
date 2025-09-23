package routes

import (
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/middleware"
	"github.com/A4GOD-AMHG/LoveApp-Backend/controller"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")

	auth := r.NewRoute().Subrouter()
	auth.Use(middleware.AuthMiddleware)
	auth.HandleFunc("/change-password", controller.ChangePasswordHandler).Methods("POST")
	auth.HandleFunc("/todos", controller.CreateTodoHandler).Methods("POST")
	auth.HandleFunc("/todos", controller.ListTodosHandler).Methods("GET")
	auth.HandleFunc("/todos/{id}", controller.DeleteTodoHandler).Methods("DELETE")
	auth.HandleFunc("/todos/{id}/complete", controller.CompleteTodoHandler).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}).Methods("GET")
}
