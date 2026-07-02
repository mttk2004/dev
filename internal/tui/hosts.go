package tui

import (
	"bufio"
	"fmt"
	"os"

	"dev/internal/hosts"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

type hostsAction string

const (
	hostsActionAdd    hostsAction = "add"
	hostsActionRemove hostsAction = "remove"
	hostsActionBack   hostsAction = "back"
)

// RunHostsAction runs the local domain / hosts file manager.
func RunHostsAction() bool {
	for {
		ui.ClearScreen()
		fmt.Println("🌐 Local Domain Manager (/etc/hosts)")
		fmt.Println("------------------------------------")

		domains, err := hosts.GetLocalDomains()
		if err != nil {
			ui.Error("Failed to read local domains from /etc/hosts: %v", err)
		} else {
			if len(domains) == 0 {
				ui.Subtle("No custom local domains configured.")
			} else {
				fmt.Println("Configured local domains:")
				for _, d := range domains {
					fmt.Printf("  • %-20s -> %s\n", d.Domain, d.IP)
				}
			}
		}
		fmt.Println()

		var action hostsAction
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[hostsAction]().
					Title("Select an action").
					Options(
						huh.NewOption("➕ Add new local domain", hostsActionAdd),
						huh.NewOption("➖ Remove a local domain", hostsActionRemove),
						huh.NewOption("🚪 Back to Main Menu", hostsActionBack),
					).
					Value(&action),
			),
		).WithTheme(huh.ThemeCatppuccin())

		err = form.Run()
		if err != nil || action == hostsActionBack {
			return false
		}

		switch action {
		case hostsActionAdd:
			handleAddDomain()
		case hostsActionRemove:
			handleRemoveDomain(domains)
		}
	}
}

func handleAddDomain() {
	var domain string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter new local domain name").
				Description("e.g. my-app.local, api.test").
				Value(&domain),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || domain == "" {
		return
	}

	ui.ActionHeader("➕", "Adding local domain mapping for", domain)
	err := hosts.AddLocalDomain(domain)
	ui.ActionResult(domain, err, "add local domain")
	pauseHosts()
}

func handleRemoveDomain(domains []hosts.HostEntry) {
	if len(domains) == 0 {
		ui.Warning("No custom local domains configured to remove.")
		pauseHosts()
		return
	}

	var options []huh.Option[string]
	for _, d := range domains {
		options = append(options, huh.NewOption(d.Domain, d.Domain))
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select local domain to remove").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || selected == "" {
		return
	}

	ui.ActionHeader("➖", "Removing local domain mapping for", selected)
	err := hosts.RemoveLocalDomain(selected)
	ui.ActionResult(selected, err, "remove local domain")
	pauseHosts()
}

func pauseHosts() {
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
