package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB() {
	dbPath := "./data/loveapp.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	os.MkdirAll("./data", 0755)

	var err error
	Db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=ON")
	if err != nil {
		log.Fatal(err)
	}
}
