package tui

import (
	"github.com/charmbracelet/huh"
)

// DashboardAction represents the user's choice in the main menu.
type DashboardAction string

const (
	ActionInstall   DashboardAction = "install"
	ActionUpdate    DashboardAction = "update"
	ActionUninstall DashboardAction = "uninstall"
	ActionJDK       DashboardAction = "jdk"
	ActionClean     DashboardAction = "clean"
	ActionSearch    DashboardAction = "search"
	ActionServices  DashboardAction = "services"
	ActionScaffold  DashboardAction = "scaffold"
	ActionConfig    DashboardAction = "config"
	ActionExit      DashboardAction = "exit"
)

// RunDashboard displays the main interactive menu for the dev tool.
func RunDashboard() (DashboardAction, error) {
	var action DashboardAction

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[DashboardAction]().
				Title("🛠️  Main Menu").
				Description("Choose an action to manage your web dev environment").
				Options(
					huh.NewOption("📦 Install packages", ActionInstall),
					huh.NewOption("🔄 Update packages", ActionUpdate),
					huh.NewOption("🧹 Uninstall packages", ActionUninstall),
					huh.NewOption("☕ Manage JDK Versions", ActionJDK),
					huh.NewOption("🧼 Clean Dev Caches", ActionClean),
					huh.NewOption("🔍 Search for a package", ActionSearch),
					huh.NewOption("⚙️  Manage Services", ActionServices),
					huh.NewOption("✨ Create New Project", ActionScaffold),
					huh.NewOption("🔧 Configuration", ActionConfig),
					huh.NewOption("🚪 Exit", ActionExit),
				).
				Value(&action),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return ActionExit, nil
		}
		return ActionExit, err
	}

	return action, nil
}

// RunSearchPrompt displays a text input for the user to search for a package.
func RunSearchPrompt() (string, error) {
	var query string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔍 Search Package").
				Description("Enter a package name to search in Arch Linux repositories").
				Placeholder("e.g. neovim, tmux, apache...").
				Value(&query),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		return "", err
	}

	return query, nil
}
