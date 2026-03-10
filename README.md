# 🚀 dev - Web Development Environment Manager

`dev` is a blazing-fast, interactive Terminal User Interface (TUI) tool written in Go. It is designed to automate the tedious tasks of setting up and managing web development environments on **Arch Linux** (with **Zsh**).

**The Philosophy:** Zero subcommands to memorize. Forget `install`, `update`, or `doctor` arguments. Just type `dev` and let the unified interactive dashboard handle the rest.

## ✨ Features

- **One-Command Dashboard:** Run `dev` to open a centralized control panel for your entire dev machine.
- **Smart Diagnostics:** Instantly see installed tools, their versions, paths, and `$PATH` misconfigurations in a beautiful auto-sizing table.
- **Package Management:** Install, update, or uninstall dev tools using a keyboard-driven checklist. Automatically detects and uses AUR helpers (`yay` or `paru`) if available.
- **Service Manager:** Start, stop, enable, or disable background services (`systemctl`) like Docker or PostgreSQL directly from the UI.
- **Project Scaffolding:** Spin up new projects in seconds (Next.js, React, Vue, Express, Laravel, Django, Spring Boot, Go API) with automated dependency isolation.
- **Automated `$PATH`:** Automatically injects necessary environment variables into your `~/.zshrc`.

## 📦 Supported Stack

- **Runtimes/Langs:** Node.js (via `fnm`), Bun, Go, PHP, Python, Java (JDK).
- **Package Managers:** Composer, Maven.
- **Databases/Cache:** PostgreSQL, MariaDB, Redis.
- **DevOps/Servers:** Docker, Nginx.

## 🛠️ Prerequisites

- **OS:** Arch Linux.
- **Shell:** Zsh.
- **Go:** Version 1.21+ (to build from source).

## 🚀 Installation

```bash
git clone https://github.com/yourusername/dev.git
cd dev
go mod tidy
go build -o dev
sudo mv dev /usr/local/bin/
```

## 📖 Usage

Drop the CLI arguments. Just run:

```bash
dev
```

From the interactive menu, you can navigate using your arrow keys and `Enter` to:
1. **📦 Install packages** (Smartly filters out already installed tools).
2. **🔄 Update packages** (Updates Arch packages & standalone tools).
3. **🧹 Uninstall / Clean packages**
4. **🔍 Search for a package** (Queries Pacman/AUR directly).
5. **⚙️ Manage Services** (Toggle running states of local databases & servers).
6. **✨ Create New Project** (Scaffold boilerplate for 9+ different frameworks).

## 📂 Project Structure

This project follows a clean Go architecture separating the UI from system logic:
- `cmd/root.go`: The single entrypoint launching the TUI.
- `internal/tui/`: Interactive UI components built with `charmbracelet/huh`.
- `internal/pkgmanager/`: Arch-native package installation and AUR routing.
- `internal/scaffold/`: Project template generation logic.
- `internal/system/`: OS-level utilities (Zsh configs, systemctl, command detection).
- `internal/doctor/`: System diagnostic and health reporting.