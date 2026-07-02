package tui

import (
	"bufio"
	"fmt"
	"os"

	"dev/internal/jdk"
	"dev/internal/pkgmanager"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

type jdkActionType string

const (
	jdkActionSwitch    jdkActionType = "switch"
	jdkActionInstall   jdkActionType = "install"
	jdkActionUninstall jdkActionType = "uninstall"
	jdkActionBack      jdkActionType = "back"
)

// RunJDKAction manages the JDK versions menu and returns true if the dashboard
// needs to pause (already handled internally here, but return true if we performed actions).
func RunJDKAction() bool {
	for {
		ui.ClearScreen()

		installed, err := jdk.GetInstalledJDKs()
		if err != nil {
			ui.Error("Failed to get installed JDKs: %v", err)
		}

		// Print current JDK status
		fmt.Println("☕ JDK Version Manager")
		fmt.Println("--------------------")
		if len(installed) == 0 {
			ui.Warning("No JDK versions detected on this system.")
		} else {
			fmt.Println("Installed JDK versions:")
			for _, j := range installed {
				if j.IsDefault {
					fmt.Printf("  ● \033[32m%s (Active)\033[0m\n", j.EnvName)
				} else {
					fmt.Printf("  ○ %s\n", j.EnvName)
				}
			}
		}
		fmt.Println()

		var choice jdkActionType
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[jdkActionType]().
					Title("Choose an action").
					Options(
						huh.NewOption("🔄 Switch default JDK version", jdkActionSwitch),
						huh.NewOption("📥 Install new JDK version", jdkActionInstall),
						huh.NewOption("🧹 Uninstall a JDK version", jdkActionUninstall),
						huh.NewOption("🚪 Back to Main Menu", jdkActionBack),
					).
					Value(&choice),
			),
		).WithTheme(huh.ThemeCatppuccin())

		err = form.Run()
		if err != nil || choice == jdkActionBack {
			return false
		}

		switch choice {
		case jdkActionSwitch:
			handleSwitchJDK(installed)
		case jdkActionInstall:
			handleInstallJDK(installed)
		case jdkActionUninstall:
			handleUninstallJDK(installed)
		}
	}
}

func handleSwitchJDK(installed []jdk.JDKStatus) {
	if len(installed) == 0 {
		ui.Warning("No JDK versions are installed. Please install one first.")
		pause()
		return
	}
	if len(installed) == 1 {
		ui.Warning("Only one JDK (%s) is installed. There are no other versions to switch to.", installed[0].EnvName)
		pause()
		return
	}

	var options []huh.Option[string]
	options = append(options, huh.NewOption("🚪 Back to JDK Menu", "back"))
	for _, j := range installed {
		label := j.EnvName
		if j.IsDefault {
			label += " (current default)"
		}
		options = append(options, huh.NewOption(label, j.EnvName))
	}

	var target string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select default JDK version").
				Description("This will update the default java environment symlinks").
				Options(options...).
				Value(&target),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || target == "back" {
		return
	}

	ui.ActionHeader("🔄", "Setting default JDK to", target)
	err := pkgmanager.RunInteractive("sudo", "archlinux-java", "set", target)
	ui.ActionResult(target, err, "set default")
	pause()
}

func handleInstallJDK(installed []jdk.JDKStatus) {
	// Detect project Java version requirements in current directory
	cwd, err := os.Getwd()
	detectedVer := 0
	detectedFile := ""
	if err == nil {
		detectedVer, detectedFile = jdk.DetectProjectJavaVersion(cwd)
	}

	// Create lookup map of installed versions
	isInstalledMap := make(map[string]bool)
	for _, j := range installed {
		isInstalledMap[j.EnvName] = true
	}

	var options []huh.Option[string]
	var recommendedPkg string

	// Try to find if the detected version matches any recommended JDK
	for _, r := range jdk.RecommendedJDKs {
		if isInstalledMap[r.EnvName] {
			continue // Already installed
		}

		label := fmt.Sprintf("JDK %d (%s)", r.Version, r.Description)
		isRecommended := false

		if detectedVer > 0 && r.Version == detectedVer {
			label = fmt.Sprintf("⭐ JDK %d (LTS) - Recommended for your current project (detected in %s)", r.Version, detectedFile)
			isRecommended = true
		} else if detectedVer == 0 && r.Version == 21 {
			label = "⭐ JDK 21 (LTS) - Recommended (Latest LTS)"
			isRecommended = true
		}

		// Prepend recommended options to top of list if matches
		if isRecommended {
			options = append([]huh.Option[string]{huh.NewOption(label, r.PackageName)}, options...)
			if recommendedPkg == "" {
				recommendedPkg = r.PackageName
			}
		} else {
			options = append(options, huh.NewOption(label, r.PackageName))
		}
	}

	if len(options) == 0 {
		ui.Success("All recommended JDK versions are already installed!")
		pause()
		return
	}

	// Add Back option at the very beginning of options
	options = append([]huh.Option[string]{huh.NewOption("🚪 Back to JDK Menu", "back")}, options...)

	var selectedPkg string
	// Set default selection to recommended if available
	if recommendedPkg != "" {
		selectedPkg = recommendedPkg
	} else {
		selectedPkg = "back"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a JDK version to install").
				Description("Suggested versions are marked with ⭐").
				Options(options...).
				Value(&selectedPkg),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || selectedPkg == "back" {
		return
	}

	ui.ActionHeader("📥", "Installing", selectedPkg)
	err = pkgmanager.PacmanInstall(selectedPkg)
	ui.ActionResult(selectedPkg, err, "install")

	// Post install suggestion: if there was no active version or they want to set it default
	newInstalled, _ := jdk.GetInstalledJDKs()
	if len(newInstalled) == 1 {
		// Auto-set default if it's the only one installed
		newEnv := newInstalled[0].EnvName
		ui.ActionHeader("🔄", "Setting default JDK to", newEnv)
		err = pkgmanager.RunInteractive("sudo", "archlinux-java", "set", newEnv)
		ui.ActionResult(newEnv, err, "set default")
	} else if len(newInstalled) > 1 {
		// Find env name matching installed package
		var newEnv string
		for _, r := range jdk.RecommendedJDKs {
			if r.PackageName == selectedPkg {
				newEnv = r.EnvName
				break
			}
		}
		if newEnv != "" {
			var confirm bool
			confirmForm := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Would you like to set %s as the default JDK version?", newEnv)).
						Value(&confirm),
				),
			).WithTheme(huh.ThemeCatppuccin())

			if err := confirmForm.Run(); err == nil && confirm {
				ui.ActionHeader("🔄", "Setting default JDK to", newEnv)
				err = pkgmanager.RunInteractive("sudo", "archlinux-java", "set", newEnv)
				ui.ActionResult(newEnv, err, "set default")
			}
		}
	}

	pause()
}

func handleUninstallJDK(installed []jdk.JDKStatus) {
	if len(installed) == 0 {
		ui.Warning("No JDK versions are installed to uninstall.")
		pause()
		return
	}

	var options []huh.Option[string]
	var activeEnv string
	for _, j := range installed {
		label := j.EnvName
		if j.IsDefault {
			label += " (Active)"
			activeEnv = j.EnvName
		}
		options = append(options, huh.NewOption(label, j.EnvName))
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select JDK version(s) to uninstall").
				Description("Space to select, Enter to confirm. Press Esc or Enter with no selection to go back.").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || len(selected) == 0 {
		return
	}

	// Safety check: uninstalling active JDK or all JDKs
	hasActive := false
	for _, sel := range selected {
		if sel == activeEnv {
			hasActive = true
			break
		}
	}

	if hasActive || len(installed) == len(selected) {
		var confirm bool
		var title string
		if len(installed) == len(selected) {
			title = "⚠️ Bạn đang gỡ cài đặt TẤT CẢ JDK hiện có. Điều này sẽ làm lỗi các công cụ phụ thuộc như Maven/Gradle. Tiếp tục?"
		} else {
			title = "⚠️ Bạn đang gỡ phiên bản JDK mặc định (Active) đang hoạt động. Tiếp tục?"
		}

		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(title).
					Value(&confirm),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := confirmForm.Run(); err != nil || !confirm {
			return
		}
	}

	for _, env := range selected {
		pkgName, err := jdk.GetPackageOwningJDK(env)
		if err != nil {
			ui.Error("Could not resolve package for %s: %v", env, err)
			continue
		}

		ui.ActionHeader("🧹", "Uninstalling", env)
		err = pkgmanager.PacmanRemove(pkgName)
		ui.ActionResult(env, err, "uninstall")

		if err != nil {
			fmt.Println("\n💡 Gợi ý: Việc gỡ cài đặt thất bại có thể do ràng buộc phụ thuộc (Dependency Check).")
			fmt.Println("   Nếu các công cụ khác (như Maven, Gradle, IDE...) đang yêu cầu môi trường Java,")
			fmt.Println("   bạn cần cài đặt một phiên bản JDK thay thế khác trước khi gỡ cài đặt phiên bản này.")
		}
	}

	// Try running fix if there are still jdks installed
	remaining, _ := jdk.GetInstalledJDKs()
	if len(remaining) > 0 {
		hasActiveRemaining := false
		for _, j := range remaining {
			if j.IsDefault {
				hasActiveRemaining = true
				break
			}
		}
		if !hasActiveRemaining {
			// Auto set or run fix
			ui.ActionHeader("🔧", "Fixing default Java symlink", "")
			err := pkgmanager.RunInteractive("sudo", "archlinux-java", "fix")
			if err != nil {
				// Let user choose a new default manually
				fmt.Println("No default JDK is active now. Please set one.")
			}
		}
	}

	pause()
}

func pause() {
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
