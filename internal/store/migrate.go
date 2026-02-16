package store

import "database/sql"

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS active_task (
			task_name TEXT PRIMARY KEY,
			start_time DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS task_log (
			id TEXT PRIMARY KEY,
			task_name TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			duration_seconds INTEGER NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
