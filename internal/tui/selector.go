package tui

import (
	"fmt"
	"strings"

	"dev/internal/system"

	"github.com/charmbracelet/huh"
)

// RunSelector opens an interactive TUI form to select items from a list
// based on the provided action (install, update, uninstall) and returns the chosen item names.
func RunSelector(action string, choices []string) ([]string, error) {
	var selected []string

	var availableOptions []huh.Option[string]
	var unavailablePkgs []string

	// Determine which packages are applicable based on the action
	for _, choice := range choices {
		cmdToCheck := choice

		// Map package names to their actual executable names for checking
		switch choice {
		case "node":
			// We use fnm to manage node
			cmdToCheck = "fnm"
		case "jdk":
			cmdToCheck = "java"
		case "postgresql":
			cmdToCheck = "psql"
		case "redis":
			cmdToCheck = "redis-cli"
		case "maven":
			cmdToCheck = "mvn"
		}

		isInstalled := system.CommandExists(cmdToCheck)

		if action == "install" {
			if isInstalled {
				unavailablePkgs = append(unavailablePkgs, choice)
			} else {
				availableOptions = append(availableOptions, huh.NewOption(choice, choice))
			}
		} else if action == "update" || action == "uninstall" {
			if isInstalled {
				availableOptions = append(availableOptions, huh.NewOption(choice, choice))
			} else {
				unavailablePkgs = append(unavailablePkgs, choice)
			}
		}
	}

	var fields []huh.Field

	// If there are packages that don't apply, display them in a Note (read-only)
	if len(unavailablePkgs) > 0 {
		title := "✨ Already Installed"
		if action == "update" || action == "uninstall" {
			title = "❌ Not Installed"
		}
		fields = append(fields, huh.NewNote().
			Title(title).
			Description(strings.Join(unavailablePkgs, ", ")))
	}

	// If there are still packages left to process, show the multiselect
	if len(availableOptions) > 0 {
		title := fmt.Sprintf("📦 Select packages to %s", action)
		fields = append(fields, huh.NewMultiSelect[string]().
			Title(title).
			Description("Use Up/Down arrows to navigate, Space to select, Enter to confirm.").
			Options(availableOptions...).
			Value(&selected))
	} else {
		// If no packages are available for the action
		msg := "You have already installed all available packages."
		if action == "update" || action == "uninstall" {
			msg = fmt.Sprintf("No packages are currently installed to %s.", action)
		}
		fields = append(fields, huh.NewNote().
			Title("🎉 All caught up!").
			Description(msg))
	}

	// Create the form
	form := huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(huh.ThemeCatppuccin())

	// Run the form
	err := form.Run()
	if err != nil {
		// Handle user cancellation (e.g., pressing Esc or Ctrl+C)
		if err == huh.ErrUserAborted {
			return nil, fmt.Errorf("user aborted")
		}
		return nil, err
	}

	return selected, nil
}
