package store

import (
	"database/sql"
	"errors"
	"time"
)

func (s *Store) UpdateTaskLog(
	id int64,
	taskName *string,
	startTime *time.Time,
	endTime *time.Time,
) (TaskLogEntry, error) {
	var entry TaskLogEntry

	err := s.db.QueryRow(
		`SELECT id, task_name, start_time, end_time, duration_seconds
		 FROM task_log
		 WHERE id = ?`,
		id,
	).Scan(
		&entry.ID,
		&entry.TaskName,
		&entry.StartTime,
		&entry.EndTime,
		&entry.DurationSeconds,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskLogEntry{}, ErrLogNotFound
		}
		return TaskLogEntry{}, err
	}

	updatedName := entry.TaskName
	updatedStart := entry.StartTime
	updatedEnd := entry.EndTime

	if taskName != nil {
		updatedName = *taskName
	}
	if startTime != nil {
		updatedStart = *startTime
	}
	if endTime != nil {
		updatedEnd = *endTime
	}

	if updatedEnd.Before(updatedStart) {
		return TaskLogEntry{}, ErrInvalidTimeRange
	}

	durationSeconds := int(updatedEnd.Sub(updatedStart).Seconds())

	_, err = s.db.Exec(
		`UPDATE task_log
		 SET task_name = ?, start_time = ?, end_time = ?, duration_seconds = ?
		 WHERE id = ?`,
		updatedName,
		updatedStart,
		updatedEnd,
		durationSeconds,
		id,
	)
	if err != nil {
		return TaskLogEntry{}, err
	}

	entry.TaskName = updatedName
	entry.StartTime = updatedStart
	entry.EndTime = updatedEnd
	entry.DurationSeconds = durationSeconds
	return entry, nil
}
