package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppendToZshrc appends a given line to the user's ~/.zshrc file if it doesn't already exist.
func AppendToZshrc(line string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	zshrcPath := filepath.Join(home, ".zshrc")

	// Check if line already exists
	exists, err := lineExistsInFile(zshrcPath, line)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not read .zshrc: %v", err)
	}

	if exists {
		return nil // Already exists, nothing to do
	}

	// Append to file
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open .zshrc for writing: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + line + "\n"); err != nil {
		return fmt.Errorf("could not write to .zshrc: %v", err)
	}

	return nil
}

// lineExistsInFile checks if an exact line exists in a file.
func lineExistsInFile(filePath, searchLine string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(searchLine) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// HasEnvVarInZshrc checks if a specific environment variable or string exists in ~/.zshrc.
func HasEnvVarInZshrc(envStr string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	zshrcPath := filepath.Join(home, ".zshrc")
	f, err := os.Open(zshrcPath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), envStr) {
			return true
		}
	}

	return false
}
