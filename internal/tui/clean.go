package tui

import (
	"fmt"
	"os"
	"os/exec"

	"dev/internal/clean"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

// RunCleanAction launches the interactive cache cleaning dashboard.
func RunCleanAction() bool {
	ui.ClearScreen()
	fmt.Println("🧼 Clean Dev Caches")
	fmt.Println("-------------------")
	fmt.Println("Scanning for system developer caches, please wait...")

	items := clean.GetCacheItems()

	var options []huh.Option[string]
	itemMap := make(map[string]clean.CacheItem)

	for _, item := range items {
		itemMap[item.ID] = item
		label := fmt.Sprintf("%s (%s) - %s", item.Name, clean.FormatSize(item.Size), item.Description)
		options = append(options, huh.NewOption(label, item.ID))
	}

	// Add special actions
	options = append(options, huh.NewOption("🐳 Docker Prune (Clean stopped containers, dangling images)", "docker_prune"))
	options = append(options, huh.NewOption("📦 Scan local node_modules (Search recursively in current dir)", "scan_node"))

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select items to clean").
				Description("Space to select, Enter to confirm. Press Esc to go back.").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil || len(selected) == 0 {
		return false
	}

	// Check if "scan_node" was selected
	hasScanNode := false
	var finalCleanList []string
	for _, sel := range selected {
		if sel == "scan_node" {
			hasScanNode = true
		} else {
			finalCleanList = append(finalCleanList, sel)
		}
	}

	// 1. If scan_node was selected, perform the recursive local scan
	if hasScanNode {
		ui.ClearScreen()
		fmt.Println("🔍 Scanning for local 'node_modules' (depth <= 4)...")
		cwd, err := os.Getwd()
		if err == nil {
			nodeItems, err := clean.ScanNodeModules(cwd)
			if err != nil {
				ui.Error("Failed to scan node_modules: %v", err)
				pause()
			} else if len(nodeItems) == 0 {
				ui.Success("No node_modules folders found in current directory tree.")
				pause()
			} else {
				var nodeOptions []huh.Option[string]
				nodeMap := make(map[string]clean.CacheItem)
				for _, item := range nodeItems {
					nodeMap[item.ID] = item
					label := fmt.Sprintf("%s - %s", item.Name, clean.FormatSize(item.Size))
					nodeOptions = append(nodeOptions, huh.NewOption(label, item.ID))
				}

				var nodeSelected []string
				nodeForm := huh.NewForm(
					huh.NewGroup(
						huh.NewMultiSelect[string]().
							Title("Select node_modules to delete").
							Description("Space to select, Enter to confirm. Press Esc to cancel.").
							Options(nodeOptions...).
							Value(&nodeSelected),
					),
				).WithTheme(huh.ThemeCatppuccin())

				if err := nodeForm.Run(); err == nil && len(nodeSelected) > 0 {
					for _, sel := range nodeSelected {
						item := nodeMap[sel]
						ui.ActionHeader("🧹", "Deleting", item.Description)
						err := os.RemoveAll(item.Path)
						ui.ActionResult(item.Name, err, "delete")
					}
					pause()
				}
			}
		}
	}

	// 2. Perform other clean actions
	if len(finalCleanList) > 0 {
		ui.ClearScreen()
		for _, sel := range finalCleanList {
			if sel == "docker_prune" {
				ui.ActionHeader("🐳", "Pruning", "Docker resources")
				cmd := exec.Command("docker", "system", "prune", "-a", "-f", "--volumes")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				ui.ActionResult("Docker system prune", err, "prune")
			} else {
				item := itemMap[sel]
				ui.ActionHeader("🧹", "Cleaning", item.Name)
				var err error
				if item.ID == "pacman" {
					// Clear ALL packages from cache using pacman -Scc
					cmd := exec.Command("sudo", "pacman", "-Scc", "--noconfirm")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()

					// Clean up leftover corrupt temp download-* files
					exec.Command("sudo", "rm", "-f", "/var/cache/pacman/pkg/download-*").Run()
				} else {
					err = os.RemoveAll(item.Path)
				}
				ui.ActionResult(item.Name, err, "clean")
			}
		}
		pause()
	}

	return false
}
