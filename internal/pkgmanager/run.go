package pkgmanager

import (
	"os"
	"os/exec"
)

// RunInteractive executes a command directly attached to the system's standard streams.
// This allows native prompts (like sudo password requests or pacman confirmations)
// to work perfectly without being intercepted or broken by Go's background execution
// or Bubbletea's terminal handling.
func RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
