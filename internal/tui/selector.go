package tui

import (
	"fmt"
	"strings"

	"dev/internal/registry"

	"github.com/charmbracelet/huh"
)

// RunSelector opens an interactive TUI form to select items from a list
// based on the provided action (install, update, uninstall) and returns the chosen item names.
func RunSelector(action string, choices []string) ([]string, error) {
	var selected []string

	var availableOptions []huh.Option[string]
	var unavailablePkgs []string

	// Map choices to registry packages for easy lookup
	pkgMap := make(map[string]registry.Package)
	for _, p := range registry.Packages {
		pkgMap[p.ID] = p
	}

	// Determine which packages are applicable based on the action
	for _, choice := range choices {
		pkg, ok := pkgMap[choice]
		if !ok {
			continue // Skip unknown packages
		}

		isInstalled := pkg.IsInstalled()

		if action == "install" {
			if isInstalled {
				unavailablePkgs = append(unavailablePkgs, pkg.DisplayName)
			} else {
				availableOptions = append(availableOptions, huh.NewOption(pkg.DisplayName, pkg.ID))
			}
		} else if action == "update" || action == "uninstall" {
			if isInstalled {
				availableOptions = append(availableOptions, huh.NewOption(pkg.DisplayName, pkg.ID))
			} else {
				unavailablePkgs = append(unavailablePkgs, pkg.DisplayName)
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
			Description("Use Up/Down arrows to navigate, Space to select, Enter to confirm. Press Esc to go back.").
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
