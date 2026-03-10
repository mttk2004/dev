package cmd

import (
	"bufio"
	"fmt"
	"os"

	"dev/internal/doctor"
	"dev/internal/pkgmanager"
	"dev/internal/registry"
	"dev/internal/scaffold"
	"dev/internal/system"
	"dev/internal/tui"

	"github.com/spf13/cobra"
)

// getPackageByID retrieves a package from the registry by its ID
func getPackageByID(id string) (*registry.Package, error) {
	for _, p := range registry.Packages {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("package %s not found", id)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "A professional CLI to automate web dev environment setup",
	Long: `dev is a professional interactive CLI tool designed to automate tedious
tasks related to web development on Arch Linux, such as installing,
updating, and managing packages and services.

Just run 'dev' to launch the interactive dashboard!`,
	Run: func(cmd *cobra.Command, args []string) {
		// Dynamically build choices from the registry
		var choices []string
		for _, p := range registry.Packages {
			choices = append(choices, p.ID)
		}

		for {
			// Clear screen before each iteration to keep it fresh and update the doctor report
			fmt.Print("\033[H\033[2J")

			// 1. Run Doctor Check and print the table report
			report := doctor.RunChecks()
			report.Print()

			// 2. Open Main Dashboard Menu
			action, err := tui.RunDashboard()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			var actionTaken bool

			switch action {
			case tui.ActionInstall:
				selected, err := tui.RunSelector("install", choices)
				if err != nil || len(selected) == 0 {
					continue
				}
				for _, pkgID := range selected {
					p, err := getPackageByID(pkgID)
					if err != nil {
						continue
					}
					fmt.Printf("\n📦 Installing %s...\n", p.DisplayName)
					if err := p.Install(); err != nil {
						fmt.Printf("❌ Failed to install %s: %v\n", p.DisplayName, err)
					} else {
						fmt.Printf("✅ Successfully installed %s!\n", p.DisplayName)
					}
				}
				actionTaken = true

			case tui.ActionUpdate:
				selected, err := tui.RunSelector("update", choices)
				if err != nil || len(selected) == 0 {
					continue
				}
				for _, pkgID := range selected {
					p, err := getPackageByID(pkgID)
					if err != nil {
						continue
					}
					fmt.Printf("\n🔄 Updating %s...\n", p.DisplayName)
					if err := p.Update(); err != nil {
						fmt.Printf("❌ Failed to update %s: %v\n", p.DisplayName, err)
					} else {
						fmt.Printf("✅ Successfully updated %s!\n", p.DisplayName)
					}
				}
				actionTaken = true

			case tui.ActionUninstall:
				selected, err := tui.RunSelector("uninstall", choices)
				if err != nil || len(selected) == 0 {
					continue
				}
				for _, pkgID := range selected {
					p, err := getPackageByID(pkgID)
					if err != nil {
						continue
					}
					fmt.Printf("\n🧹 Uninstalling %s...\n", p.DisplayName)
					if err := p.Remove(); err != nil {
						fmt.Printf("❌ Failed to uninstall %s: %v\n", p.DisplayName, err)
					} else {
						fmt.Printf("✅ Successfully uninstalled %s!\n", p.DisplayName)
					}
				}
				actionTaken = true

			case tui.ActionSearch:
				query, err := tui.RunSearchPrompt()
				if err != nil || query == "" {
					continue
				}
				pm := system.GetPackageManager()
				fmt.Printf("\n🔍 Searching for '%s' using %s...\n", query, pm)
				err = pkgmanager.RunInteractive(pm, "-Ss", query)
				if err != nil {
					fmt.Printf("\nSearch finished.\n")
				}
				actionTaken = true

			case tui.ActionServices:
				err := tui.RunServiceManager()
				if err != nil {
					fmt.Printf("Error managing services: %v\n", err)
					actionTaken = true
				}
				// Service manager has its own interactive loop and pauses, no need to pause here again

			case tui.ActionScaffold:
				pType, pDir, pName, err := tui.RunScaffoldPrompt()
				if err != nil {
					if err.Error() != "user aborted" {
						fmt.Printf("Error scaffolding project: %v\n", err)
						actionTaken = true
					}
					continue
				}
				if err := scaffold.CreateProject(pType, pDir, pName); err != nil {
					fmt.Printf("❌ Failed to create project: %v\n", err)
				}
				actionTaken = true

			case tui.ActionExit:
				fmt.Print("\033[H\033[2J") // Clean exit
				fmt.Println("Goodbye! 👋")
				os.Exit(0)
			}

			// Pause so the user can read the logs of the action they just performed
			if actionTaken && action != tui.ActionServices {
				fmt.Println("\nPress Enter to return to Dashboard...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
