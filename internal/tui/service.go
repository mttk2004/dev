package tui

import (
	"fmt"

	"dev/internal/registry"
	"dev/internal/system"

	"github.com/charmbracelet/huh"
)

// RunServiceManager displays an interactive menu to start, stop, enable, or disable system services.
func RunServiceManager() error {
	for {
		var options []huh.Option[string]
		hasInstalledServices := false

		for _, pkg := range registry.Packages {
			if len(pkg.Services) == 0 {
				continue
			}
			if !pkg.IsInstalled() {
				continue // Skip if not installed
			}

			for _, srvName := range pkg.Services {
				hasInstalledServices = true
				status := system.GetServiceStatus(srvName)

				activeIcon := "🔴"
				if status.Active {
					activeIcon = "🟢"
				}

				enabledIcon := "[Disabled]"
				if status.Enabled {
					enabledIcon = "[Enabled]"
				}

				label := fmt.Sprintf("%s %-12s %s", activeIcon, srvName, enabledIcon)
				options = append(options, huh.NewOption(label, srvName))
			}
		}

		if !hasInstalledServices {
			// Inform user that no manageble services are installed yet
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewNote().
						Title("⚙️ Manage Services").
						Description("You haven't installed any background services (Docker, PostgreSQL, Redis, Nginx, MariaDB) yet."),
				),
			).WithTheme(huh.ThemeCatppuccin())
			return form.Run()
		}

		options = append(options, huh.NewOption("🚪 Back to Main Menu", "back"))

		var selectedService string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("⚙️  Manage Services").
					Description("Select a service to change its state").
					Options(options...).
					Value(&selectedService),
			),
		).WithTheme(huh.ThemeCatppuccin())

		err := form.Run()
		if err != nil || selectedService == "back" {
			if err == huh.ErrUserAborted {
				return nil
			}
			return err
		}

		// Choose action for the selected service
		var action string
		actionForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("⚙️  Action for %s", selectedService)).
					Options(
						huh.NewOption("▶️  Start", "start"),
						huh.NewOption("⏹️  Stop", "stop"),
						huh.NewOption("🔄 Enable (Start on boot)", "enable"),
						huh.NewOption("❌ Disable", "disable"),
						huh.NewOption("🔙 Cancel", "cancel"),
					).
					Value(&action),
			),
		).WithTheme(huh.ThemeCatppuccin())

		err = actionForm.Run()
		if err != nil || action == "cancel" {
			if err == huh.ErrUserAborted {
				continue // Go back to service selection
			}
			return err
		}

		// Execute action
		fmt.Printf("\nExecuting %s on %s...\n", action, selectedService)
		switch action {
		case "start":
			err = system.StartService(selectedService)
		case "stop":
			err = system.StopService(selectedService)
		case "enable":
			err = system.EnableService(selectedService)
		case "disable":
			err = system.DisableService(selectedService)
		}

		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully executed %s on %s.\n", action, selectedService)
		}

		// Prompt user to continue to refresh the list
		var ack bool
		huh.NewConfirm().Title("Press Enter to continue...").Value(&ack).Run()
	}
}
