/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"fmt"
	"time"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

var (
	logsToday bool
	logsWeek  bool
	logsDays  int
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logged tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		filterCount := 0
		if logsToday {
			filterCount++
		}
		if logsWeek {
			filterCount++
		}
		if logsDays != 0 {
			filterCount++
		}
		if filterCount > 1 {
			return fmt.Errorf("use only one of --today, --week, or --days")
		}
		if logsDays < 0 {
			return fmt.Errorf("--days must be >= 0")
		}

		st, err := store.Open()
		if err != nil {
			return err
		}

		var since *time.Time
		now := time.Now()
		switch {
		case logsToday:
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			since = &startOfDay
		case logsWeek:
			weekAgo := now.AddDate(0, 0, -7)
			since = &weekAgo
		case logsDays > 0:
			daysAgo := now.Add(-time.Duration(logsDays) * 24 * time.Hour)
			since = &daysAgo
		}

		logs, err := st.GetTaskLogs(since)
		if err != nil {
			return fmt.Errorf("could not get task logs: %w", err)
		}

		if len(logs) == 0 {
			fmt.Println("No logs found.")
			return nil
		}

		for _, entry := range logs {
			duration := time.Duration(entry.DurationSeconds) * time.Second
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60

			fmt.Printf(
				"#%d %s | %s -> %s | %dh %dm\n",
				entry.ID,
				entry.TaskName,
				entry.StartTime.Local().Format("Jan 2 3:04 PM"),
				entry.EndTime.Local().Format("Jan 2 3:04 PM"),
				hours,
				minutes,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVar(&logsToday, "today", false, "show only today's logs")
	logsCmd.Flags().BoolVar(&logsWeek, "week", false, "show logs from last 7 days")
	logsCmd.Flags().IntVar(&logsDays, "days", 0, "show logs from the last N days")
}
