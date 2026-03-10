package doctor

import (
	"fmt"

	"dev/internal/registry"
	"dev/internal/system"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// CheckResult represents the outcome of a single system check.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
	Version string
	Path    string
}

// Report aggregates multiple check results.
type Report struct {
	Results []CheckResult
}

// RunChecks performs a series of system diagnostics and returns a Report.
func RunChecks() Report {
	var results []CheckResult

	for _, pkg := range registry.Packages {
		passed := pkg.IsInstalled()
		var msg string
		var version string
		var path string

		if passed {
			if pkg.ID == "node" {
				if system.CommandExists("node") {
					msg = "fnm and node are installed and in PATH"
					version = system.GetCommandVersion("node")
					path = system.GetCommandPath("node")
				} else {
					passed = false
					msg = "fnm is installed, but node is not in PATH. Try restarting your terminal."
					version = pkg.GetVersion()
					path = pkg.GetPath()
				}
			} else {
				msg = fmt.Sprintf("%s is installed and in PATH", pkg.CheckCmd)
				version = pkg.GetVersion()
				path = pkg.GetPath()
			}
		} else {
			if pkg.ID == "bun" && system.HasEnvVarInZshrc("BUN_INSTALL") {
				msg = "BUN_INSTALL is in .zshrc, but bun is not in PATH. Try restarting your terminal."
			} else {
				msg = fmt.Sprintf("%s is missing or not in PATH. Run 'dev install %s'.", pkg.CheckCmd, pkg.ID)
			}
			version = "-"
			path = "-"
		}

		results = append(results, CheckResult{
			Name:    pkg.DisplayName,
			Passed:  passed,
			Message: msg,
			Version: version,
			Path:    path,
		})
	}

	return Report{Results: results}
}

// Print outputs the report to the console in a readable format using a lipgloss table.
func (r Report) Print() {
	fmt.Println("\n🩺 System Diagnostics Report:\n")

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers("STATUS", "PACKAGE", "VERSION", "WHICH", "DETAILS").
		Wrap(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			// Set maximum widths for columns to ensure they wrap properly instead of overflowing
			switch col {
			case 2: // VERSION
				style = style.Width(15)
			case 3: // WHICH
				style = style.Foreground(lipgloss.Color("240")).Width(30)
			case 4: // DETAILS
				style = style.Width(35)
			}

			if row == 0 {
				return style.Bold(true).Foreground(lipgloss.Color("63"))
			}
			return style
		})

	allPassed := true
	for _, res := range r.Results {
		status := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓ OK")
		if !res.Passed {
			allPassed = false
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ ERR")
		}

		version := res.Version
		if version == "" || version == "unknown" {
			version = "-"
		}
		path := res.Path
		if path == "" {
			path = "-"
		}

		t.Row(status, res.Name, version, path, res.Message)
	}

	fmt.Println(t.Render())

	if allPassed {
		fmt.Println("\n✨ " + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("Your system is ready for web development!"))
	} else {
		fmt.Println("\n⚠️  " + lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render("Some issues were found. Please review the ERR items above."))
	}
	fmt.Println()
}
