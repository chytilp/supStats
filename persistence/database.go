package persistence

import (
	"database/sql"
	"fmt"
	"os"

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
	dropTable := "DROP TABLE indexdata;"
	err := createTable(db, dropTable)
	if err != nil {
		fmt.Printf("Drop table failed: %v\n", err)
		os.Exit(1)
	}
	// indexType: [tiobe, pypl]
	createDataQuery := `CREATE TABLE IF NOT EXISTS indexdata(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		lang TEXT NOT NULL,
		indexType TEXT NOT NULL,
		order1 INTEGER NOT NULL,
		orderChange TEXT,
		orderPrevYear INTEGER,
		rating DECIMAL(5,2) NOT NULL,
        ratingChange DECIMAL(5,2),
		date DATETIME NOT NULL);`
	return createTable(db, createDataQuery)
}
