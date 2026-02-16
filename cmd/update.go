/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

var (
	updateName  string
	updateStart string
	updateEnd   string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [log-id]",
	Short: "Update a task log entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := strings.TrimSpace(args[0])
		if !store.IsValidLogID(id) {
			return fmt.Errorf("log-id must be an 8-character alphanumeric value")
		}
		if updateName == "" && updateStart == "" && updateEnd == "" {
			return fmt.Errorf("provide at least one of --name, --start, or --end")
		}

		var namePtr *string
		if updateName != "" {
			trimmed := strings.TrimSpace(updateName)
			if trimmed == "" {
				return fmt.Errorf("--name cannot be empty")
			}
			namePtr = &trimmed
		}

		var startPtr *time.Time
		if updateStart != "" {
			startTime, err := parseDateTimeValue(updateStart, "--start")
			if err != nil {
				return err
			}
			startPtr = &startTime
		}

		var endPtr *time.Time
		if updateEnd != "" {
			endTime, err := parseDateTimeValue(updateEnd, "--end")
			if err != nil {
				return err
			}
			endPtr = &endTime
		}

		st, err := store.Open()
		if err != nil {
			return err
		}

		entry, err := st.UpdateTaskLog(id, namePtr, startPtr, endPtr)
		if err != nil {
			if errors.Is(err, store.ErrLogNotFound) {
				return fmt.Errorf("log with id %s not found", id)
			}
			if errors.Is(err, store.ErrInvalidTimeRange) {
				return fmt.Errorf("start time cannot be after end time")
			}
			return fmt.Errorf("could not update log: %w", err)
		}

		duration := time.Duration(entry.DurationSeconds) * time.Second
		printSuccess("Updated log #%s (%s)", entry.ID, entry.TaskName)
		printField("start", formatDateTime(entry.StartTime))
		printField("end", formatDateTime(entry.EndTime))
		printField("total", formatDuration(duration))
		return nil
	},
}

func parseDateTimeValue(input string, flagName string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04",
		"2006-01-02 3:04 PM",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, input, time.Local); err == nil {
			return t, nil
		}
	}

	clockOnlyLayouts := []string{
		"15:04",
		"3:04 PM",
	}
	now := time.Now()
	for _, layout := range clockOnlyLayouts {
		if t, err := time.ParseInLocation(layout, input, time.Local); err == nil {
			return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local), nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"invalid %s value. use RFC3339, '2006-01-02 15:04', '2006-01-02 3:04 PM', '15:04', or '3:04 PM'",
		flagName,
	)
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateName, "name", "", "new task name for the log")
	updateCmd.Flags().StringVar(&updateStart, "start", "", "new start time for the log")
	updateCmd.Flags().StringVar(&updateEnd, "end", "", "new end time for the log")
}
