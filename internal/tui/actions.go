package tui

import (
	"fmt"

	"dev/internal/pkgmanager"
	"dev/internal/registry"
	"dev/internal/scaffold"
	"dev/internal/system"
	"dev/internal/ui"
)

// RunInstallAction handles the full install flow: selector → install each selected package.
// Returns true if any action was performed (caller should pause for user to read output).
func RunInstallAction(choices []string) bool {
	selected, err := RunSelector("install", choices)
	if err != nil || len(selected) == 0 {
		return false
	}

	for _, pkgID := range selected {
		p, err := registry.LookupByID(pkgID)
		if err != nil {
			ui.Warning("%v", err)
			continue
		}
		ui.ActionHeader("📦", "Installing", p.DisplayName)
		ui.ActionResult(p.DisplayName, p.Install(), "install")
	}

	return true
}

// RunUpdateAction handles the full update flow: selector → update each selected package.
// Returns true if any action was performed.
func RunUpdateAction(choices []string) bool {
	selected, err := RunSelector("update", choices)
	if err != nil || len(selected) == 0 {
		return false
	}

	for _, pkgID := range selected {
		p, err := registry.LookupByID(pkgID)
		if err != nil {
			ui.Warning("%v", err)
			continue
		}
		ui.ActionHeader("🔄", "Updating", p.DisplayName)
		ui.ActionResult(p.DisplayName, p.Update(), "update")
	}

	return true
}

// RunUninstallAction handles the full uninstall flow: selector → remove each selected package.
// Returns true if any action was performed.
func RunUninstallAction(choices []string) bool {
	selected, err := RunSelector("uninstall", choices)
	if err != nil || len(selected) == 0 {
		return false
	}

	for _, pkgID := range selected {
		p, err := registry.LookupByID(pkgID)
		if err != nil {
			ui.Warning("%v", err)
			continue
		}
		ui.ActionHeader("🧹", "Uninstalling", p.DisplayName)
		ui.ActionResult(p.DisplayName, p.Remove(), "uninstall")
	}

	return true
}

// RunSearchAction handles the package search flow: prompt → search via pacman/AUR helper.
// Returns true if any action was performed.
func RunSearchAction() bool {
	query, err := RunSearchPrompt()
	if err != nil || query == "" {
		return false
	}

	pm := system.GetPackageManager()
	fmt.Printf("\n🔍 Searching for '%s' using %s...\n", query, pm)

	if err := pkgmanager.RunInteractive(pm, "-Ss", query); err != nil {
		ui.Subtle("Search finished.")
	}

	return true
}

// RunServiceAction handles the service management flow.
// Returns true if an error occurred that needs user attention (the service
// manager has its own internal pause loop for normal operation).
func RunServiceAction() bool {
	if err := RunServiceManager(); err != nil {
		ui.Error("Error managing services: %v", err)
		return true
	}
	return false
}

// RunScaffoldAction handles the project scaffolding flow: prompt → create project.
// Returns true if any action was performed.
func RunScaffoldAction() bool {
	pType, pDir, pName, err := RunScaffoldPrompt()
	if err != nil {
		if err.Error() != "user aborted" {
			ui.Error("Error scaffolding project: %v", err)
			return true
		}
		return false
	}

	if err := scaffold.CreateProject(pType, pDir, pName); err != nil {
		ui.Error("Failed to create project: %v", err)
	}

	return true
}
