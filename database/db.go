package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTables()
	return DB
}

func createTables() {
	createAlbumsTable := `CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		price REAL NOT NULL
	);`

	createAPIKeysTable := `CREATE TABLE IF NOT EXISTS api_keys (
		key TEXT PRIMARY KEY UNIQUE,
		can_access_secret BOOLEAN NOT NULL DEFAULT 0,
		can_add_album BOOLEAN NOT NULL DEFAULT 0,
		can_view_album BOOLEAN NOT NULL DEFAULT 0
	);`

	if _, err := DB.Exec(createAlbumsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createAPIKeysTable); err != nil {
		log.Fatal(err)
	}
}
