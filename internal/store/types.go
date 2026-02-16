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

type TaskLogEntry struct {
	ID              int64
	TaskName        string
	StartTime       time.Time
	EndTime         time.Time
	DurationSeconds int
}
