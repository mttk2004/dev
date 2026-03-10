package updater

import (
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"v0.1.0", "0.1.0"},
		{"1.2.3", "1.2.3"},
		{"  v2.0.0  ", "2.0.0"},
		{"v", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeVersion(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseVersionParts(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"1.2.3", []int{1, 2, 3}},
		{"0.1.0", []int{0, 1, 0}},
		{"10.20.30", []int{10, 20, 30}},
		{"1.2.3-beta", []int{1, 2, 3}},
		{"1.2", []int{1, 2}},
		{"5", []int{5}},
		{"", []int{0}},
		{"abc", []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseVersionParts(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("parseVersionParts(%q) = %v, want %v", tt.input, got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("parseVersionParts(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name    string
		latest  string
		current string
		want    bool
	}{
		{"newer major", "2.0.0", "1.0.0", true},
		{"newer minor", "1.2.0", "1.1.0", true},
		{"newer patch", "1.0.1", "1.0.0", true},
		{"same version", "1.0.0", "1.0.0", false},
		{"older major", "1.0.0", "2.0.0", false},
		{"older minor", "1.1.0", "1.2.0", false},
		{"older patch", "1.0.0", "1.0.1", false},
		{"multi-digit newer", "1.10.0", "1.9.0", true},
		{"multi-digit equal", "10.20.30", "10.20.30", false},
		{"different length newer", "1.2.3", "1.2", true},
		{"different length older", "1.2", "1.2.3", false},
		{"different length equal", "1.2.0", "1.2", false},
		{"zero vs zero", "0.0.0", "0.0.0", false},
		{"big jump", "3.0.0", "0.1.0", true},
		{"prerelease stripped", "1.2.3", "1.2.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNewer(tt.latest, tt.current)
			if got != tt.want {
				t.Errorf("isNewer(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}

func TestFormatNotification_NoUpdate(t *testing.T) {
	// Reset global state
	checkerMu.Lock()
	checkerResult = nil
	checkerMu.Unlock()

	got := FormatNotification()
	if got != "" {
		t.Errorf("FormatNotification() with no result = %q, want empty string", got)
	}
}

func TestFormatNotification_NoNewVersion(t *testing.T) {
	checkerMu.Lock()
	checkerResult = &UpdateInfo{
		LatestVersion:  "v1.0.0",
		CurrentVersion: "v1.0.0",
		HasUpdate:      false,
	}
	checkerMu.Unlock()

	got := FormatNotification()
	if got != "" {
		t.Errorf("FormatNotification() with same version = %q, want empty string", got)
	}
}

func TestFormatNotification_HasUpdate(t *testing.T) {
	checkerMu.Lock()
	checkerResult = &UpdateInfo{
		LatestVersion:  "v1.2.0",
		CurrentVersion: "v1.0.0",
		HasUpdate:      true,
	}
	checkerMu.Unlock()

	got := FormatNotification()
	expected := "💡 New version available: v1.2.0 (you have v1.0.0)"
	if got != expected {
		t.Errorf("FormatNotification() = %q, want %q", got, expected)
	}
}

func TestGetUpdateInfo_Nil(t *testing.T) {
	checkerMu.Lock()
	checkerResult = nil
	checkerMu.Unlock()

	got := GetUpdateInfo()
	if got != nil {
		t.Errorf("GetUpdateInfo() = %v, want nil", got)
	}
}

func TestGetUpdateInfo_ReturnsResult(t *testing.T) {
	expected := &UpdateInfo{
		LatestVersion:  "v2.0.0",
		CurrentVersion: "v1.0.0",
		ReleaseURL:     "https://github.com/kietdev/dev/releases/tag/v2.0.0",
		HasUpdate:      true,
	}

	checkerMu.Lock()
	checkerResult = expected
	checkerMu.Unlock()

	got := GetUpdateInfo()
	if got == nil {
		t.Fatal("GetUpdateInfo() = nil, want non-nil")
	}
	if got.LatestVersion != expected.LatestVersion {
		t.Errorf("LatestVersion = %q, want %q", got.LatestVersion, expected.LatestVersion)
	}
	if got.CurrentVersion != expected.CurrentVersion {
		t.Errorf("CurrentVersion = %q, want %q", got.CurrentVersion, expected.CurrentVersion)
	}
	if got.ReleaseURL != expected.ReleaseURL {
		t.Errorf("ReleaseURL = %q, want %q", got.ReleaseURL, expected.ReleaseURL)
	}
	if got.HasUpdate != expected.HasUpdate {
		t.Errorf("HasUpdate = %v, want %v", got.HasUpdate, expected.HasUpdate)
	}
}
