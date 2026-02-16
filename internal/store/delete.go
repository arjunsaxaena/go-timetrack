package store

import "time"

func (s *Store) DeleteLogsSince(since time.Time) (int64, error) {
	result, err := s.db.Exec(`DELETE FROM task_log WHERE end_time >= ?`, since)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Store) DeleteLogByID(id string) error {
	result, err := s.db.Exec(`DELETE FROM task_log WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrLogNotFound
	}
	return nil
}

func (s *Store) DeleteActiveTask(task string) error {
	result, err := s.db.Exec(`DELETE FROM active_task WHERE task_name = ?`, task)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotActive
	}
	return nil
}

func (s *Store) DeleteAllData() (deletedLogs int64, deletedActive int64, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	logResult, err := tx.Exec(`DELETE FROM task_log`)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	activeResult, err := tx.Exec(`DELETE FROM active_task`)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	deletedLogs, err = logResult.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	deletedActive, err = activeResult.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	return deletedLogs, deletedActive, nil
}
