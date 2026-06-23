package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// SetupGitAndSSH sets up global git config and ensures an SSH key exists.
// Returns the public SSH key string.
func SetupGitAndSSH(name, email string) (string, error) {
	// Configure git user.name
	cmd := exec.Command("git", "config", "--global", "user.name", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to set git user.name: %w: %s", err, string(out))
	}

	// Configure git user.email
	cmd = exec.Command("git", "config", "--global", "user.email", email)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to set git user.email: %w: %s", err, string(out))
	}

	// Setup SSH key
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home dir: %w", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("could not create .ssh directory: %w", err)
	}

	privKeyPath := filepath.Join(sshDir, "id_ed25519")
	pubKeyPath := privKeyPath + ".pub"

	// Ensure SSH keypair exists
	if _, err := os.Stat(privKeyPath); err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("could not stat SSH private key: %w", err)
		}

		cmd = exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", privKeyPath, "-N", "")
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("failed to generate SSH key: %w: %s", err, string(out))
		}
	}

	// Read public key (private key may exist while the .pub file is missing)
	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			cmd = exec.Command("ssh-keygen", "-y", "-f", privKeyPath)
			out, derr := cmd.CombinedOutput()
			if derr != nil {
				return "", fmt.Errorf("failed to derive SSH public key: %w: %s", derr, string(out))
			}
			pubKey = out
			if werr := os.WriteFile(pubKeyPath, pubKey, 0644); werr != nil {
				return "", fmt.Errorf("could not write derived public key: %w", werr)
			}
		} else {
			return "", fmt.Errorf("could not read public key: %w", err)
		}
	}

	return string(pubKey), nil
}
