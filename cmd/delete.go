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
	deleteAll    bool
	deleteToday  bool
	deleteDays   int
	deleteID     string
	deleteActive string
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete task logs or active tasks",
	Example: `  tt delete --today
  tt delete --days 7
  tt delete --id a1b2c3d4
  tt delete --active "deep work"
  tt delete --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args

		modeCount := 0
		if deleteAll {
			modeCount++
		}
		if deleteToday {
			modeCount++
		}
		daysFlagSet := cmd.Flags().Changed("days")
		if daysFlagSet {
			if deleteDays < 0 {
				return fmt.Errorf("--days must be >= 0")
			}
			modeCount++
		}
		if strings.TrimSpace(deleteID) != "" {
			modeCount++
		}
		if strings.TrimSpace(deleteActive) != "" {
			modeCount++
		}

		if modeCount == 0 {
			return fmt.Errorf("pick one delete mode: --all, --today, --days, --id, or --active")
		}
		if modeCount > 1 {
			return fmt.Errorf("use only one delete mode at a time")
		}

		st, err := store.Open()
		if err != nil {
			return err
		}

		now := time.Now()
		switch {
		case deleteAll:
			deletedLogs, deletedActive, err := st.DeleteAllData()
			if err != nil {
				return fmt.Errorf("could not delete all data: %w", err)
			}
			printSuccess("Deleted all tracked data")
			printField("logs", fmt.Sprintf("%d", deletedLogs))
			printField("active", fmt.Sprintf("%d", deletedActive))
			return nil

		case deleteToday:
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			deleted, err := st.DeleteLogsSince(startOfDay)
			if err != nil {
				return fmt.Errorf("could not delete today's logs: %w", err)
			}
			if deleted == 0 {
				printEmpty("No logs found for today.")
				return nil
			}
			printSuccess("Deleted today's logs")
			printField("count", fmt.Sprintf("%d", deleted))
			return nil

		case daysFlagSet:
			startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			since := startOfToday.AddDate(0, 0, -deleteDays)
			deleted, err := st.DeleteLogsSince(since)
			if err != nil {
				return fmt.Errorf("could not delete logs for last %d days: %w", deleteDays, err)
			}
			if deleted == 0 {
				printEmpty("No logs found in the last %d days.", deleteDays)
				return nil
			}
			printSuccess("Deleted logs from today - %d days", deleteDays)
			printField("count", fmt.Sprintf("%d", deleted))
			return nil

		case strings.TrimSpace(deleteID) != "":
			id := strings.TrimSpace(deleteID)
			if !store.IsValidLogID(id) {
				return fmt.Errorf("--id must be an 8-character alphanumeric value")
			}
			if err := st.DeleteLogByID(id); err != nil {
				if errors.Is(err, store.ErrLogNotFound) {
					return fmt.Errorf("log with id %s not found", id)
				}
				return fmt.Errorf("could not delete log %s: %w", id, err)
			}
			printSuccess("Deleted log %s", uiID(id))
			return nil

		case strings.TrimSpace(deleteActive) != "":
			task := strings.TrimSpace(deleteActive)
			if task == "" {
				return fmt.Errorf("--active cannot be empty")
			}
			if err := st.DeleteActiveTask(task); err != nil {
				if errors.Is(err, store.ErrTaskNotActive) {
					return fmt.Errorf("active task %q not found", task)
				}
				return fmt.Errorf("could not delete active task %q: %w", task, err)
			}
			printSuccess("Deleted active task %q", task)
			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteAll, "all", false, "delete all logs and active tasks")
	deleteCmd.Flags().BoolVar(&deleteToday, "today", false, "delete logs for today")
	deleteCmd.Flags().IntVar(&deleteDays, "days", 0, "delete logs from today - N days")
	deleteCmd.Flags().StringVar(&deleteID, "id", "", "delete a specific log by id")
	deleteCmd.Flags().StringVar(&deleteActive, "active", "", "delete an active task by name")

	_ = deleteCmd.RegisterFlagCompletionFunc("active", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		st, err := store.Open()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		suggestions, err := st.GetActiveTaskNameSuggestions(toComplete, 20)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})
}
