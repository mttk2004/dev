package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"dev/internal/system"
)

// ProjectType represents the type of project to scaffold.
type ProjectType string

const (
	ProjectNextJS    ProjectType = "nextjs"
	ProjectLaravel   ProjectType = "laravel"
	ProjectGoAPI     ProjectType = "go-api"
	ProjectReactVite   ProjectType = "react-vite"
	ProjectVueVite     ProjectType = "vue-vite"
	ProjectExpress     ProjectType = "express"
	ProjectDjango      ProjectType = "django"
	ProjectSpringBoot  ProjectType = "spring-boot"
	ProjectReactRouter ProjectType = "react-router"
)

// CreateProject scaffolds a new project of the given type with the given name in the specified parent directory.
func CreateProject(pType ProjectType, parentDir, name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if parentDir == "" {
		parentDir = "." // Default to current directory
	}

	fullPath := filepath.Join(parentDir, name)

	// Ensure the target directory does not already exist
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", fullPath)
	}

	// Ensure the parent directory exists
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %v", err)
	}

	switch pType {
	case ProjectNextJS:
		return scaffoldNextJS(parentDir, name)
	case ProjectLaravel:
		return scaffoldLaravel(parentDir, name)
	case ProjectGoAPI:
		return scaffoldGoAPI(parentDir, name)
	case ProjectReactVite:
		return scaffoldReactVite(parentDir, name)
	case ProjectVueVite:
		return scaffoldVueVite(parentDir, name)
	case ProjectExpress:
		return scaffoldExpress(parentDir, name)
	case ProjectDjango:
		return scaffoldDjango(parentDir, name)
	case ProjectSpringBoot:
		return scaffoldSpringBoot(parentDir, name)
	case ProjectReactRouter:
		return scaffoldReactRouter(parentDir, name)
	default:
		return fmt.Errorf("unsupported project type: %s", pType)
	}
}

func scaffoldNextJS(parentDir, name string) error {
	if !system.CommandExists("bun") {
		return fmt.Errorf("bun is required to scaffold Next.js. Please install it first")
	}
	cmd := exec.Command("bun", "create", "next-app", name)
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding Next.js project '%s' using bun", name))
}

func scaffoldLaravel(parentDir, name string) error {
	if !system.CommandExists("composer") {
		return fmt.Errorf("composer is required to scaffold Laravel. Please install it first")
	}
	cmd := exec.Command("composer", "create-project", "laravel/laravel", name)
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding Laravel project '%s' using composer", name))
}

func scaffoldGoAPI(parentDir, name string) error {
	if !system.CommandExists("go") {
		return fmt.Errorf("go is required to scaffold a Go API. Please install it first")
	}

	fullPath := filepath.Join(parentDir, name)
	fmt.Printf("\n🚀 Scaffolding Go API project '%s' at '%s'...\n", name, fullPath)

	// Create root directory
	if err := os.Mkdir(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create standard Go project layout
	dirs := []string{
		"cmd/api",
		"internal/handler",
		"internal/service",
		"internal/repository",
		"pkg",
		"config",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(fullPath, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize go modules
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = fullPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod init failed: %v", err)
	}

	// Create a basic main.go file
	mainContent := fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"🚀 Hello from %s API!\")\n}\n", name)
	if err := os.WriteFile(filepath.Join(fullPath, "cmd/api/main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}

	fmt.Println("✅ Go API scaffolding completed successfully!")
	return nil
}

func scaffoldReactVite(parentDir, name string) error {
	if !system.CommandExists("bun") {
		return fmt.Errorf("bun is required. Please install it first")
	}
	cmd := exec.Command("bun", "create", "vite", name, "--template", "react-ts")
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding React (Vite) project '%s'", name))
}

func scaffoldVueVite(parentDir, name string) error {
	if !system.CommandExists("bun") {
		return fmt.Errorf("bun is required. Please install it first")
	}
	cmd := exec.Command("bun", "create", "vite", name, "--template", "vue-ts")
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding Vue (Vite) project '%s'", name))
}

func scaffoldExpress(parentDir, name string) error {
	if !system.CommandExists("bun") {
		return fmt.Errorf("bun is required. Please install it first")
	}
	cmd := exec.Command("bunx", "express-generator", name)
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding Express project '%s'", name))
}

func scaffoldDjango(parentDir, name string) error {
	if !system.CommandExists("python") {
		return fmt.Errorf("python is required. Please install it first")
	}

	fullPath := filepath.Join(parentDir, name)
	fmt.Printf("\n🚀 Scaffolding Django project '%s' at '%s'...\n", name, fullPath)

	if err := os.Mkdir(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create virtualenv
	cmdVenv := exec.Command("python", "-m", "venv", "venv")
	cmdVenv.Dir = fullPath
	if err := runCmdSilent(cmdVenv); err != nil {
		return fmt.Errorf("failed to create virtual environment: %v", err)
	}

	// Install django inside venv
	pipPath := filepath.Join(fullPath, "venv", "bin", "pip")
	cmdPip := exec.Command(pipPath, "install", "django")
	cmdPip.Dir = fullPath
	if err := runCmdSilent(cmdPip); err != nil {
		return fmt.Errorf("failed to install django: %v", err)
	}

	// Run django-admin to start the project
	djangoAdminPath := filepath.Join(fullPath, "venv", "bin", "django-admin")
	cmdDjango := exec.Command(djangoAdminPath, "startproject", "core", ".")
	cmdDjango.Dir = fullPath
	if err := runCmdSilent(cmdDjango); err != nil {
		return fmt.Errorf("failed to start django project: %v", err)
	}

	fmt.Println("✅ Django scaffolding completed successfully!")
	return nil
}

func scaffoldSpringBoot(parentDir, name string) error {
	if !system.CommandExists("curl") {
		return fmt.Errorf("curl is required. Please install it first")
	}
	if !system.CommandExists("unzip") {
		return fmt.Errorf("unzip is required. Please install it first")
	}

	fullPath := filepath.Join(parentDir, name)
	fmt.Printf("\n🚀 Scaffolding Spring Boot project '%s' at '%s'...\n", name, fullPath)

	if err := os.Mkdir(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	zipPath := filepath.Join(fullPath, "demo.zip")
	cmdCurl := exec.Command("curl", "-s", "https://start.spring.io/starter.zip", "-d", "name="+name, "-d", "artifactId="+name, "-o", zipPath)
	cmdCurl.Dir = fullPath
	if err := runCmdSilent(cmdCurl); err != nil {
		return fmt.Errorf("failed to download Spring Boot template: %v", err)
	}

	cmdUnzip := exec.Command("unzip", "-q", "demo.zip")
	cmdUnzip.Dir = fullPath
	if err := runCmdSilent(cmdUnzip); err != nil {
		return fmt.Errorf("failed to unzip Spring Boot template: %v", err)
	}

	os.Remove(zipPath)

	fmt.Println("✅ Spring Boot scaffolding completed successfully!")
	return nil
}

func scaffoldReactRouter(parentDir, name string) error {
	if !system.CommandExists("npx") {
		return fmt.Errorf("npx is required. Please install node first")
	}
	cmd := exec.Command("npx", "create-react-router@latest", name)
	cmd.Dir = parentDir
	return runCmd(cmd, fmt.Sprintf("🚀 Scaffolding React Router project '%s'", name))
}

// runCmd executes a command with standard IO attached.
func runCmd(cmd *exec.Cmd, msg string) error {
	fmt.Printf("\n%s...\n", msg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCmdSilent executes a command capturing output only on error.
func runCmdSilent(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
