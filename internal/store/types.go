package store

import (
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

type ActiveTask struct {
	Name      string
	StartTime time.Time
}
