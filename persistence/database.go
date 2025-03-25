package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase(databaseFilePath string) error {
	fmt.Println("in create database")
	err := CreateDatabaseFileIfNotExists(databaseFilePath)
	if err != nil {
		return err
	}
	db, err := GetDatabase(databaseFilePath)
	if err != nil {
		return err
	}
	err = createLanguageTable(db)
	if err != nil {
		return err
	}
	err = createDataTable(db)
	if err != nil {
		return err
	}
	return nil
}

func createLanguageTable(db *sql.DB) error {
	createLanguageQuery := `CREATE TABLE IF NOT EXISTS language(
		language_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		parent_id INTEGER);`

	return createTable(db, createLanguageQuery)
}

func createDataTable(db *sql.DB) error {
	createDataQuery := `CREATE TABLE IF NOT EXISTS data(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		language_id INTEGER NOT NULL,
		count INTEGER NOT NULL,
		date TEXT NOT NULL);`
	return createTable(db, createDataQuery)
}

func createTable(db *sql.DB, query string) error {
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
