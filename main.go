package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initConfig()
	initDB()
	defer db.Close()

	if err := migrate(); err != nil {
		log.Fatal(err)
	}
	if err := seed(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	registerRoutes(r)

	addr := ":8080"
	log.Printf("listening %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
