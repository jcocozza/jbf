package sqlite

import (
	"database/sql"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "test.db"

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	Schema(db)
	return db, nil
}

//go:embed schema.sql
var schema string

//go:embed truncate.sql
var truncate string

func Schema(db *sql.DB) {
	_, err := db.Exec(truncate)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(schema)
	if err != nil {
		panic(err)
	}
}
