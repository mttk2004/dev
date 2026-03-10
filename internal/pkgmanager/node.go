package pkgmanager

import (
	"fmt"
	"os/exec"

	"dev/internal/system"
)

// InstallNode installs FNM (Fast Node Manager) and the LTS version of Node.js.
func InstallNode() error {
	// 1. Install fnm via Arch Linux pacman
	if err := PacmanInstall("fnm"); err != nil {
		return fmt.Errorf("failed to install fnm: %v", err)
	}

	// 2. Add fnm environment setup to .zshrc
	zshrcLine := `eval "$(fnm env --use-on-cd --shell zsh)"`
	if err := system.AppendToZshrc(zshrcLine); err != nil {
		return fmt.Errorf("failed to setup fnm in .zshrc: %v", err)
	}

	// 3. Install Node.js LTS using fnm
	// Since fnm is installed via pacman, it should be available in PATH immediately.
	cmdInstall := exec.Command("fnm", "install", "--lts")
	if output, err := cmdInstall.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Node.js LTS via fnm: %v\nOutput: %s", err, string(output))
	}

	// 4. Set the LTS version as default
	cmdDefault := exec.Command("fnm", "default", "lts-latest")
	if output, err := cmdDefault.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set default Node.js version: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// UpdateNode updates fnm and installs the latest LTS Node.js version.
func UpdateNode() error {
	// Update fnm package itself
	if err := PacmanUpdate("fnm"); err != nil {
		return err
	}

	// Fetch and install the newest LTS release
	cmdInstall := exec.Command("fnm", "install", "--lts")
	if output, err := cmdInstall.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update Node.js LTS via fnm: %v\nOutput: %s", err, string(output))
	}

	cmdDefault := exec.Command("fnm", "default", "lts-latest")
	if output, err := cmdDefault.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set default Node.js version: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// RemoveNode uninstalls fnm and its configurations.
func RemoveNode() error {
	// Uninstall fnm package
	if err := PacmanRemove("fnm"); err != nil {
		return err
	}

	// Note: We leave the downloaded node versions in ~/.local/share/fnm
	// and the .zshrc eval line intact as a safe default,
	// but the user might want to clean them up manually.
	return nil
}
