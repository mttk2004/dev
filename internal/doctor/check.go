package doctor

import (
	"fmt"

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

	// 1. Check Node.js and fnm
	if system.CommandExists("fnm") {
		if system.CommandExists("node") {
			results = append(results, CheckResult{Name: "Node.js (fnm)", Passed: true, Message: "fnm and node are installed and in PATH", Version: system.GetCommandVersion("node"), Path: system.GetCommandPath("node")})
		} else {
			results = append(results, CheckResult{Name: "Node.js (fnm)", Passed: false, Message: "fnm is installed, but node is not in PATH. Try restarting your terminal.", Version: system.GetCommandVersion("fnm"), Path: system.GetCommandPath("fnm")})
		}
	} else {
		results = append(results, CheckResult{Name: "Node.js (fnm)", Passed: false, Message: "fnm is missing or not in PATH. Run 'dev install node'."})
	}

	// 2. Check Bun
	if system.CommandExists("bun") {
		results = append(results, CheckResult{Name: "Bun", Passed: true, Message: "bun is installed and in PATH", Version: system.GetCommandVersion("bun"), Path: system.GetCommandPath("bun")})
	} else {
		if system.HasEnvVarInZshrc("BUN_INSTALL") {
			results = append(results, CheckResult{Name: "Bun", Passed: false, Message: "BUN_INSTALL is in .zshrc, but bun is not in PATH. Try restarting your terminal."})
		} else {
			results = append(results, CheckResult{Name: "Bun", Passed: false, Message: "bun is missing or not in PATH. Run 'dev install bun'."})
		}
	}

	// 3. Check Composer
	if system.CommandExists("composer") {
		results = append(results, CheckResult{Name: "Composer", Passed: true, Message: "composer is installed and in PATH", Version: system.GetCommandVersion("composer"), Path: system.GetCommandPath("composer")})
	} else {
		results = append(results, CheckResult{Name: "Composer", Passed: false, Message: "composer is missing or not in PATH. Run 'dev install composer'."})
	}

	// 4. Check JDK (Java)
	if system.CommandExists("java") {
		results = append(results, CheckResult{Name: "JDK (Java)", Passed: true, Message: "java is installed and in PATH", Version: system.GetCommandVersion("java"), Path: system.GetCommandPath("java")})
	} else {
		results = append(results, CheckResult{Name: "JDK (Java)", Passed: false, Message: "java is missing or not in PATH. Run 'dev install jdk'."})
	}

	// 5. Check Go
	if system.CommandExists("go") {
		results = append(results, CheckResult{Name: "Go", Passed: true, Message: "go is installed and in PATH", Version: system.GetCommandVersion("go"), Path: system.GetCommandPath("go")})
	} else {
		results = append(results, CheckResult{Name: "Go", Passed: false, Message: "go is missing or not in PATH. Run 'dev install go'."})
	}

	// 6. Check PHP
	if system.CommandExists("php") {
		results = append(results, CheckResult{Name: "PHP", Passed: true, Message: "php is installed and in PATH", Version: system.GetCommandVersion("php"), Path: system.GetCommandPath("php")})
	} else {
		results = append(results, CheckResult{Name: "PHP", Passed: false, Message: "php is missing or not in PATH. Run 'dev install php'."})
	}

	// 7. Check Docker
	if system.CommandExists("docker") {
		results = append(results, CheckResult{Name: "Docker", Passed: true, Message: "docker is installed and in PATH", Version: system.GetCommandVersion("docker"), Path: system.GetCommandPath("docker")})
	} else {
		results = append(results, CheckResult{Name: "Docker", Passed: false, Message: "docker is missing or not in PATH. Run 'dev install docker'."})
	}

	// 8. Check PostgreSQL
	if system.CommandExists("psql") {
		results = append(results, CheckResult{Name: "PostgreSQL", Passed: true, Message: "postgresql is installed and in PATH", Version: system.GetCommandVersion("psql"), Path: system.GetCommandPath("psql")})
	} else {
		results = append(results, CheckResult{Name: "PostgreSQL", Passed: false, Message: "postgresql is missing or not in PATH. Run 'dev install postgresql'."})
	}

	// 9. Check Redis
	if system.CommandExists("redis-cli") {
		results = append(results, CheckResult{Name: "Redis", Passed: true, Message: "redis is installed and in PATH", Version: system.GetCommandVersion("redis-cli"), Path: system.GetCommandPath("redis-cli")})
	} else {
		results = append(results, CheckResult{Name: "Redis", Passed: false, Message: "redis is missing or not in PATH. Run 'dev install redis'."})
	}

	// 10. Check Nginx
	if system.CommandExists("nginx") {
		results = append(results, CheckResult{Name: "Nginx", Passed: true, Message: "nginx is installed and in PATH", Version: system.GetCommandVersion("nginx"), Path: system.GetCommandPath("nginx")})
	} else {
		results = append(results, CheckResult{Name: "Nginx", Passed: false, Message: "nginx is missing or not in PATH. Run 'dev install nginx'."})
	}

	// 11. Check Python
	if system.CommandExists("python") {
		results = append(results, CheckResult{Name: "Python", Passed: true, Message: "python is installed and in PATH", Version: system.GetCommandVersion("python"), Path: system.GetCommandPath("python")})
	} else {
		results = append(results, CheckResult{Name: "Python", Passed: false, Message: "python is missing or not in PATH. Run 'dev install python'."})
	}

	// 12. Check Maven
	if system.CommandExists("mvn") {
		results = append(results, CheckResult{Name: "Maven", Passed: true, Message: "maven is installed and in PATH", Version: system.GetCommandVersion("mvn"), Path: system.GetCommandPath("mvn")})
	} else {
		results = append(results, CheckResult{Name: "Maven", Passed: false, Message: "maven is missing or not in PATH. Run 'dev install maven'."})
	}

	// 13. Check MariaDB
	if system.CommandExists("mariadb") {
		results = append(results, CheckResult{Name: "MariaDB", Passed: true, Message: "mariadb is installed and in PATH", Version: system.GetCommandVersion("mariadb"), Path: system.GetCommandPath("mariadb")})
	} else {
		results = append(results, CheckResult{Name: "MariaDB", Passed: false, Message: "mariadb is missing or not in PATH. Run 'dev install mariadb'."})
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
				style = style.Width(20)
			case 3: // WHICH
				style = style.Foreground(lipgloss.Color("240")).Width(25)
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
		if version == "" {
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
