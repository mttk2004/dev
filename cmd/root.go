package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"dev/internal/doctor"
	"dev/internal/registry"
	"dev/internal/tui"
	"dev/internal/ui"
	"dev/internal/updater"
	"dev/internal/version"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dev",
	Short:   "A professional CLI to automate web dev environment setup",
	Version: version.Version,
	Long: `dev is a professional interactive CLI tool designed to automate tedious
tasks related to web development on Arch Linux, such as installing,
updating, and managing packages and services.

Just run 'dev' to launch the interactive dashboard!`,
	Run: runDashboard,
}

// runDashboard is the main loop that drives the interactive TUI.
func runDashboard(cmd *cobra.Command, args []string) {
	updater.StartAsyncCheck()

	choices := registry.IDs()

	for {
		ui.ClearScreen()
		updater.WaitForResult(3 * time.Second)

		report := doctor.RunChecks()
		report.Print()

		action, err := tui.RunDashboard()
		if err != nil {
			ui.Error("Dashboard error: %v", err)
			os.Exit(1)
		}

		actionTaken := handleAction(action, choices)

		if actionTaken && action != tui.ActionServices {
			fmt.Println("\nPress Enter to return to Dashboard...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}

// handleAction dispatches to the appropriate action handler and returns
// true if an action was performed that requires pausing for the user.
func handleAction(action tui.DashboardAction, choices []string) bool {
	switch action {
	case tui.ActionInstall:
		return tui.RunInstallAction(choices)
	case tui.ActionUpdate:
		return tui.RunUpdateAction(choices)
	case tui.ActionUninstall:
		return tui.RunUninstallAction(choices)
	case tui.ActionSearch:
		return tui.RunSearchAction()
	case tui.ActionServices:
		return tui.RunServiceAction()
	case tui.ActionConfig:
		return tui.RunConfigAction()
	case tui.ActionScaffold:
		return tui.RunScaffoldAction()
	case tui.ActionExit:
		ui.ClearScreen()
		fmt.Println("Goodbye! 👋")
		os.Exit(0)
	}
	return false
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
