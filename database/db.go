package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	dbPath := "./data/app.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	os.MkdirAll("./data", 0755)

	var err error
	db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=ON")
	if err != nil {
		log.Fatal(err)
	}
}
