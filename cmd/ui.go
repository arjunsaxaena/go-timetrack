package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var uiColorEnabled = detectColorSupport()

func detectColorSupport() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func uiColor(code string, text string) string {
	if !uiColorEnabled {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}

func uiMuted(text string) string {
	return uiColor("90", text)
}

func uiAccent(text string) string {
	return uiColor("36", text)
}

func uiGood(text string) string {
	return uiColor("32", text)
}

func uiWarn(text string) string {
	return uiColor("33", text)
}

func uiID(text string) string {
	return uiColor("35", text)
}

func printSection(title string) {
	line := strings.Repeat("-", len(title)+4)
	fmt.Println(uiMuted(line))
	fmt.Println(uiAccent("  " + title))
	fmt.Println(uiMuted(line))
}

func printSuccess(msg string, args ...any) {
	fmt.Println(uiGood("[OK] " + fmt.Sprintf(msg, args...)))
}

func printInfo(msg string, args ...any) {
	fmt.Println(uiAccent("[i] " + fmt.Sprintf(msg, args...)))
}

func printEmpty(msg string, args ...any) {
	fmt.Println(uiWarn("[ ] " + fmt.Sprintf(msg, args...)))
}

func printField(label string, value string) {
	fmt.Printf("  %-7s %s\n", label+":", value)
}

func formatDateTime(t time.Time) string {
	return t.Local().Format("Jan 2, 3:04 PM")
}

func formatClock(t time.Time) string {
	return t.Local().Format("3:04 PM")
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	totalSeconds := int(d.Round(time.Second).Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
