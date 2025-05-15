package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type CreateTableFunc func(db *sql.DB) error

func createDatabase(databaseFilePath string, createTableFunc *CreateTableFunc) error {
	fmt.Println("in create database")
	err := CreateDatabaseFileIfNotExists(databaseFilePath)
	if err != nil {
		return err
	}
	var db *sql.DB
	db, err = GetDatabase(databaseFilePath)
	if err != nil {
		return err
	}
	if createTableFunc != nil {
		f := *createTableFunc
		err = f(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSupDatabase(databaseFilePath string) error {
	var fce CreateTableFunc = createSupDataTable
	return createDatabase(databaseFilePath, &fce)
}

func CreateIndexesDatabase(databaseFilePath string) error {
	var fce CreateTableFunc = createIndexDataTable
	return createDatabase(databaseFilePath, &fce)
}

func createSupDataTable(db *sql.DB) error {
	// type: 1 - category, 2 - technology
	// version: 24 - old format, 25 - new format
	createDataQuery := `CREATE TABLE IF NOT EXISTS supdata(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		language TEXT NOT NULL,
		type INTEGER NOT NULL, 
		count INTEGER NOT NULL,
		version INTEGER NOT NULL,
		date DATETIME NOT NULL);`
	return createTable(db, createDataQuery)
}

func createIndexDataTable(db *sql.DB) error {
	// indexType: 1 - tiobe, 2 - pypl
	createDataQuery := `CREATE TABLE IF NOT EXISTS indexdata(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		language TEXT NOT NULL,
		indexType INTEGER NOT NULL,
		rating REAL NOT NULL,
		order INTEGER NOT NULL,
		date DATETIME NOT NULL);`
	return createTable(db, createDataQuery)
}
