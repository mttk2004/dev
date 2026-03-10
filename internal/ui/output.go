package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Styles used across all output helpers.
var (
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	styleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	styleInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	styleSubtle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// Success prints a success message to stdout.
//
//	✅ Successfully installed Node.js!
func Success(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleSuccess.Render("✅ " + msg))
}

// Error prints an error message to stdout.
//
//	❌ Failed to install Node.js: exit status 1
func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleError.Render("❌ " + msg))
}

// Warning prints a warning message to stdout.
//
//	⚠️  Package "xyz" not found in registry
func Warning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleWarning.Render("⚠️  " + msg))
}

// Info prints an informational message to stdout.
//
//	📦 Installing Node.js...
func Info(icon string, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleInfo.Render(icon + " " + msg))
}

// Subtle prints a dimmed/secondary message to stdout.
//
//	Search finished.
func Subtle(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleSubtle.Render(msg))
}

// Fatalf prints an error message and returns a formatted error.
// This does NOT call os.Exit — the caller decides what to do.
func Fatalf(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(styleError.Render("❌ " + msg))
	return fmt.Errorf("%s", msg)
}

// ActionHeader prints a prominent header before an action starts.
//
//	📦 Installing Docker...
func ActionHeader(icon string, action string, target string) {
	fmt.Printf("\n%s %s %s...\n", icon, action, target)
}

// ActionResult prints the result of an action on a specific target.
func ActionResult(target string, err error, successVerb string) {
	if err != nil {
		Error("Failed to %s %s: %v", successVerb, target, err)
	} else {
		Success("Successfully %s %s!", successVerb, target)
	}
}

// ClearScreen sends ANSI escape codes to clear the terminal.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
