package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS documents (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		file_hash TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		page_count INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL,
		chunk_count INTEGER NOT NULL DEFAULT 0,
		uploaded_at DATETIME NOT NULL,
		indexed_at DATETIME,
		error_message TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_documents_file_hash ON documents(file_hash);
	`
	_, err := db.Exec(query)
	return err
}
