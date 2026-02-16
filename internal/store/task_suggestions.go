package store

import "strings"

func (s *Store) GetTaskNameSuggestions(prefix string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20
	}

	prefix = strings.TrimSpace(prefix)

	rows, err := s.db.Query(
		`SELECT name
		 FROM (
			SELECT DISTINCT task_name AS name FROM task_log
			UNION
			SELECT DISTINCT task_name AS name FROM active_task
		 )
		 WHERE (? = '' OR LOWER(name) LIKE LOWER(?) || '%')
		 ORDER BY name ASC
		 LIMIT ?`,
		prefix,
		prefix,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suggestions := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suggestions, nil
}

func (s *Store) GetActiveTaskNameSuggestions(prefix string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20
	}

	prefix = strings.TrimSpace(prefix)

	rows, err := s.db.Query(
		`SELECT task_name
		 FROM active_task
		 WHERE (? = '' OR LOWER(task_name) LIKE LOWER(?) || '%')
		 ORDER BY task_name ASC
		 LIMIT ?`,
		prefix,
		prefix,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suggestions := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suggestions, nil
}
