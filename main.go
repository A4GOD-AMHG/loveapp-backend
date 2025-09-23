package main

import (
	"log"
	"net/http"

	"github.com/A4GOD-AMHG/LoveApp-Backend/database"
	"github.com/A4GOD-AMHG/LoveApp-Backend/config"
	"github.com/A4GOD-AMHG/LoveApp-Backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.InitConfig()
	database.InitDB()
	defer database.Db.Close()

	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}
	if err := database.Seed(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	addr := ":8080"
	log.Printf("listening %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
