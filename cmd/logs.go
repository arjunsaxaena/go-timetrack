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
	logsToday    bool
	logsWeek     bool
	logsDays     int
	logsSeparate bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logged tasks (grouped by default)",
	Example: `  tt logs
  tt logs --separate
  tt logs --today
  tt logs --week
  tt logs --days 14`,
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

		if logsSeparate {
			logs, err := st.GetTaskLogs(since)
			if err != nil {
				return fmt.Errorf("could not get task logs: %w", err)
			}

			if len(logs) == 0 {
				printEmpty("No logs found.")
				return nil
			}

			printSection("Task Logs (Separate Sessions)")
			for i, entry := range logs {
				duration := time.Duration(entry.DurationSeconds) * time.Second
				fmt.Printf("# %s %s\n", uiID(entry.ID), entry.TaskName)
				printField("start", formatDateTime(entry.StartTime))
				printField("end", formatDateTime(entry.EndTime))
				printField("total", formatDuration(duration))
				if i < len(logs)-1 {
					fmt.Println()
				}
			}
			return nil
		}

		groups, err := st.GetTaskLogGroups(since)
		if err != nil {
			return fmt.Errorf("could not get grouped task logs: %w", err)
		}
		if len(groups) == 0 {
			printEmpty("No logs found.")
			return nil
		}

		printSection("Task Logs (Grouped)")
		for i, group := range groups {
			total := time.Duration(group.DurationSeconds) * time.Second
			fmt.Printf("%d) %s\n", i+1, group.TaskName)
			printField("total", formatDuration(total))
			printField("sessions", fmt.Sprintf("%d", group.SessionCount))
			if i < len(groups)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVar(&logsToday, "today", false, "show only today's logs")
	logsCmd.Flags().BoolVar(&logsWeek, "week", false, "show logs from last 7 days")
	logsCmd.Flags().IntVar(&logsDays, "days", 0, "show logs from the last N days")
	logsCmd.Flags().BoolVar(&logsSeparate, "separate", false, "show each log session separately")
}
