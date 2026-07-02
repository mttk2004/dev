package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
				Description("Use Up/Down arrows to navigate, Enter to confirm. Press Esc to go back.").
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
	sysName, sysEmail := system.GetGlobalGitConfig()

	// Check if SSH key exists
	hasSSHKey := false
	if home, err := os.UserHomeDir(); err == nil {
		if _, err := os.Stat(filepath.Join(home, ".ssh", "id_ed25519")); err == nil {
			hasSSHKey = true
		}
	}

	// If all are already configured, ask before overwriting
	if sysName != "" && sysEmail != "" && hasSSHKey {
		var confirm bool
		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Git & SSH đã được cấu hình trước đó.\n(Name: %s, Email: %s)\nBạn có muốn cấu hình lại không?", sysName, sysEmail)).
					Description("Press Enter to confirm, Esc to cancel.").
					Value(&confirm),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := confirmForm.Run(); err != nil || !confirm {
			return false
		}
	}

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

	// Fallback to system-level configuration as default values
	if name == "" {
		name = sysName
	}
	if email == "" {
		email = sysEmail
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Git User Name").
				Description("Enter your full name for git commits. Press Enter to confirm, Esc to go back.").
				Value(&name),
			huh.NewInput().
				Title("Git User Email").
				Description("Enter your email address for git commits and SSH key. Press Enter to confirm, Esc to go back.").
				Value(&email),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		return false // Aborted
	}

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

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
