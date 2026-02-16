/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"fmt"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.Open()
		if err != nil {
			return err
		}

		tasks, err := st.GetActiveTasks()
		if err != nil {
			return fmt.Errorf("could not get active tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No active tasks.")
			return nil
		}

		for _, task := range tasks {
			fmt.Printf("%s - started at %s\n", task.Name, task.StartTime.Local().Format("3:04 PM"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
