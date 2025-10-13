package commands

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/csotherden/strftime"

	"github.com/chytilp/supStats/persistence"
)

type RemoveMonthCommand struct {
	DB        *sql.DB
	Month     time.Time
	Version   *int
	IndexType *string
}

func NewSupdataRemoveMonth(db *sql.DB, month time.Time, version int) RemoveMonthCommand {
	return RemoveMonthCommand{
		DB:      db,
		Month:   month,
		Version: &version,
	}
}

func NewIndexRemoveMonth(db *sql.DB, month time.Time, indexType string) RemoveMonthCommand {
	return RemoveMonthCommand{
		DB:        db,
		Month:     month,
		IndexType: &indexType,
	}
}

func (r *RemoveMonthCommand) getFromTo() (string, string) {
	fmt.Printf("month: %v\n", r.Month)
	fromDt := strftime.Format("%Y-%m-%d %H:%M:%S", r.Month)
	end := r.Month.AddDate(0, 1, 0)
	toDt := strftime.Format("%Y-%m-%d %H:%M:%S", end)
	return fromDt, toDt
}

func (r *RemoveMonthCommand) Run() error {
	fromDt, toDt := r.getFromTo()
	fmt.Printf("%s - %s\n", fromDt, toDt)
	var err error
	if r.Version != nil {
		supdataTable := persistence.SupDataTable{DB: r.DB}
		err = supdataTable.RemoveMonth(fromDt, toDt, *r.Version)
	} else {
		indexTable := persistence.IndexDataTable{DB: r.DB}
		err = indexTable.RemoveMonth(fromDt, toDt, *r.IndexType)
	}
	if err != nil {
		fmt.Printf("Error in removeMonth db method: %v\n", err)
		return err
	}
	return nil
}
