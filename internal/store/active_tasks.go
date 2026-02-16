package store

func (s *Store) GetActiveTasks() ([]ActiveTask, error) {
	var tasks []ActiveTask

	rows, err := s.db.Query(
		`SELECT task_name, start_time FROM active_task ORDER BY start_time ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task ActiveTask
		if err := rows.Scan(&task.Name, &task.StartTime); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
