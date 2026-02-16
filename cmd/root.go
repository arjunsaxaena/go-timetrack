/*
Copyright © 2026 ARJUN SAXENA arjunsaxena04@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "tt",
	Short:        "Track work time from your terminal",
	Long:         "TimeTrack is a lightweight CLI to start, stop, inspect, and edit task time logs.",
	SilenceUsage: true,
	Example: `  tt start "project setup"
  tt status
  tt stop "project setup"
  tt logs --today
  tt dash --month
  tt update a1b2c3d4 --name "setup review" --end "6:30 PM"
  tt delete --today`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tt.yaml)")

}
