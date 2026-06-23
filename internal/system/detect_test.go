package system

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Standard semver with v prefix
		{"semver with v", "v1.2.3", "v1.2.3"},
		{"semver without v", "1.2.3", "1.2.3"},
		{"two-part version", "1.2", "1.2"},

		// Real-world version strings from common tools
		{"node version", "v22.11.0", "v22.11.0"},
		{"go version output", "go version go1.26.1 linux/amd64", "1.26.1"},
		{"python version", "Python 3.12.7", "3.12.7"},
		{"docker version", "Docker version 27.3.1, build ce12230", "27.3.1"},
		{"php version", "PHP 8.3.13 (cli) (built: Nov 2024)", "8.3.13"},
		{"java version", "openjdk version \"21.0.5\" 2024-10-15", "21.0.5"},
		{"nginx version", "nginx version: nginx/1.26.2", "1.26.2"},
		{"redis version", "redis-cli 7.2.6", "7.2.6"},
		{"psql version", "psql (PostgreSQL) 16.4", "16.4"},
		{"composer version", "Composer version 2.7.9 2024-09-25", "2.7.9"},
		{"bun version", "1.1.33", "1.1.33"},
		{"fnm version", "fnm 1.37.2", "1.37.2"},
		{"maven version", "Apache Maven 3.9.9 (8e8579a9e76f7d015ee5ec7bfcdc97d260186937)", "3.9.9"},
		{"mariadb version", "mariadb  Ver 15.1 Distrib 11.5.2-MariaDB", "15.1"},

		// Edge cases
		{"version in parens", "(v2.0.0)", "v2.0.0"},
		{"version after equals", "version=3.4.5", "3.4.5"},
		{"multi-digit major", "100.200.300", "100.200.300"},
		{"four-part version", "1.2.3.4", "1.2.3.4"},

		// Fallback behavior — no version pattern found
		{"no version short", "hello", "hello"},
		{"no version exact 20", "12345678901234567890", "12345678901234567890"},
		{"no version long truncated", "this is a very long string without any version number at all", "this is a very lo..."},
		{"empty string", "", ""},
		{"only text", "some random output", "some random output"},

		// Multiple versions — should pick the first
		{"multiple versions", "v1.0.0 upgraded to v2.0.0", "v1.0.0"},

		// Version with pre-release suffix
		{"version with beta", "v1.2.3-beta.1 extra", "v1.2.3"},
		{"version with rc", "tool 3.0.0-rc1 (release)", "3.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseVersion(tt.input)
			if got != tt.expected {
				t.Errorf("ParseVersion(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseVersion_Idempotent(t *testing.T) {
	// Parsing an already-clean version should return it unchanged
	versions := []string{"1.2.3", "v1.0.0", "22.11.0"}
	for _, v := range versions {
		first := ParseVersion(v)
		second := ParseVersion(first)
		if first != second {
			t.Errorf("ParseVersion is not idempotent: ParseVersion(%q)=%q, ParseVersion(%q)=%q",
				v, first, first, second)
		}
	}
}

func TestGetPackageManager(t *testing.T) {
	// We can only verify that the function returns a valid value.
	// On the test machine, at least one of yay/paru/pacman should exist
	// (or the function falls through to "pacman" as default).
	pm := GetPackageManager()

	validPMs := map[string]bool{
		"yay":    true,
		"paru":   true,
		"pacman": true,
	}

	if !validPMs[pm] {
		t.Errorf("GetPackageManager() = %q, want one of yay/paru/pacman", pm)
	}
}

func TestGetPackageManager_Priority(t *testing.T) {
	// Verify the function prefers yay over paru over pacman.
	// If yay exists, it must be returned.
	pm := GetPackageManager()

	if CommandExists("yay") && pm != "yay" {
		t.Errorf("GetPackageManager() = %q, want %q (yay is installed and should have priority)", pm, "yay")
	} else if !CommandExists("yay") && CommandExists("paru") && pm != "paru" {
		t.Errorf("GetPackageManager() = %q, want %q (paru is installed, yay is not)", pm, "paru")
	} else if !CommandExists("yay") && !CommandExists("paru") && pm != "pacman" {
		t.Errorf("GetPackageManager() = %q, want %q (no AUR helper installed)", pm, "pacman")
	}
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name   string
		cmd    string
		exists bool
	}{
		// These should always exist on any Linux system
		{"sh exists", "bash", true},
		{"ls exists", "ls", true},
		// These should never exist
		{"nonexistent command", "this_command_definitely_does_not_exist_xyz_abc_123", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommandExists(tt.cmd)
			if got != tt.exists {
				t.Errorf("CommandExists(%q) = %v, want %v", tt.cmd, got, tt.exists)
			}
		})
	}
}

func TestGetCommandPath(t *testing.T) {
	// A command that exists should return a non-empty absolute path
	path := GetCommandPath("bash")
	if path == "" {
		t.Error("GetCommandPath(\"bash\") = \"\", want non-empty path")
	}
	// Path should be absolute (starts with /)
	if len(path) > 0 && path[0] != '/' {
		t.Errorf("GetCommandPath(\"bash\") = %q, want absolute path starting with /", path)
	}

	// A command that doesn't exist should return empty string
	path = GetCommandPath("this_command_definitely_does_not_exist_xyz_abc_123")
	if path != "" {
		t.Errorf("GetCommandPath(nonexistent) = %q, want \"\"", path)
	}
}

func TestGetCommandVersion(t *testing.T) {
	// We can't predict exact output, but we can verify it doesn't panic
	// and returns something for known commands

	// "bash" should return something (not panic)
	ver := GetCommandVersion("bash")
	if ver == "" {
		t.Error("GetCommandVersion(\"bash\") = \"\", want non-empty string")
	}

	// A nonexistent command should return "unknown"
	ver = GetCommandVersion("this_command_definitely_does_not_exist_xyz_abc_123")
	if ver != "unknown" {
		t.Errorf("GetCommandVersion(nonexistent) = %q, want %q", ver, "unknown")
	}
}
