package cmd

import (
	"fmt"
	"strings"
	"time"
	"tt/internal/store"

	"github.com/spf13/cobra"
)

var (
	dashboardToday bool
	dashboardWeek  bool
	dashboardMonth bool
	dashboardAll   bool
	dashboardSince string
)

var dashboardCmd = &cobra.Command{
	Use:   "dash",
	Short: "Show time spent per task",
	Example: `  tt dash
  tt dash --week
  tt dash --month
  tt dash --all
  tt dash --week --since 2026-02-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args

		periodCount := 0
		if dashboardToday {
			periodCount++
		}
		if dashboardWeek {
			periodCount++
		}
		if dashboardMonth {
			periodCount++
		}
		if dashboardAll {
			periodCount++
		}
		if periodCount > 1 {
			return fmt.Errorf("use only one of --today, --week, --month, or --all")
		}

		now := time.Now()
		since, periodLabel, err := dashboardSinceStart(now, cmd.Flags().Changed("since"))
		if err != nil {
			return err
		}

		st, err := store.Open()
		if err != nil {
			return err
		}

		var sincePtr *time.Time
		if !since.IsZero() {
			sincePtr = &since
		}

		rows, totalSeconds, err := st.GetTaskDurationSummary(sincePtr)
		if err != nil {
			return fmt.Errorf("could not build dashboard: %w", err)
		}
		if len(rows) == 0 {
			printEmpty("No logs found for %s.", periodLabel)
			return nil
		}

		printSection("Dashboard")
		printField("period", periodLabel)
		if sincePtr != nil {
			printField("since", since.Local().Format("2006-01-02"))
		}
		printField("total", formatDuration(time.Duration(totalSeconds)*time.Second))
		shareBaseSeconds, shareBaseLabel := dashboardShareBaseSeconds(now, since, cmd.Flags().Changed("since"), periodLabel, totalSeconds)
		printField("base", shareBaseLabel)
		fmt.Println()

		// Histogram is intentionally disabled for now.
		// chartLines, legendLines := dashboardVerticalHistogram(rows, shareBaseSeconds)
		// for _, line := range chartLines {
		// 	fmt.Println(line)
		// }
		// for _, line := range legendLines {
		// 	fmt.Println(line)
		// }
		// fmt.Println()

		for i, row := range rows {
			pct := (float64(row.DurationSeconds) / float64(shareBaseSeconds)) * 100
			fmt.Printf("%d) %s\n", i+1, row.TaskName)
			printField("time", formatDuration(time.Duration(row.DurationSeconds)*time.Second))
			printField("share", fmt.Sprintf("%.1f%%", pct))
			if i < len(rows)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func dashboardSinceStart(now time.Time, sinceFlagSet bool) (time.Time, string, error) {
	var periodStart time.Time
	periodLabel := "today"

	switch {
	case dashboardAll:
		periodLabel = "all time"
	case dashboardWeek:
		periodLabel = "this week"
		periodStart = startOfCurrentWeek(now)
	case dashboardMonth:
		periodLabel = "this month"
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	default:
		periodLabel = "today"
		periodStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	if sinceFlagSet {
		value := strings.TrimSpace(dashboardSince)
		if value == "" {
			return time.Time{}, "", fmt.Errorf("--since cannot be empty")
		}
		parsed, err := time.ParseInLocation("2006-01-02", value, now.Location())
		if err != nil {
			return time.Time{}, "", fmt.Errorf("invalid --since value. use YYYY-MM-DD")
		}

		if periodStart.IsZero() || parsed.After(periodStart) {
			periodStart = parsed
		}
	}

	return periodStart, periodLabel, nil
}

func startOfCurrentWeek(now time.Time) time.Time {
	weekday := int(now.Weekday())
	daysSinceMonday := (weekday + 6) % 7
	base := now.AddDate(0, 0, -daysSinceMonday)
	return time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())
}

func dashboardShareBaseSeconds(now time.Time, since time.Time, sinceFlagSet bool, periodLabel string, totalSeconds int) (int, string) {
	switch periodLabel {
	case "today":
		return 24 * 60 * 60, "24h"
	case "this week":
		return 7 * 24 * 60 * 60, "168h"
	case "this month":
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
		daysInMonth := int(firstOfNextMonth.Sub(firstOfMonth).Hours() / 24)
		baseSeconds := daysInMonth * 24 * 60 * 60
		return baseSeconds, fmt.Sprintf("%dh (%d days)", daysInMonth*24, daysInMonth)
	case "all time":
		if sinceFlagSet && !since.IsZero() {
			seconds := int(time.Since(since).Seconds())
			if seconds < 1 {
				seconds = 1
			}
			return seconds, fmt.Sprintf("since %s", since.Local().Format("2006-01-02"))
		}
		if totalSeconds > 0 {
			return totalSeconds, "tracked total"
		}
		return 1, "tracked total"
	default:
		return 24 * 60 * 60, "24h"
	}
}

func dashboardVerticalHistogram(rows []store.TaskDurationSummary, shareBaseSeconds int) ([]string, []string) {
	const (
		maxBars = 12
		height  = 10
	)

	type bar struct {
		name    string
		seconds int
		pct     float64
	}

	bars := make([]bar, 0, len(rows))
	for _, row := range rows {
		pct := (float64(row.DurationSeconds) / float64(shareBaseSeconds)) * 100
		bars = append(bars, bar{name: row.TaskName, seconds: row.DurationSeconds, pct: pct})
	}

	visible := bars
	if len(bars) > maxBars {
		visible = append([]bar{}, bars[:maxBars-1]...)
		others := bar{name: "others"}
		for _, b := range bars[maxBars-1:] {
			others.seconds += b.seconds
			others.pct += b.pct
		}
		visible = append(visible, others)
	}

	chart := []string{"Histogram (% of base):"}
	for level := height; level >= 1; level-- {
		threshold := (float64(level) / float64(height)) * 100
		line := fmt.Sprintf("%3.0f%% |", threshold)
		for _, b := range visible {
			if b.pct >= threshold {
				line += " " + uiAccent("█") + " "
			} else {
				line += "   "
			}
		}
		chart = append(chart, line)
	}

	axis := "    +" + strings.Repeat("---", len(visible))
	labels := "     "
	for i := range visible {
		labels += fmt.Sprintf("%2d ", i+1)
	}
	chart = append(chart, axis, labels)

	legend := []string{"Legend:"}
	for i, b := range visible {
		legend = append(
			legend,
			fmt.Sprintf(
				"  %2d) %s - %.1f%% (%s)",
				i+1,
				b.name,
				b.pct,
				formatDuration(time.Duration(b.seconds)*time.Second),
			),
		)
	}

	return chart, legend
}

func init() {
	rootCmd.AddCommand(dashboardCmd)

	dashboardCmd.Flags().BoolVar(&dashboardToday, "today", false, "show dashboard for today (default)")
	dashboardCmd.Flags().BoolVar(&dashboardWeek, "week", false, "show dashboard for the current week")
	dashboardCmd.Flags().BoolVar(&dashboardMonth, "month", false, "show dashboard for the current month")
	dashboardCmd.Flags().BoolVar(&dashboardAll, "all", false, "show dashboard for all-time logs")
	dashboardCmd.Flags().StringVar(&dashboardSince, "since", "", "show data since date (YYYY-MM-DD)")
}
