package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Initialise() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to datab", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	err = db.Ping()

	if err != nil {
		log.Fatal("Error pinging database", err)
	}

	fmt.Println("Connected to the database")
}

func GetDB() *sql.DB {
	return db
}
