package pkgmanager

import (
	"fmt"
	"strings"
)

// Install executes the actual installation of a given package.
func Install(pkgName string) error {
	if pkgName == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	switch strings.ToLower(pkgName) {
	case "bun":
		return InstallBun()
	case "node":
		return InstallNode()
	case "composer":
		return PacmanInstall("composer")
	case "jdk":
		return PacmanInstall("jdk-openjdk")
	case "go":
		return PacmanInstall("go")
	case "php":
		return PacmanInstall("php")
	case "docker":
		return PacmanInstall("docker")
	case "postgresql":
		return PacmanInstall("postgresql")
	case "redis":
		return PacmanInstall("redis")
	case "nginx":
		return PacmanInstall("nginx")
	case "python":
		return PacmanInstall("python")
	case "maven":
		return PacmanInstall("maven")
	case "mariadb":
		return PacmanInstall("mariadb")
	default:
		return fmt.Errorf("unsupported package: %s", pkgName)
	}
}

// Update upgrades a given package.
func Update(pkgName string) error {
	if pkgName == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	switch strings.ToLower(pkgName) {
	case "bun":
		return UpdateBun()
	case "node":
		return UpdateNode()
	case "composer":
		return PacmanUpdate("composer")
	case "jdk":
		return PacmanUpdate("jdk-openjdk")
	case "go":
		return PacmanUpdate("go")
	case "php":
		return PacmanUpdate("php")
	case "docker":
		return PacmanUpdate("docker")
	case "postgresql":
		return PacmanUpdate("postgresql")
	case "redis":
		return PacmanUpdate("redis")
	case "nginx":
		return PacmanUpdate("nginx")
	case "python":
		return PacmanUpdate("python")
	case "maven":
		return PacmanUpdate("maven")
	case "mariadb":
		return PacmanUpdate("mariadb")
	default:
		return fmt.Errorf("unsupported package: %s", pkgName)
	}
}

// UpdateAll upgrades all managed packages and system dependencies.
func UpdateAll() error {
	// Update system packages via pacman
	if err := PacmanUpdateAll(); err != nil {
		return err
	}

	// Attempt to update standalone tools like Bun (ignore errors if not installed)
	_ = UpdateBun()

	return nil
}

// Remove uninstalls a given package.
func Remove(pkgName string) error {
	if pkgName == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	switch strings.ToLower(pkgName) {
	case "bun":
		return RemoveBun()
	case "node":
		return RemoveNode()
	case "composer":
		return PacmanRemove("composer")
	case "jdk":
		return PacmanRemove("jdk-openjdk")
	case "go":
		return PacmanRemove("go")
	case "php":
		return PacmanRemove("php")
	case "docker":
		return PacmanRemove("docker")
	case "postgresql":
		return PacmanRemove("postgresql")
	case "redis":
		return PacmanRemove("redis")
	case "nginx":
		return PacmanRemove("nginx")
	case "python":
		return PacmanRemove("python")
	case "maven":
		return PacmanRemove("maven")
	case "mariadb":
		return PacmanRemove("mariadb")
	default:
		return fmt.Errorf("unsupported package: %s", pkgName)
	}
}

// CleanAll cleans up the system cache (e.g., pacman cache).
func CleanAll() error {
	err := RunInteractive("sudo", "pacman", "-Scc", "--noconfirm")
	if err != nil {
		return fmt.Errorf("pacman cache clean failed: %v", err)
	}
	return nil
}
