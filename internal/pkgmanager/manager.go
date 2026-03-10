package pkgmanager

import (
	"fmt"
)

// UpdateAll upgrades all managed packages and system dependencies.
func UpdateAll() error {
	// Update system packages via pacman/AUR helper
	if err := PacmanUpdateAll(); err != nil {
		return err
	}

	// Attempt to update standalone tools (ignore errors if not installed)
	_ = UpdateBun()
	_ = UpdateNode()

	return nil
}

// CleanAll cleans up the system cache (e.g., pacman cache).
func CleanAll() error {
	// We use pacman directly for cache cleaning as AUR helpers typically wrap this anyway
	err := RunInteractive("sudo", "pacman", "-Scc", "--noconfirm")
	if err != nil {
		return fmt.Errorf("package cache clean failed: %v", err)
	}
	return nil
}
