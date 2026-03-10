package tui

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"dev/internal/scaffold"

	"github.com/charmbracelet/huh"
)

// RunScaffoldPrompt opens a TUI form to gather information for a new project.
// It returns the selected project type, parent directory, and project name.
func RunScaffoldPrompt() (scaffold.ProjectType, string, string, error) {
	var pType scaffold.ProjectType
	var pDir string
	var pName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[scaffold.ProjectType]().
				Title("✨ Choose a project template").
				Options(
					huh.NewOption("Next.js (React Framework via Bun)", scaffold.ProjectNextJS),
					huh.NewOption("React Router (via npx)", scaffold.ProjectReactRouter),
					huh.NewOption("React (Vite via Bun)", scaffold.ProjectReactVite),
					huh.NewOption("Vue (Vite via Bun)", scaffold.ProjectVueVite),
					huh.NewOption("Express (Node.js via Bun)", scaffold.ProjectExpress),
					huh.NewOption("Laravel (PHP Framework via Composer)", scaffold.ProjectLaravel),
					huh.NewOption("Django (Python Framework via venv)", scaffold.ProjectDjango),
					huh.NewOption("Spring Boot (Java via start.spring.io)", scaffold.ProjectSpringBoot),
					huh.NewOption("Go API (Standard Layout)", scaffold.ProjectGoAPI),
				).
				Value(&pType),
			huh.NewInput().
				Title("📂 Parent Directory").
				Description("Where should the project be created? (Leave blank for current directory)").
				Placeholder(".").
				Validate(func(s string) error {
					dir := strings.TrimSpace(s)
					if dir == "" || dir == "." {
						return nil // Current directory is always safe
					}

					// Reject obvious traversal patterns early for better UX
					cleaned := filepath.Clean(dir)
					if strings.HasPrefix(cleaned, "..") {
						return errors.New("path cannot traverse above current directory with '..'")
					}

					// Full validation against safe base dirs
					abs, err := filepath.Abs(cleaned)
					if err != nil {
						return fmt.Errorf("invalid path: %v", err)
					}

					// Use a dummy project name to validate the resolved parent
					if err := scaffold.ValidatePath(abs); err != nil {
						return errors.New("path is outside allowed directories (home dir or current working directory)")
					}

					return nil
				}).
				Value(&pDir),
			huh.NewInput().
				Title("📁 Project Name").
				Description("Enter the name of your new project (will create a directory)").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("project name cannot be empty")
					}
					// Prevent path traversal via project name
					if strings.Contains(s, "/") || strings.Contains(s, "\\") {
						return errors.New("project name cannot contain slashes")
					}
					if strings.Contains(s, "..") {
						return errors.New("project name cannot contain '..'")
					}
					// Reject hidden directories
					if strings.HasPrefix(strings.TrimSpace(s), ".") {
						return errors.New("project name cannot start with '.'")
					}
					return nil
				}).
				Value(&pName),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		return "", "", "", err
	}

	dir := strings.TrimSpace(pDir)
	if dir == "" {
		dir = "."
	}

	return pType, dir, strings.TrimSpace(pName), nil
}
