package pkgmanager

import (
	"fmt"

	"dev/internal/system"
)

// PacmanInstall installs a package using Arch Linux's package manager (yay/paru/pacman).
func PacmanInstall(pkg string) error {
	pm := system.GetPackageManager()
	var err error

	if pm == "pacman" {
		err = RunInteractive("sudo", "pacman", "-S", "--noconfirm", pkg)
	} else {
		// AUR helpers don't need sudo
		err = RunInteractive(pm, "-S", "--noconfirm", pkg)
	}

	if err != nil {
		return fmt.Errorf("pacman install failed for %s: %v", pkg, err)
	}
	return nil
}

// PacmanRemove uninstalls a package and its unused dependencies.
func PacmanRemove(pkg string) error {
	pm := system.GetPackageManager()
	var err error

	if pm == "pacman" {
		err = RunInteractive("sudo", "pacman", "-Rns", "--noconfirm", pkg)
	} else {
		err = RunInteractive(pm, "-Rns", "--noconfirm", pkg)
	}

	if err != nil {
		return fmt.Errorf("pacman remove failed for %s: %v", pkg, err)
	}
	return nil
}

// PacmanUpdateAll performs a full system upgrade, including AUR packages if available.
func PacmanUpdateAll() error {
	pm := system.GetPackageManager()
	var err error

	if pm == "pacman" {
		err = RunInteractive("sudo", "pacman", "-Syu", "--noconfirm")
	} else {
		err = RunInteractive(pm, "-Syu", "--noconfirm")
	}

	if err != nil {
		return fmt.Errorf("pacman system update failed: %v", err)
	}
	return nil
}

// PacmanUpdate updates a specific package.
func PacmanUpdate(pkg string) error {
	// Installing an already installed package syncs it to the latest version.
	return PacmanInstall(pkg)
}
