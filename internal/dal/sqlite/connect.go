package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "jbf.db"

func CreateDB() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	err = Schema(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func Connect() (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database does not exist. please init first.")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectAndClean() (*sql.DB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	err = Schema(db)
	if err != nil {
		return nil, err
	}
	err = Truncate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//go:embed schema.sql
var schema string

func Schema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}

//go:embed truncate.sql
var truncate string

func Truncate(db *sql.DB) error {
	_, err := db.Exec(truncate)
	return err
}
