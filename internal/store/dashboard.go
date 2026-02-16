package store

import "time"

func (s *Store) GetTaskDurationSummary(since *time.Time) ([]TaskDurationSummary, int, error) {
	query := `SELECT task_name, SUM(duration_seconds) as total_seconds
		FROM task_log`
	args := []any{}

	if since != nil {
		query += ` WHERE end_time >= ?`
		args = append(args, *since)
	}

	query += ` GROUP BY task_name ORDER BY total_seconds DESC, task_name ASC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var summaries []TaskDurationSummary
	totalSeconds := 0
	for rows.Next() {
		var row TaskDurationSummary
		if err := rows.Scan(&row.TaskName, &row.DurationSeconds); err != nil {
			return nil, 0, err
		}
		totalSeconds += row.DurationSeconds
		summaries = append(summaries, row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return summaries, totalSeconds, nil
}
