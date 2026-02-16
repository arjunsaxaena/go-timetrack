package store

import (
	"database/sql"
	"errors"
	"time"
)

func (s *Store) StopTask(task string) (time.Duration, error) {
	var startTime time.Time

	err := s.db.QueryRow(
		`SELECT start_time FROM active_task WHERE task_name = ?`,
		task,
	).Scan(&startTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTaskNotActive
		}
		return 0, err
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	logID, err := generateUniqueLogIDTx(tx, "task_log")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(
		`INSERT INTO task_log (id, task_name, start_time, end_time, duration_seconds)
		 VALUES (?, ?, ?, ?, ?)`,
		logID,
		task,
		startTime,
		endTime,
		int(duration.Seconds()),
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(
		`DELETE FROM active_task WHERE task_name = ?`,
		task,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return duration, nil
}
