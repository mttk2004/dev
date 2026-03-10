package pkgmanager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"dev/internal/system"
)

// InstallBun installs Bun using the official curl script and configures zsh.
func InstallBun() error {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://bun.sh/install | bash")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install bun: %v\nOutput: %s", err, string(output))
	}

	// Add to .zshrc
	err = system.AppendToZshrc(`export BUN_INSTALL="$HOME/.bun"`)
	if err != nil {
		return fmt.Errorf("failed to set BUN_INSTALL in .zshrc: %v", err)
	}

	err = system.AppendToZshrc(`export PATH="$BUN_INSTALL/bin:$PATH"`)
	if err != nil {
		return fmt.Errorf("failed to add bun to PATH in .zshrc: %v", err)
	}

	return nil
}

// UpdateBun upgrades Bun to the latest version.
func UpdateBun() error {
	// Bun has a built-in upgrade command
	cmd := exec.Command("bun", "upgrade")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update bun: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// RemoveBun uninstalls Bun by removing its directory.
func RemoveBun() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	bunDir := filepath.Join(home, ".bun")
	err = os.RemoveAll(bunDir)
	if err != nil {
		return fmt.Errorf("failed to remove bun directory (%s): %v", bunDir, err)
	}

	// Note: We leave the .zshrc exports as cleaning them automatically can be risky.
	// The user can manually remove them if desired.
	return nil
}
