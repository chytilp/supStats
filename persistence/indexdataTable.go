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
	insertSQL := `INSERT INTO indexdata(lang, indexType, order1, orderChange, orderPrevYear, rating, ratingChange, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := d.DB.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(newRow.Language, newRow.Type, newRow.Order, newRow.OrderChange, newRow.OrderPrevYear, newRow.Rating, newRow.RatingChange, newRow.Month)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("Inserted data item: %s\n", newRow.String())
	return nil
}

func (d *IndexDataTable) ExistsDate(date string, indexType string) (*bool, error) {
	var count int
	if err := d.DB.QueryRow(`SELECT COUNT(id) FROM indexdata WHERE date = ? and indexType = ?`, date, indexType).Scan(&count); err != nil {
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

func (d *IndexDataTable) RemoveMonth(fromDate string, toDate string, indexType string) error {
	fmt.Printf("Index table -> fromDt: %s, toDt: %s, indexType: %s\n", fromDate, toDate, indexType)
	removeSQL := `DELETE FROM indexdata WHERE date >= ? AND date < ? AND indexType = ?`
	statement, err := d.DB.Prepare(removeSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(fromDate, toDate, indexType)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}

/*
func (d *IndexDataTable) Select(fromDate string, toDate string, indexType string) ([]*Record2[float64], error) {
	var err error
	var rows *sql.Rows
	sql := "SELECT lang, rating, date FROM indexdata WHERE date >= $1 AND date < $2 AND AND indexType = $3"
	rows, err = d.DB.Query(sql, fromDate, toDate, indexType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records := make([]*Record2[float64], 8500)
	var index int = 0
	for rows.Next() {
		var record Record2[float64]
		err = rows.Scan(&record.lang, &record.value, &record.date)
		if err != nil {
			return nil, err
		}
		records[index] = &record
		index += 1
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	records = records[:index]
	return records, nil
}
*/
