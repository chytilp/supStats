package persistence

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chytilp/supStats/model"
)

type SupDataTable struct {
	DB *sql.DB
}

func (d *SupDataTable) InsertRow(newRow model.SupdataRow) error {
	insertSQL := `INSERT INTO supdata(language, type, count, version, date) VALUES (?, ?, ?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(newRow.Language, newRow.Type, newRow.Count, newRow.Version, newRow.Date)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("Inserted data item: language=%s, type=%d, count=%d, date=%s\n", newRow.Language, newRow.Type, newRow.Count, newRow.Date)
	return nil
}

func (d *SupDataTable) ExistsDate(date string, version int) (*bool, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM supdata WHERE date = ? AND version = ?`, date, version).Scan(&count); err != nil {
		return nil, err
	}
	exists := count > 0
	return &exists, nil
}

func (d *SupDataTable) GetRows(version int) (int, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM supdata WHERE version = ?`, version).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}

func (d *SupDataTable) RemoveMonth(fromDate string, toDate string, version int) error {
	fmt.Printf("Supdata table -> fromDt: %s, toDt: %s, vesrion: %d\n", fromDate, toDate, version)
	removeSQL := `DELETE FROM supdata WHERE date >= ? AND date < ? AND version = ?`
	statement, err := d.DB.Prepare(removeSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(fromDate, toDate, version)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}
