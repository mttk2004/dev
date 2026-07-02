package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"dev/internal/docker"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

type containerAction string

const (
	actionStart   containerAction = "start"
	actionStop    containerAction = "stop"
	actionRestart containerAction = "restart"
	actionLogs    containerAction = "logs"
	actionBack    containerAction = "back"
)

// RunDockerAction handles the main Docker dashboard menu.
func RunDockerAction() bool {
	ui.ClearScreen()
	fmt.Println("🐳 Docker Container Dashboard")
	fmt.Println("----------------------------")

	if !docker.IsDockerInstalled() {
		ui.Warning("Docker is not installed on this system.")
		fmt.Println("You can install Docker from the '📦 Install packages' menu.")
		pauseClean()
		return false
	}

	if !docker.IsDockerRunning() {
		ui.Warning("Docker daemon is not running.")
		var startService bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Would you like to start the Docker systemd service?").
					Description("Press Enter to confirm, Esc to cancel.").
					Value(&startService),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := form.Run(); err == nil && startService {
			ui.ActionHeader("⚙️", "Starting", "docker.service")
			err := exec.Command("sudo", "systemctl", "start", "docker").Run()
			ui.ActionResult("Docker service", err, "start")
			pauseClean()
			if err != nil {
				return false
			}
		} else {
			return false
		}
	}

	for {
		ui.ClearScreen()
		fmt.Println("🐳 Docker Container Dashboard")
		fmt.Println("----------------------------")

		containers, err := docker.GetContainers()
		if err != nil {
			ui.Error("Failed to fetch Docker containers: %v", err)
			pauseClean()
			return false
		}

		if len(containers) == 0 {
			ui.Success("No Docker containers found on the system.")
			pauseClean()
			return false
		}

		var options []huh.Option[string]
		options = append(options, huh.NewOption("🚪 Back to Main Menu", "back"))

		containerMap := make(map[string]docker.Container)
		for _, c := range containers {
			containerMap[c.ID] = c
			statusIndicator := "🔴"
			if c.IsRunning {
				statusIndicator = "🟢"
			}
			label := fmt.Sprintf("%s %s (%s) - Image: %s", statusIndicator, c.Name, c.Status, c.Image)
			options = append(options, huh.NewOption(label, c.ID))
		}

		var selectedID string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a container to manage").
					Description("Use Up/Down arrows to navigate, Enter to confirm. Press Esc to go back.").
					Options(options...).
					Value(&selectedID),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := form.Run(); err != nil || selectedID == "back" {
			return false
		}

		c := containerMap[selectedID]
		handleContainerOps(c)
	}
}

func handleContainerOps(c docker.Container) {
	for {
		ui.ClearScreen()
		statusIndicator := "🔴 Stopped"
		if c.IsRunning {
			statusIndicator = "🟢 Running"
		}
		fmt.Printf("🐳 Managing: %s (%s)\n", c.Name, statusIndicator)
		fmt.Printf("ID: %s | Image: %s\n", c.ID, c.Image)
		fmt.Println("-------------------------------------------")

		var action containerAction
		var options []huh.Option[containerAction]

		if c.IsRunning {
			options = append(options, huh.NewOption("⏹️  Stop Container", actionStop))
		} else {
			options = append(options, huh.NewOption("▶️  Start Container", actionStart))
		}
		options = append(options, huh.NewOption("🔄 Restart Container", actionRestart))
		options = append(options, huh.NewOption("📄 View Logs (last 50 lines)", actionLogs))
		options = append(options, huh.NewOption("🚪 Back to Dashboard", actionBack))

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[containerAction]().
					Title("Choose action").
					Description("Use Up/Down arrows to navigate, Enter to confirm. Press Esc to go back.").
					Options(options...).
					Value(&action),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := form.Run(); err != nil || action == actionBack {
			return
		}

		switch action {
		case actionStart:
			ui.ActionHeader("▶️", "Starting", c.Name)
			err := docker.StartContainer(c.ID)
			ui.ActionResult(c.Name, err, "start")
			if err == nil {
				c.IsRunning = true
				c.Status = "Up less than a second"
			}
			pauseClean()

		case actionStop:
			ui.ActionHeader("⏹️", "Stopping", c.Name)
			err := docker.StopContainer(c.ID)
			ui.ActionResult(c.Name, err, "stop")
			if err == nil {
				c.IsRunning = false
				c.Status = "Exited (0) less than a second ago"
			}
			pauseClean()

		case actionRestart:
			ui.ActionHeader("🔄", "Restarting", c.Name)
			err := docker.RestartContainer(c.ID)
			ui.ActionResult(c.Name, err, "restart")
			if err == nil {
				c.IsRunning = true
				c.Status = "Up less than a second (restarted)"
			}
			pauseClean()

		case actionLogs:
			ui.ClearScreen()
			fmt.Printf("📄 Logs for container: %s (last 50 lines)\n", c.Name)
			fmt.Println("----------------------------------------------------------------------")
			logs, err := docker.GetContainerLogs(c.ID)
			if err != nil {
				ui.Error("Failed to fetch logs: %v", err)
			} else {
				fmt.Print(logs)
			}
			fmt.Println("----------------------------------------------------------------------")
			pauseClean()
		}
	}
}

func pauseClean() {
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
