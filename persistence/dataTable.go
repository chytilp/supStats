package persistence

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chytilp/supStats/model"
)

type DataTable struct {
	DB *sql.DB
}

func (d *DataTable) InsertDataItem(newRow model.DataRow) error {
	insertSQL := `INSERT INTO data(language_id, count, date) VALUES (?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(newRow.LanguageId, newRow.Count, newRow.Date)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("Inserted data item: language_id=%d, count=%d, date=%s\n", newRow.LanguageId, newRow.Count, newRow.Date)
	return nil
}

func (d *DataTable) ExistsDate(date string) (*bool, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM data WHERE date = ?`, date).Scan(&count); err != nil {
		return nil, err
	}
	exists := count > 0
	return &exists, nil
}

func (d *DataTable) GetRows() (int, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM data`).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
