package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath_CurrentDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Path inside cwd should be valid
	safePath := filepath.Join(cwd, "my-project")
	if err := ValidatePath(safePath); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", safePath, err)
	}
}

func TestValidatePath_HomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// Path inside home should be valid
	safePath := filepath.Join(home, "projects", "my-app")
	if err := ValidatePath(safePath); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", safePath, err)
	}
}

func TestValidatePath_HomeDirItself(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// Home dir itself should be valid (it's an allowed base)
	if err := ValidatePath(home); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", home, err)
	}
}

func TestValidatePath_CwdItself(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	if err := ValidatePath(cwd); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", cwd, err)
	}
}

func TestValidatePath_TraversalAboveRoot(t *testing.T) {
	// Attempts to escape to /etc should fail
	err := ValidatePath("/etc/passwd")
	if err == nil {
		t.Error("ValidatePath(/etc/passwd) = nil, want error")
	}
}

func TestValidatePath_TraversalWithDotDot(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// Build a relative path that climbs above the home directory entirely.
	// From cwd, we calculate how many ".." segments are needed to reach /,
	// then append "etc" which resolves to /etc — clearly outside home & cwd.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Count depth of cwd from root
	depth := 0
	dir := cwd
	for dir != "/" && dir != "." {
		dir = filepath.Dir(dir)
		depth++
	}

	// Build enough ".." to reach root, then go into "etc"
	traversal := ""
	for i := 0; i <= depth; i++ {
		traversal = filepath.Join(traversal, "..")
	}
	traversal = filepath.Join(traversal, "etc")

	// Sanity: the resolved path must actually be outside home
	abs, _ := filepath.Abs(traversal)
	if abs == home || len(abs) > len(home) && abs[:len(home)+1] == home+string(filepath.Separator) {
		t.Skipf("resolved path %q is still under home %q; cannot test traversal escape on this layout", abs, home)
	}

	err = ValidatePath(traversal)
	if err == nil {
		t.Errorf("ValidatePath(%q) [resolves to %q] = nil, want error", traversal, abs)
	}
}

func TestValidatePath_RootPath(t *testing.T) {
	err := ValidatePath("/")
	if err == nil {
		t.Error("ValidatePath(/) = nil, want error")
	}
}

func TestValidatePath_SystemPaths(t *testing.T) {
	dangerous := []string{
		"/tmp/evil-project",
		"/var/log/something",
		"/usr/bin/hijack",
		"/opt/sneaky",
		"/root/gotcha",
	}

	for _, p := range dangerous {
		t.Run(p, func(t *testing.T) {
			err := ValidatePath(p)
			if err == nil {
				t.Errorf("ValidatePath(%q) = nil, want error for unsafe path", p)
			}
		})
	}
}

func TestValidatePath_NestedSubdirInCwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Deeply nested path under cwd should be fine
	deep := filepath.Join(cwd, "a", "b", "c", "d", "project")
	if err := ValidatePath(deep); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", deep, err)
	}
}

func TestValidatePath_NestedSubdirInHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	deep := filepath.Join(home, "dev", "projects", "2024", "my-app")
	if err := ValidatePath(deep); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", deep, err)
	}
}

func TestValidatePath_DotDotInMiddle(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// cwd/subdir/../project resolves to cwd/project which is still safe
	sneaky := filepath.Join(cwd, "subdir", "..", "project")
	if err := ValidatePath(sneaky); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil (resolves inside cwd)", sneaky, err)
	}
}

func TestValidatePath_SimilarPrefixAttack(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// /home/username_evil should NOT match /home/username
	// This tests that we use proper separator-aware prefix matching
	fakeHome := home + "_evil/project"
	err = ValidatePath(fakeHome)
	if err == nil {
		t.Errorf("ValidatePath(%q) = nil, want error (similar prefix attack)", fakeHome)
	}
}

func TestSafeBaseDirs_ContainsCwdAndHome(t *testing.T) {
	dirs := SafeBaseDirs()

	if len(dirs) < 2 {
		t.Fatalf("SafeBaseDirs() returned %d dirs, want at least 2", len(dirs))
	}

	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()

	foundCwd := false
	foundHome := false
	for _, d := range dirs {
		if d == cwd {
			foundCwd = true
		}
		if d == home {
			foundHome = true
		}
	}

	if !foundCwd {
		t.Errorf("SafeBaseDirs() does not contain cwd %q", cwd)
	}
	if !foundHome {
		t.Errorf("SafeBaseDirs() does not contain home %q", home)
	}
}

func TestValidatePath_RelativeInCwd(t *testing.T) {
	// A simple relative path like "my-project" should resolve inside cwd and be valid
	if err := ValidatePath("my-project"); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", "my-project", err)
	}
}

func TestValidatePath_DotSlashRelative(t *testing.T) {
	if err := ValidatePath("./my-project"); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil", "./my-project", err)
	}
}

func TestValidatePath_EmptyString(t *testing.T) {
	// Empty string resolves to cwd via filepath.Abs
	if err := ValidatePath(""); err != nil {
		t.Errorf("ValidatePath(%q) = %v, want nil (resolves to cwd)", "", err)
	}
}
