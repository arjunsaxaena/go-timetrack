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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [task]",
	Short: "Stop tracking a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]

		st, err := store.Open()
		if err != nil {
			return err
		}

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
