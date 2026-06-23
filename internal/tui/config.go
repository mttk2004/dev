package tui

import (
	"fmt"

	"dev/internal/system"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

// RunConfigAction handles the configuration menu.
// Returns true if an action was performed that requires pausing for the user.
func RunConfigAction() bool {
	var selectedAction string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🔧 Configuration").
				Description("Select an aspect to configure").
				Options(
					huh.NewOption("Git Configuration", "git"),
					huh.NewOption("🚪 Back to Main Menu", "back"),
				).
				Value(&selectedAction),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil || selectedAction == "back" {
		return false
	}

	if selectedAction == "git" {
		return RunGitConfig()
	}

	return false
}

// RunGitConfig handles the Git/SSH setup flow.
func RunGitConfig() bool {
	cfg, err := system.LoadConfig()
	if err != nil {
		ui.Warning("Could not load config: %v", err)
		// We can still continue with empty config
	}

	var name string
	var email string

	if cfg != nil {
		name = cfg.Git.Name
		email = cfg.Git.Email
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Git User Name").
				Description("Enter your full name for git commits").
				Value(&name),
			huh.NewInput().
				Title("Git User Email").
				Description("Enter your email address for git commits and SSH key").
				Value(&email),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		return false // Aborted
	}

	if name == "" || email == "" {
		ui.Warning("Name and Email cannot be empty.")
		return true
	}

	// Save the config
	if cfg == nil {
		cfg = &system.Config{}
	}
	cfg.Git.Name = name
	cfg.Git.Email = email

	if err := system.SaveConfig(cfg); err != nil {
		ui.Error("Failed to save config: %v", err)
		// Try to continue anyway
	}

	// Execute Setup
	ui.ActionHeader("🔧", "Configuring", "Git and SSH")
	fmt.Println("Running git config --global...")
	fmt.Println("Ensuring SSH ed25519 key exists...")

	pubKey, err := system.SetupGitAndSSH(name, email)

	if err != nil {
		ui.ActionResult("Git & SSH Setup", err, "setup")
		return true
	}

	ui.ActionResult("Git & SSH Setup", nil, "setup")

	fmt.Println("\n🔑 Your Public SSH Key (id_ed25519.pub):")
	fmt.Println("----------------------------------------------------------------------")
	fmt.Print(pubKey)
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("You can copy the above key and paste it into GitHub -> Settings -> SSH and GPG keys.")

	return true
}
