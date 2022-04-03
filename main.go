package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
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
	if _, err := tx.Exec(`INSERT INTO tbl1(title, data) VALUES(?, ?)`, "Hello", text); err != nil {
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

	const latestVersion = 1 // must be incremented whenever schema changes
	var ver int
	if err := tx.QueryRow(`SELECT user_version FROM pragma_user_version LIMIT 1`).Scan(&ver); err != nil {
		return fmt.Errorf("reading current user_version value: %w", err)
	}
	if ver == latestVersion {
		return nil
	}
	if ver > latestVersion {
		return fmt.Errorf("database schema version %d is newer than supported %d", ver, latestVersion)
	}
	// we may either have an existing database with ver=0 (the initial schema),
	// or a newly created empty database, also with ver=0; need to distinguish
	// these two cases
	const tableExistsQuery = `SELECT 1 FROM pragma_table_list WHERE name='tbl1' AND type='table' AND schema='main'`
	var discard int
	switch err := tx.QueryRow(tableExistsQuery).Scan(&discard); err {
	case nil: // table exists, need to migrate schema
		if err := migrateSchema(tx, ver); err != nil {
			return fmt.Errorf("migrating schema from (ver %d->%d): %w", ver, latestVersion, err)
		}
	case sql.ErrNoRows: // fresh database
	default:
		return fmt.Errorf("querying table list: %w", err)
	}

	for _, s := range [...]string{ // don't forget to increment latestVersion const on any changes
		`CREATE TABLE IF NOT EXISTS tbl1(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			data TEXT NOT NULL
		)`,
	} {
		if _, err := tx.Exec(s); err != nil {
			return fmt.Errorf("statement %q: %w", s, err)
		}
	}
	if _, err := tx.Exec(`PRAGMA user_version=` + strconv.Itoa(latestVersion)); err != nil {
		return fmt.Errorf("setting new user_version: %w", err)
	}
	return tx.Commit()
}

func migrateSchema(tx *sql.Tx, ver int) error {
	if ver < 0 {
		return fmt.Errorf("invalid version %d, negative values are not supported", ver)
	}
	// https://sqlite.org/lang_altertable.html#making_other_kinds_of_table_schema_changes
	switch ver {
	case 0:
		for _, s := range [...]string{
			`CREATE TABLE new_tbl1(
				id INTEGER PRIMARY KEY,
				title TEXT NOT NULL, -- legacy migrated records have this set to ":FIXME:"
				data TEXT NOT NULL
			)`,
			`INSERT INTO new_tbl1 SELECT id, ':FIXME:', data FROM tbl1`,
			`DROP TABLE tbl1`,
			`ALTER TABLE new_tbl1 RENAME TO tbl1`,
		} {
			if _, err := tx.Exec(s); err != nil {
				return fmt.Errorf("statement %q: %w", s, err)
			}
		}
	default:
		return fmt.Errorf("schema migration from version %d is not supported", ver)
	}
	return nil
}
