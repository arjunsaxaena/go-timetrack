package store

func (s *Store) StartTask(task string) error {
	result, err := s.db.Exec(
		`INSERT INTO active_task (task_name, start_time)
		 VALUES (?, datetime('now'))
		 ON CONFLICT(task_name) DO NOTHING`,
		task,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskAlreadyActive
	}

	return nil
}
