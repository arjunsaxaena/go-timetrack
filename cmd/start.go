/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [task]",
	Short: "Start tracking a task",
	Example: `  tt start "deep work"
  tt start "meeting"`,
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		st, err := store.Open()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		suggestions, err := st.GetTaskNameSuggestions(toComplete, 20)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]

		st, err := store.Open()
		if err != nil {
			return err
		}

		if err := st.StartTask(task); err != nil {
			if errors.Is(err, store.ErrTaskAlreadyActive) {
				return fmt.Errorf("task %q is already active", task)
			}
			return fmt.Errorf("could not start task: %w", err)
		}

		printSuccess("Started task %q", task)
		printInfo("Use %q to see active timers.", "tt status")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
