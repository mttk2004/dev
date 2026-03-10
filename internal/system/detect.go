package system

import (
	"os/exec"
	"strings"
)

// CommandExists checks if a specific command is available in the system's PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetCommandPath returns the absolute path of a command if it exists in PATH.
func GetCommandPath(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return ""
	}
	return path
}

// GetCommandVersion attempts to get the version string of a command.
func GetCommandVersion(cmd string) string {
	var out []byte
	var err error

	// Special cases for commands that need specific flags
	switch cmd {
	case "java":
		out, err = exec.Command(cmd, "-version").CombinedOutput()
	case "go":
		out, err = exec.Command(cmd, "version").CombinedOutput()
	default:
		out, err = exec.Command(cmd, "--version").CombinedOutput()
		if err != nil {
			// Fallback to -v
			out, err = exec.Command(cmd, "-v").CombinedOutput()
		}
	}

	if err != nil && len(out) == 0 {
		return "unknown"
	}

	// Extract the first line of the output
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		// Truncate if it's too long to keep the table clean
		versionStr := lines[0]
		if len(versionStr) > 40 {
			return versionStr[:37] + "..."
		}
		return versionStr
	}
	return "unknown"
}

// GetPackageManager returns the best available package manager for Arch Linux.
// It prioritizes AUR helpers (yay, paru) over default pacman.
func GetPackageManager() string {
	if CommandExists("yay") {
		return "yay"
	}
	if CommandExists("paru") {
		return "paru"
	}
	return "pacman"
}
