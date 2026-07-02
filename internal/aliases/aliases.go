package aliases

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Alias represents a key-value shell command alias.
type Alias struct {
	Name  string
	Value string
}

func getZshrcPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".zshrc"), nil
}

// GetAliases reads ~/.zshrc and returns all defined aliases.
func GetAliases() ([]Alias, error) {
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(zshrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var aliases []Alias
	re := regexp.MustCompile(`^alias\s+([a-zA-Z0-9_\-]+)\s*=\s*['"]?(.*?)['"]?\s*$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "alias ") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			aliases = append(aliases, Alias{
				Name:  matches[1],
				Value: matches[2],
			})
		}
	}

	return aliases, scanner.Err()
}

// AddAlias adds or updates a shell alias in ~/.zshrc.
func AddAlias(name, value string) error {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" || value == "" {
		return fmt.Errorf("name and value cannot be empty")
	}

	zshrcPath, err := getZshrcPath()
	if err != nil {
		return err
	}

	var lines []string
	exists := false
	newLine := fmt.Sprintf("alias %s='%s'", name, value)

	file, err := os.Open(zshrcPath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(fmt.Sprintf(`^alias\s+%s\s*=`, regexp.QuoteMeta(name)))
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(strings.TrimSpace(line)) {
				lines = append(lines, newLine)
				exists = true
			} else {
				lines = append(lines, line)
			}
		}
		file.Close()
	}

	if !exists {
		f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString("\n" + newLine + "\n")
		return err
	}

	return os.WriteFile(zshrcPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// RemoveAlias deletes a shell alias from ~/.zshrc.
func RemoveAlias(name string) error {
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return err
	}

	file, err := os.Open(zshrcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	re := regexp.MustCompile(fmt.Sprintf(`^alias\s+%s\s*=`, regexp.QuoteMeta(name)))
	scanner := bufio.NewScanner(file)
	removed := false

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(strings.TrimSpace(line)) {
			removed = true
			continue
		}
		lines = append(lines, line)
	}

	if !removed {
		return fmt.Errorf("alias %s not found", name)
	}

	return os.WriteFile(zshrcPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
