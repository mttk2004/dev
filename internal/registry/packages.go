package registry

import (
	"dev/internal/pkgmanager"
	"dev/internal/system"
)

// Package represents a centralized definition for a managed tool.
type Package struct {
	ID          string   // Unique identifier (e.g., "node")
	DisplayName string   // Display name for the UI (e.g., "Node.js (fnm)")
	CheckCmd    string   // Command to check for existence in PATH (e.g., "fnm" or "node")
	PacmanPkg   string   // Pacman package name (if applicable)
	Services    []string // Associated systemd services (e.g., ["docker"])
	Install     func() error
	Remove      func() error
	Update      func() error
}

// IsInstalled checks if the package's primary command is available in PATH.
func (p Package) IsInstalled() bool {
	return system.CommandExists(p.CheckCmd)
}

// GetVersion returns the version of the package.
func (p Package) GetVersion() string {
	return system.GetCommandVersion(p.CheckCmd)
}

// GetPath returns the installation path of the package.
func (p Package) GetPath() string {
	return system.GetCommandPath(p.CheckCmd)
}

// Packages is the central registry of all available tools.
// Adding a new tool here automatically propagates it to the UI, Doctor, and Managers.
var Packages = []Package{
	{
		ID:          "node",
		DisplayName: "Node.js (fnm)",
		CheckCmd:    "fnm",
		PacmanPkg:   "fnm",
		Install:     pkgmanager.InstallNode,
		Remove:      pkgmanager.RemoveNode,
		Update:      pkgmanager.UpdateNode,
	},
	{
		ID:          "bun",
		DisplayName: "Bun",
		CheckCmd:    "bun",
		PacmanPkg:   "",
		Install:     pkgmanager.InstallBun,
		Remove:      pkgmanager.RemoveBun,
		Update:      pkgmanager.UpdateBun,
	},
	{
		ID:          "composer",
		DisplayName: "Composer",
		CheckCmd:    "composer",
		PacmanPkg:   "composer",
		Install:     func() error { return pkgmanager.PacmanInstall("composer") },
		Remove:      func() error { return pkgmanager.PacmanRemove("composer") },
		Update:      func() error { return pkgmanager.PacmanUpdate("composer") },
	},
	{
		ID:          "jdk",
		DisplayName: "JDK (Java)",
		CheckCmd:    "java",
		PacmanPkg:   "jdk-openjdk",
		Install:     func() error { return pkgmanager.PacmanInstall("jdk-openjdk") },
		Remove:      func() error { return pkgmanager.PacmanRemove("jdk-openjdk") },
		Update:      func() error { return pkgmanager.PacmanUpdate("jdk-openjdk") },
	},
	{
		ID:          "go",
		DisplayName: "Go",
		CheckCmd:    "go",
		PacmanPkg:   "go",
		Install:     func() error { return pkgmanager.PacmanInstall("go") },
		Remove:      func() error { return pkgmanager.PacmanRemove("go") },
		Update:      func() error { return pkgmanager.PacmanUpdate("go") },
	},
	{
		ID:          "php",
		DisplayName: "PHP",
		CheckCmd:    "php",
		PacmanPkg:   "php",
		Install:     func() error { return pkgmanager.PacmanInstall("php") },
		Remove:      func() error { return pkgmanager.PacmanRemove("php") },
		Update:      func() error { return pkgmanager.PacmanUpdate("php") },
	},
	{
		ID:          "docker",
		DisplayName: "Docker",
		CheckCmd:    "docker",
		PacmanPkg:   "docker",
		Services:    []string{"docker"},
		Install:     func() error { return pkgmanager.PacmanInstall("docker") },
		Remove:      func() error { return pkgmanager.PacmanRemove("docker") },
		Update:      func() error { return pkgmanager.PacmanUpdate("docker") },
	},
	{
		ID:          "postgresql",
		DisplayName: "PostgreSQL",
		CheckCmd:    "psql",
		PacmanPkg:   "postgresql",
		Services:    []string{"postgresql"},
		Install:     func() error { return pkgmanager.PacmanInstall("postgresql") },
		Remove:      func() error { return pkgmanager.PacmanRemove("postgresql") },
		Update:      func() error { return pkgmanager.PacmanUpdate("postgresql") },
	},
	{
		ID:          "redis",
		DisplayName: "Redis",
		CheckCmd:    "redis-cli",
		PacmanPkg:   "redis",
		Services:    []string{"redis"},
		Install:     func() error { return pkgmanager.PacmanInstall("redis") },
		Remove:      func() error { return pkgmanager.PacmanRemove("redis") },
		Update:      func() error { return pkgmanager.PacmanUpdate("redis") },
	},
	{
		ID:          "nginx",
		DisplayName: "Nginx",
		CheckCmd:    "nginx",
		PacmanPkg:   "nginx",
		Services:    []string{"nginx"},
		Install:     func() error { return pkgmanager.PacmanInstall("nginx") },
		Remove:      func() error { return pkgmanager.PacmanRemove("nginx") },
		Update:      func() error { return pkgmanager.PacmanUpdate("nginx") },
	},
	{
		ID:          "python",
		DisplayName: "Python",
		CheckCmd:    "python",
		PacmanPkg:   "python",
		Install:     func() error { return pkgmanager.PacmanInstall("python") },
		Remove:      func() error { return pkgmanager.PacmanRemove("python") },
		Update:      func() error { return pkgmanager.PacmanUpdate("python") },
	},
	{
		ID:          "maven",
		DisplayName: "Maven",
		CheckCmd:    "mvn",
		PacmanPkg:   "maven",
		Install:     func() error { return pkgmanager.PacmanInstall("maven") },
		Remove:      func() error { return pkgmanager.PacmanRemove("maven") },
		Update:      func() error { return pkgmanager.PacmanUpdate("maven") },
	},
	{
		ID:          "mariadb",
		DisplayName: "MariaDB",
		CheckCmd:    "mariadb",
		PacmanPkg:   "mariadb",
		Services:    []string{"mariadb"},
		Install:     func() error { return pkgmanager.PacmanInstall("mariadb") },
		Remove:      func() error { return pkgmanager.PacmanRemove("mariadb") },
		Update:      func() error { return pkgmanager.PacmanUpdate("mariadb") },
	},
}
