package database

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dataDir string, migrationsDir string) (*Database, error) {
	db, err := sql.Open("sqlite3", path.Join(dataDir, "gift-cards.db"))
	if err != nil {
		return nil, err
	}
	dbInstance := &Database{db: db}
	if err = dbInstance.migrate(migrationsDir); err != nil {
		return nil, err
	}
	if err = dbInstance.prepareCards(); err != nil {
		return nil, err
	}
	return dbInstance, nil
}
