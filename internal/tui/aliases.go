package tui

import (
	"bufio"
	"fmt"
	"os"

	"dev/internal/aliases"
	"dev/internal/ui"

	"github.com/charmbracelet/huh"
)

type aliasesAction string

const (
	aliasesActionAdd    aliasesAction = "add"
	aliasesActionRemove aliasesAction = "remove"
	aliasesActionBack   aliasesAction = "back"
)

// RunAliasesAction runs the shell alias / zshrc manager.
func RunAliasesAction() bool {
	for {
		ui.ClearScreen()
		fmt.Println("🔗 Shell Alias Manager (~/.zshrc)")
		fmt.Println("---------------------------------")

		list, err := aliases.GetAliases()
		if err != nil {
			ui.Error("Failed to read aliases from ~/.zshrc: %v", err)
		} else {
			if len(list) == 0 {
				ui.Subtle("No custom aliases found in ~/.zshrc.")
			} else {
				fmt.Println("Configured Zsh aliases:")
				for _, a := range list {
					fmt.Printf("  • alias \033[36m%s\033[0m='%s'\n", a.Name, a.Value)
				}
			}
		}
		fmt.Println()

		var action aliasesAction
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[aliasesAction]().
					Title("Select an action").
					Options(
						huh.NewOption("➕ Add a new alias", aliasesActionAdd),
						huh.NewOption("➖ Remove an alias", aliasesActionRemove),
						huh.NewOption("🚪 Back to Main Menu", aliasesActionBack),
					).
					Value(&action),
			),
		).WithTheme(huh.ThemeCatppuccin())

		err = form.Run()
		if err != nil || action == aliasesActionBack {
			return false
		}

		switch action {
		case aliasesActionAdd:
			handleAddAlias()
		case aliasesActionRemove:
			handleRemoveAlias(list)
		}
	}
}

func handleAddAlias() {
	var name string
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter alias shortcut name").
				Description("e.g. gco, art, dps").
				Value(&name),
			huh.NewInput().
				Title("Enter alias command value").
				Description("e.g. git checkout, php artisan, docker ps").
				Value(&value),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || name == "" || value == "" {
		return
	}

	ui.ActionHeader("➕", "Adding alias mapping for", name)
	err := aliases.AddAlias(name, value)
	ui.ActionResult(name, err, "add alias")
	fmt.Println("\n💡 Note: Restart your terminal or run 'source ~/.zshrc' to apply new aliases.")
	pauseAliases()
}

func handleRemoveAlias(list []aliases.Alias) {
	if len(list) == 0 {
		ui.Warning("No custom aliases found to remove.")
		pauseAliases()
		return
	}

	var options []huh.Option[string]
	for _, a := range list {
		label := fmt.Sprintf("%s -> %s", a.Name, a.Value)
		options = append(options, huh.NewOption(label, a.Name))
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select alias to remove").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil || selected == "" {
		return
	}

	ui.ActionHeader("➖", "Removing alias", selected)
	err := aliases.RemoveAlias(selected)
	ui.ActionResult(selected, err, "remove alias")
	fmt.Println("\n💡 Note: Restart your terminal or run 'source ~/.zshrc' to apply changes.")
	pauseAliases()
}

func pauseAliases() {
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
