package tui

import (
	"errors"
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
				Value(&pDir),
			huh.NewInput().
				Title("📁 Project Name").
				Description("Enter the name of your new project (will create a directory)").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("project name cannot be empty")
					}
					// Simple validation to prevent path traversal issues or invalid chars
					if strings.Contains(s, "/") || strings.Contains(s, "\\") {
						return errors.New("project name cannot contain slashes")
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
