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

func (d *IndexDataTable) InsertRow(newRow model.IndexRecord) error {
	insertSQL := `INSERT INTO indexdata(language, indexType, rating, order, date) VALUES (?, ?, ?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(newRow.Language, newRow.Type, newRow.Rating, newRow.Order, newRow.Month)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("Inserted data item: language=%s, index=%s, rating=%f, order=%d, date=%s\n", newRow.Language,
		newRow.Type, newRow.Rating, newRow.Order, newRow.Month)
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

func (d *IndexDataTable) GetRows(indexType string) (int, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM indexdata where indexType = ?`, indexType).Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
