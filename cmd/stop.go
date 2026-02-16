/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"time"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [task]",
	Short: "Stop tracking a task (or all active tasks)",
	Example: `  tt stop "deep work"
  tt stop`,
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		st, err := store.Open()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		suggestions, err := st.GetActiveTaskNameSuggestions(toComplete, 20)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.Open()
		if err != nil {
			return err
		}

		if len(args) == 1 {
			task := args[0]
			duration, err := st.StopTask(task)
			if err != nil {
				if errors.Is(err, store.ErrTaskNotActive) {
					return fmt.Errorf("task %q is not active", task)
				}
				return fmt.Errorf("could not stop task: %w", err)
			}

			printSuccess("Stopped task %q", task)
			printField("spent", formatDuration(duration))
			return nil
		}

		activeTasks, err := st.GetActiveTasks()
		if err != nil {
			return fmt.Errorf("could not get active tasks: %w", err)
		}
		if len(activeTasks) == 0 {
			printEmpty("No active tasks.")
			return nil
		}

		var total time.Duration
		for _, task := range activeTasks {
			duration, err := st.StopTask(task.Name)
			if err != nil {
				return fmt.Errorf("could not stop task %q: %w", task.Name, err)
			}
			total += duration
		}

		printSuccess("Stopped all active tasks")
		printField("count", fmt.Sprintf("%d", len(activeTasks)))
		printField("total", formatDuration(total))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
