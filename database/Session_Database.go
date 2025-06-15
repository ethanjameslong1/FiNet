package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func addSessionID(id uuid.UUID, username string, maxage time.Time) (bool, error) {
	db, err := NewService(DriverName, SessionDataSource)
	r, err := db.db.Exec(SQL_INSERT_SESSION, id.String(), time.Now()+maxage)
	if err != nil {
		return false, fmt.Errorf("Error inserting User: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Error finding Rows Affected: %w", err)
	}
	if rows != 1 {
		return false, fmt.Errorf("Expected rows impacted to be 1, rows impacted %d", rows)
	}
	return true, nil
}
