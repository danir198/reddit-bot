package datastore

import (
	"database/sql"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbName string) (*SQLiteStore, error) {
	dbPath := filepath.Join(".", dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS votes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        item_id TEXT NOT NULL,
        item_type TEXT NOT NULL,
        action TEXT NOT NULL,
        bot_id TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(item_id, bot_id)
    );`

	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) RecordVote(itemID, itemType, action, botID string) error {
	query := `
    INSERT OR REPLACE INTO votes (item_id, item_type, action, bot_id)
    VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, itemID, itemType, action, botID, time.Now())
	if err != nil {
		log.Printf("Error executing RecordVote query: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("RecordVote: Rows affected: %d", rowsAffected)

	return nil
}

func (s *SQLiteStore) HasVoted(itemID, botID string) (bool, string, error) {
	query := `
    SELECT action FROM votes
    WHERE item_id = ? AND bot_id = ? LIMIT 1`

	var action string
	err := s.db.QueryRow(query, itemID, botID).Scan(&action)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}

	return true, action, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
