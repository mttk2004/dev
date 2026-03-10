package registry

import (
	"fmt"

	"dev/internal/pkgmanager"
	"dev/internal/system"
)

// Package ID constants — the single source of truth for all package identifiers.
// Use these constants everywhere instead of hardcoding string literals.
const (
	IDNode       = "node"
	IDBun        = "bun"
	IDComposer   = "composer"
	IDJDK        = "jdk"
	IDGo         = "go"
	IDPHP        = "php"
	IDDocker     = "docker"
	IDPostgreSQL = "postgresql"
	IDRedis      = "redis"
	IDNginx      = "nginx"
	IDPython     = "python"
	IDMaven      = "maven"
	IDMariaDB    = "mariadb"
)

// Package represents a centralized definition for a managed tool.
type Package struct {
	ID          string   // Unique identifier — use the ID* constants above
	DisplayName string   // Display name for the UI (e.g., "Node.js (fnm)")
	CheckCmd    string   // Command to check for existence in PATH (e.g., "fnm" or "node")
	AltCheckCmd string   // Alternative command to verify (e.g., "node" when CheckCmd is "fnm")
	PacmanPkg   string   // Pacman package name (if applicable)
	Services    []string // Associated systemd services (e.g., ["docker"])
	Install     func() error
	Remove      func() error
	Update      func() error
}

// LookupByID returns a pointer to the Package with the given ID, or an error
// if no such package exists. This replaces ad-hoc map/loop lookups scattered
// across the codebase.
func LookupByID(id string) (*Package, error) {
	for i := range Packages {
		if Packages[i].ID == id {
			return &Packages[i], nil
		}
	}
	return nil, fmt.Errorf("package %q not found in registry", id)
}

// IDs returns a slice of all registered package IDs, useful for building
// TUI choice lists without hardcoding.
func IDs() []string {
	ids := make([]string, len(Packages))
	for i, p := range Packages {
		ids[i] = p.ID
	}
	return ids
}

// IsInstalled checks if the package's primary command is available in PATH.
func (p Package) IsInstalled() bool {
	return system.CommandExists(p.CheckCmd)
}

// GetVersion returns the version of the package.
// If AltCheckCmd is set and available, it is preferred for version detection.
func (p Package) GetVersion() string {
	if p.AltCheckCmd != "" && system.CommandExists(p.AltCheckCmd) {
		return system.GetCommandVersion(p.AltCheckCmd)
	}
	return system.GetCommandVersion(p.CheckCmd)
}

// GetPath returns the installation path of the package.
// If AltCheckCmd is set and available, it is preferred for path detection.
func (p Package) GetPath() string {
	if p.AltCheckCmd != "" && system.CommandExists(p.AltCheckCmd) {
		return system.GetCommandPath(p.AltCheckCmd)
	}
	return system.GetCommandPath(p.CheckCmd)
}

// IsFullyOperational returns true if the package's primary command (and
// alternative command, if any) are both available. For most packages this
// is identical to IsInstalled(), but for node (fnm + node) it catches the
// case where fnm is installed but node is not yet in PATH.
func (p Package) IsFullyOperational() bool {
	if !p.IsInstalled() {
		return false
	}
	if p.AltCheckCmd != "" {
		return system.CommandExists(p.AltCheckCmd)
	}
	return true
}

// DiagnosticMessage returns a human-readable status message for the doctor
// report, encoding all the special-case logic that was previously scattered
// as magic strings in doctor/check.go.
func (p Package) DiagnosticMessage() string {
	if !p.IsInstalled() {
		// Special hint for bun when env is configured but binary missing
		if p.ID == IDBun && system.HasEnvVarInZshrc("BUN_INSTALL") {
			return "BUN_INSTALL is in .zshrc, but bun is not in PATH. Try restarting your terminal."
		}
		return fmt.Sprintf("%s is missing or not in PATH", p.CheckCmd)
	}

	// Installed but alternative command missing (e.g. fnm ok, node missing)
	if p.AltCheckCmd != "" && !system.CommandExists(p.AltCheckCmd) {
		return fmt.Sprintf("%s is installed, but %s is not in PATH. Try restarting your terminal.",
			p.CheckCmd, p.AltCheckCmd)
	}

	if p.AltCheckCmd != "" {
		return fmt.Sprintf("%s and %s are installed and in PATH", p.CheckCmd, p.AltCheckCmd)
	}
	return fmt.Sprintf("%s is installed and in PATH", p.CheckCmd)
}

// Packages is the central registry of all available tools.
// Adding a new tool here automatically propagates it to the UI, Doctor, and Managers.
var Packages = []Package{
	{
		ID:          IDNode,
		DisplayName: "Node.js (fnm)",
		CheckCmd:    "fnm",
		AltCheckCmd: "node",
		PacmanPkg:   "fnm",
		Install:     pkgmanager.InstallNode,
		Remove:      pkgmanager.RemoveNode,
		Update:      pkgmanager.UpdateNode,
	},
	{
		ID:          IDBun,
		DisplayName: "Bun",
		CheckCmd:    "bun",
		PacmanPkg:   "",
		Install:     pkgmanager.InstallBun,
		Remove:      pkgmanager.RemoveBun,
		Update:      pkgmanager.UpdateBun,
	},
	{
		ID:          IDComposer,
		DisplayName: "Composer",
		CheckCmd:    "composer",
		PacmanPkg:   "composer",
		Install:     func() error { return pkgmanager.PacmanInstall("composer") },
		Remove:      func() error { return pkgmanager.PacmanRemove("composer") },
		Update:      func() error { return pkgmanager.PacmanUpdate("composer") },
	},
	{
		ID:          IDJDK,
		DisplayName: "JDK (Java)",
		CheckCmd:    "java",
		PacmanPkg:   "jdk-openjdk",
		Install:     func() error { return pkgmanager.PacmanInstall("jdk-openjdk") },
		Remove:      func() error { return pkgmanager.PacmanRemove("jdk-openjdk") },
		Update:      func() error { return pkgmanager.PacmanUpdate("jdk-openjdk") },
	},
	{
		ID:          IDGo,
		DisplayName: "Go",
		CheckCmd:    "go",
		PacmanPkg:   "go",
		Install:     func() error { return pkgmanager.PacmanInstall("go") },
		Remove:      func() error { return pkgmanager.PacmanRemove("go") },
		Update:      func() error { return pkgmanager.PacmanUpdate("go") },
	},
	{
		ID:          IDPHP,
		DisplayName: "PHP",
		CheckCmd:    "php",
		PacmanPkg:   "php",
		Install:     func() error { return pkgmanager.PacmanInstall("php") },
		Remove:      func() error { return pkgmanager.PacmanRemove("php") },
		Update:      func() error { return pkgmanager.PacmanUpdate("php") },
	},
	{
		ID:          IDDocker,
		DisplayName: "Docker",
		CheckCmd:    "docker",
		PacmanPkg:   "docker",
		Services:    []string{"docker"},
		Install:     func() error { return pkgmanager.PacmanInstall("docker") },
		Remove:      func() error { return pkgmanager.PacmanRemove("docker") },
		Update:      func() error { return pkgmanager.PacmanUpdate("docker") },
	},
	{
		ID:          IDPostgreSQL,
		DisplayName: "PostgreSQL",
		CheckCmd:    "psql",
		PacmanPkg:   "postgresql",
		Services:    []string{"postgresql"},
		Install:     func() error { return pkgmanager.PacmanInstall("postgresql") },
		Remove:      func() error { return pkgmanager.PacmanRemove("postgresql") },
		Update:      func() error { return pkgmanager.PacmanUpdate("postgresql") },
	},
	{
		ID:          IDRedis,
		DisplayName: "Redis",
		CheckCmd:    "redis-cli",
		PacmanPkg:   "redis",
		Services:    []string{"redis"},
		Install:     func() error { return pkgmanager.PacmanInstall("redis") },
		Remove:      func() error { return pkgmanager.PacmanRemove("redis") },
		Update:      func() error { return pkgmanager.PacmanUpdate("redis") },
	},
	{
		ID:          IDNginx,
		DisplayName: "Nginx",
		CheckCmd:    "nginx",
		PacmanPkg:   "nginx",
		Services:    []string{"nginx"},
		Install:     func() error { return pkgmanager.PacmanInstall("nginx") },
		Remove:      func() error { return pkgmanager.PacmanRemove("nginx") },
		Update:      func() error { return pkgmanager.PacmanUpdate("nginx") },
	},
	{
		ID:          IDPython,
		DisplayName: "Python",
		CheckCmd:    "python",
		PacmanPkg:   "python",
		Install:     func() error { return pkgmanager.PacmanInstall("python") },
		Remove:      func() error { return pkgmanager.PacmanRemove("python") },
		Update:      func() error { return pkgmanager.PacmanUpdate("python") },
	},
	{
		ID:          IDMaven,
		DisplayName: "Maven",
		CheckCmd:    "mvn",
		PacmanPkg:   "maven",
		Install:     func() error { return pkgmanager.PacmanInstall("maven") },
		Remove:      func() error { return pkgmanager.PacmanRemove("maven") },
		Update:      func() error { return pkgmanager.PacmanUpdate("maven") },
	},
	{
		ID:          IDMariaDB,
		DisplayName: "MariaDB",
		CheckCmd:    "mariadb",
		PacmanPkg:   "mariadb",
		Services:    []string{"mariadb"},
		Install:     func() error { return pkgmanager.PacmanInstall("mariadb") },
		Remove:      func() error { return pkgmanager.PacmanRemove("mariadb") },
		Update:      func() error { return pkgmanager.PacmanUpdate("mariadb") },
	},
}
