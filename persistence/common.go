package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func GetDatabase(databaseFilePath string) (*sql.DB, error) {
	if database == nil {
		db, err := sql.Open("sqlite3", databaseFilePath)
		if err != nil {
			fmt.Printf("error during connecting database: %v, file: %s\n", err, databaseFilePath)
			return nil, err
		}
		database = db
	}
	return database, nil
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateDatabaseFileIfNotExists(databaseFilePath string) error {
	path, err := filepath.Abs(databaseFilePath)
	if err != nil {
		fmt.Printf("error during filepath.abs of file: %s, err: %v\n", databaseFilePath, err)
		return err
	}
	fmt.Printf("database file path: %s\n", path)
	if PathExists(path) {
		return nil
	}
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		folderPath := fmt.Sprintf("%s/", strings.Join(parts[0:len(parts)-1], "/"))
		if !PathExists(folderPath) {
			err := os.MkdirAll(folderPath, 0755)
			if err != nil {
				fmt.Printf("error during creating folder: %s, err: %v\n", folderPath, err)
				return err
			}
		}
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("error during creating file: %s, err: %v\n", path, err)
		return err
	}
	defer file.Close()
	return nil
}
