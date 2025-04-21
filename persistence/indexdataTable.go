package persistence

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chytilp/supStats/model"
)

type IndexDataTable struct {
	DB *sql.DB
}

func (d *IndexDataTable) InsertRow(newRow model.IndexdataRow) error {
	insertSQL := `INSERT INTO indexdata(language, indexType, rating, order, date) VALUES (?, ?, ?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(newRow.Language, newRow.IndexType, newRow.Rating, newRow.Order, newRow.Date)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("Inserted data item: language=%s, index=%d, rating=%f, order=%d, date=%s\n", newRow.Language, newRow.IndexType,
		newRow.Rating, newRow.Order, newRow.Date)
	return nil
}

func (d *IndexDataTable) ExistsDate(date string) (*bool, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM indexdata WHERE date = ?`, date).Scan(&count); err != nil {
		return nil, err
	}
	exists := count > 0
	return &exists, nil
}

func (d *IndexDataTable) GetRows() (int, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM indexdata`).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
