package store

import "time"

func (s *Store) GetTaskLogs(since *time.Time) ([]TaskLogEntry, error) {
	query := `SELECT id, task_name, start_time, end_time, duration_seconds FROM task_log`
	args := []any{}

	if since != nil {
		query += ` WHERE end_time >= ?`
		args = append(args, *since)
	}

	query += ` ORDER BY end_time DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []TaskLogEntry
	for rows.Next() {
		var entry TaskLogEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.TaskName,
			&entry.StartTime,
			&entry.EndTime,
			&entry.DurationSeconds,
		); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
