package persistence

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chytilp/dbForSupTest/model"
)

type LanguageTable struct {
	DB *sql.DB
}

func (l *LanguageTable) InsertLanguage(newRow model.LanguageRow) (int, error) {
	insertSQL := `INSERT INTO language(name, path, parent_id) VALUES (?, ?, ?)`
	statement, err := l.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	result, err := statement.Exec(newRow.Name, newRow.Path, newRow.ParentId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var parent int = 0
	if newRow.ParentId != nil {
		parent = *newRow.ParentId
	}
	fmt.Printf("Inserted language: name=%s, path=%s, parent=%d\n", newRow.Name, newRow.Path, parent)
	newId, err := result.LastInsertId()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return int(newId), nil
}

func (l *LanguageTable) GetLanguageId(name string) (*int, error) {
	var languageId int
	if err := l.DB.QueryRow(`SELECT language_id FROM language WHERE name = ?`, name).Scan(&languageId); err != nil {
		if err == sql.ErrNoRows {
			languageId = 0
			return &languageId, nil
		}
		return nil, err
	}
	return &languageId, nil
}

func (l *LanguageTable) GetRows() (int, error) {
	var count int
	if err := l.DB.QueryRow(`SELECT COUNT(language_id) FROM language`).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
