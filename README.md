# 🚀 dev - Web Development Environment Manager

`dev` is a professional, blazing-fast Command Line Interface (CLI) tool written in Go, specifically designed to automate the tedious tasks of setting up and managing web development environments on **Arch Linux** (with **Zsh**).

Instead of manually running `pacman`, writing `curl` scripts, or editing your `~/.zshrc` file to fix `$PATH` issues, `dev` handles everything for you behind a beautiful interactive Terminal User Interface (TUI).

## ✨ Features

- **Interactive TUI:** Select packages to install using a beautiful, keyboard-driven checklist.
- **Automated Configuration:** Automatically injects necessary environment variables (like `BUN_INSTALL` and `fnm env`) into your `~/.zshrc`.
- **Node.js Version Management:** Installs and configures Node.js using `fnm` (Fast Node Manager) by default, allowing you to easily switch between Node versions.
- **System Diagnostics:** Includes a built-in `doctor` command to scan your system for missing dependencies and `$PATH` misconfigurations.
- **Bulk Actions:** Update or clean all your development tools with a single command.

## 📦 Supported Packages

- `node` (Installs `fnm` and the latest Node.js LTS)
- `bun` (Installs via official script)
- `composer` (PHP package manager)
- `jdk` (Java Development Kit - OpenJDK)
- `go` (Go programming language)

## 🛠️ Prerequisites

- **OS:** Arch Linux (relies on `pacman` for native packages).
- **Shell:** Zsh (configuration logic targets `~/.zshrc`).
- **Go:** Version 1.20+ (to build from source).

## 🚀 Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/yourusername/dev.git
cd dev
go mod tidy
go build -o dev
```

Move the executable to your path (optional):
```bash
sudo mv dev /usr/local/bin/
```

## 📖 Usage

### 1. Install Packages
You can launch the interactive TUI to select multiple packages:
```bash
dev install
```
Or install a specific package directly:
```bash
dev install node
dev install bun
```

### 2. Check System Health
Run diagnostics to ensure all your tools are installed and correctly configured in your `$PATH`:
```bash
dev doctor
```

### 3. Update Packages
Update a specific tool, or run without arguments to update your entire system (`pacman -Syu`) and all standalone tools (like Bun):
```bash
dev update
dev update composer
```

### 4. Clean System Cache
Uninstall a specific tool, or run without arguments to clear the `pacman` package cache:
```bash
dev clean
dev clean jdk
```

## 📂 Project Structure

This project follows the Standard Go Project Layout:
- `cmd/`: Cobra CLI command definitions (`install`, `doctor`, `update`, `clean`).
- `internal/pkgmanager/`: Tool-specific installation and update logic.
- `internal/system/`: OS-level utilities for checking commands and modifying `~/.zshrc`.
- `internal/doctor/`: System diagnostic and reporting logic.
- `internal/tui/`: Bubble Tea and Lipgloss UI components (Spinners, Selectors).
- `config/`: Viper configuration management (defaulting to `~/.dev.yaml`).

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the issues page.