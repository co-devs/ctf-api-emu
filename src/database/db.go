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
	createTicksTable := `CREATE TABLE IF NOT EXISTS ticks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT
	);`
	createServicesTable := `CREATE TABLE IF NOT EXISTS services (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT UNIQUE NOT NULL
	);`
	createTeamsTable := `CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		key TEXT NOT NULL UNIQUE,
		is_admin BOOLEAN NOT NULL
	);`
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT NOT NULL,
		team_id INTEGER,
		FOREIGN KEY (team_id) REFERENCES teams(id)
	);`
	createEndpointsTable := `CREATE TABLE IF NOT EXISTS endpoints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		team_id INTEGER NOT NULL,
		service_id INT NOT NULL,
		hostname TEXT NOT NULL,
		FOREIGN KEY (service_id) REFERENCES services(id)
	);`
	createFlagsTable := `CREATE TABLE IF NOT EXISTS flags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		flag_identifier TEXT NOT NULL,
		flag TEXT UNIQUE NOT NULL,
		endpoint_id INTEGER NOT NULL,
		tick INTEGER NOT NULL,
		expiration TEXT NOT NULL,
		FOREIGN KEY (endpoint_id) REFERENCES endpoints(id),
		FOREIGN KEY (tick) REFERENCES ticks(id)
	);`
	createSubmittedFlagsTable := `CREATE TABLE IF NOT EXISTS submitted_flags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		flag_id TEXT NOT NULL,
		team_id INTEGER NOT NULL,
		-- service_id INTEGER NOT NULL,
		-- tick INTEGER NOT NULL,
		timestamp TEXT NOT NULL,
		FOREIGN KEY (flag_id) REFERENCES flags(id),
		FOREIGN KEY (team_id) REFERENCES teams(id)
	);`
	createStatusChecksTable := `CREATE TABLE IF NOT EXISTS status_checks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tick INTEGER NOT NULL,
		endpoint_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		FOREIGN KEY (endpoint_id) REFERENCES endpoints(id),
		FOREIGN KEY (tick) REFERENCES ticks(id)
	);`

	if _, err := DB.Exec(createTicksTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createServicesTable); err != nil {
		log.Fatal(err)
	}

		if _, err := DB.Exec(createTeamsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createEndpointsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createFlagsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(createSubmittedFlagsTable); err != nil {
		log.Fatal(err)
	}


	if _, err := DB.Exec(createStatusChecksTable); err != nil {
		log.Fatal(err)
	}
}
