package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	if err := initSchema(db); err != nil {
		return fmt.Errorf("schema init: %w", err)
	}
	if err := populateDB(db); err != nil {
		return err
	}
	return db.Close()
}

func populateDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	text := fmt.Sprintf("Text inserted at %s", time.Now().Format(time.RFC822Z))
	if _, err := tx.Exec(`INSERT INTO tbl1(data) VALUES(?)`, text); err != nil {
		return err
	}
	return tx.Commit()
}

func initSchema(db *sql.DB) error {
	for _, s := range [...]string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=normal`,
	} {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("statement %q: %w", s, err)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, s := range [...]string{
		`CREATE TABLE IF NOT EXISTS tbl1(
			id INTEGER PRIMARY KEY,
			data TEXT NOT NULL
		)`,
	} {
		if _, err := tx.Exec(s); err != nil {
			return fmt.Errorf("statement %q: %w", s, err)
		}
	}
	return tx.Commit()
}
