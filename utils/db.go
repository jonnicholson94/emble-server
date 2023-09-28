package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Initialise() {

	var connStr string

	if os.Getenv("ENVIRONMENT") == "development" {
		connStr = os.Getenv("TEST_DATABASE_URL")
	} else {
		connStr = os.Getenv("LIVE_DATABASE_URL")
	}

	db, _ = sql.Open("postgres", connStr)

	_ = db.Ping()

	fmt.Println("Connected to the database")

}

func GetDB() *sql.DB {
	return db
}
