package store

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
)

const (
	logIDLength   = 8
	logIDAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func IsValidLogID(id string) bool {
	if len(id) != logIDLength {
		return false
	}
	for _, ch := range id {
		if !strings.ContainsRune(logIDAlphabet, ch) {
			return false
		}
	}
	return true
}

func generateLogID() (string, error) {
	buf := make([]byte, logIDLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	out := make([]byte, logIDLength)
	for i := range buf {
		out[i] = logIDAlphabet[int(buf[i])%len(logIDAlphabet)]
	}
	return string(out), nil
}

func generateUniqueLogIDTx(tx *sql.Tx, table string) (string, error) {
	lookupQuery, err := logIDLookupQuery(table)
	if err != nil {
		return "", err
	}

	const maxAttempts = 32
	for i := 0; i < maxAttempts; i++ {
		id, err := generateLogID()
		if err != nil {
			return "", err
		}

		var found string
		err = tx.QueryRow(lookupQuery, id).Scan(&found)
		if err == nil {
			continue
		}
		if err == sql.ErrNoRows {
			return id, nil
		}
		return "", err
	}

	return "", fmt.Errorf("could not generate a unique log id")
}

func logIDLookupQuery(table string) (string, error) {
	switch table {
	case "task_log":
		return `SELECT id FROM task_log WHERE id = ?`, nil
	case "task_log_new":
		return `SELECT id FROM task_log_new WHERE id = ?`, nil
	default:
		return "", fmt.Errorf("unsupported log id table: %s", table)
	}
}
