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
	Args:  cobra.ExactArgs(1),
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

		fmt.Printf("Started task: %s\n", task)
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
