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
	ID              string
	TaskName        string
	StartTime       time.Time
	EndTime         time.Time
	DurationSeconds int
}

type TaskDurationSummary struct {
	TaskName        string
	DurationSeconds int
}

type TaskLogGroup struct {
	TaskName        string
	DurationSeconds int
	SessionCount    int
}
