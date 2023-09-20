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

	dotenvPath := os.Getenv("DOTENV_PATH")

	err := godotenv.Load(dotenvPath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var connStr string

	if os.Getenv("ENVIRONMENT") == "development" {
		connStr = os.Getenv("TEST_DATABASE_URL")
	} else {
		connStr = os.Getenv("LIVE_DATABASE_URL")
	}

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to datab", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	err = db.Ping()

	if err != nil {
		log.Fatal("Error pinging database", err)
	}

	fmt.Println("Connected to the database")

}

func GetDB() *sql.DB {
	return db
}
