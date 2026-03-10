package updater

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// normalizeVersion
// ---------------------------------------------------------------------------

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
		{"  ", ""},
		{"vv1.0.0", "v1.0.0"},
		{"V1.0.0", "V1.0.0"},
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

// ---------------------------------------------------------------------------
// parseVersionParts
// ---------------------------------------------------------------------------

func TestParseVersionParts(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"1.2.3", []int{1, 2, 3}},
		{"0.1.0", []int{0, 1, 0}},
		{"10.20.30", []int{10, 20, 30}},
		{"1.2.3-beta", []int{1, 2, 3}},
		{"1.2.3-rc.1", []int{1, 2, 3}},
		{"1.2", []int{1, 2}},
		{"5", []int{5}},
		{"", []int{0}},
		{"abc", []int{0}},
		{"0.0.0", []int{0, 0, 0}},
		{"999.888.777", []int{999, 888, 777}},
		{"1.2.3.4", []int{1, 2, 3, 4}},
		{"0", []int{0}},
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

// ---------------------------------------------------------------------------
// isNewer
// ---------------------------------------------------------------------------

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

		// Additional edge cases
		{"single digit newer", "2", "1", true},
		{"single digit same", "1", "1", false},
		{"single digit older", "1", "2", false},
		{"latest empty parsed", "", "", false},
		{"both zeroes single", "0", "0", false},
		{"four-part newer", "1.2.3.4", "1.2.3.3", true},
		{"four-part older", "1.2.3.3", "1.2.3.4", false},
		{"zero to one", "0.0.1", "0.0.0", true},
		{"major rollover", "2.0.0", "1.99.99", true},
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

func TestIsNewer_Symmetry(t *testing.T) {
	// If isNewer(a, b) is true, then isNewer(b, a) must be false
	pairs := [][2]string{
		{"2.0.0", "1.0.0"},
		{"1.1.0", "1.0.0"},
		{"1.0.1", "1.0.0"},
		{"10.0.0", "9.99.99"},
	}

	for _, p := range pairs {
		if !isNewer(p[0], p[1]) {
			t.Errorf("expected isNewer(%q, %q) = true", p[0], p[1])
		}
		if isNewer(p[1], p[0]) {
			t.Errorf("expected isNewer(%q, %q) = false (symmetry)", p[1], p[0])
		}
	}
}

// ---------------------------------------------------------------------------
// FormatNotification
// ---------------------------------------------------------------------------


func TestFormatNotification_NoUpdate(t *testing.T) {
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

func TestFormatNotification_HasUpdate_LargeJump(t *testing.T) {
	checkerMu.Lock()
	checkerResult = &UpdateInfo{
		LatestVersion:  "v10.0.0",
		CurrentVersion: "v0.1.0",
		HasUpdate:      true,
	}
	checkerMu.Unlock()

	got := FormatNotification()
	expected := "💡 New version available: v10.0.0 (you have v0.1.0)"
	if got != expected {
		t.Errorf("FormatNotification() = %q, want %q", got, expected)
	}
}

// ---------------------------------------------------------------------------
// GetUpdateInfo
// ---------------------------------------------------------------------------

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
		ReleaseURL:     "https://github.com/mttk2004/dev/releases/tag/v2.0.0",
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

func TestGetUpdateInfo_ConcurrentAccess(t *testing.T) {
	checkerMu.Lock()
	checkerResult = &UpdateInfo{
		LatestVersion:  "v3.0.0",
		CurrentVersion: "v1.0.0",
		HasUpdate:      true,
	}
	checkerMu.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := GetUpdateInfo()
			if info == nil {
				t.Error("GetUpdateInfo() returned nil during concurrent access")
			}
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// Cache file: saveCache / loadCache
// ---------------------------------------------------------------------------

func setupTempCache(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	origEnv := os.Getenv("XDG_CACHE_HOME")
	os.Setenv("XDG_CACHE_HOME", tmpDir)
	return tmpDir, func() {
		if origEnv == "" {
			os.Unsetenv("XDG_CACHE_HOME")
		} else {
			os.Setenv("XDG_CACHE_HOME", origEnv)
		}
	}
}

func TestSaveCache_CreatesFileAndDirs(t *testing.T) {
	tmpDir, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v2.0.0",
		CurrentVersion: "v1.0.0",
		ReleaseURL:     "https://example.com/release",
		CheckedAt:      time.Now().Unix(),
		HasUpdate:      true,
	}

	err := saveCache(info)
	if err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "dev", cacheFileName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("cache file not created at %q", expectedPath)
	}

	// Verify file content is valid JSON
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}

	var loaded UpdateInfo
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("cache file contains invalid JSON: %v", err)
	}

	if loaded.LatestVersion != info.LatestVersion {
		t.Errorf("cached LatestVersion = %q, want %q", loaded.LatestVersion, info.LatestVersion)
	}
	if loaded.HasUpdate != info.HasUpdate {
		t.Errorf("cached HasUpdate = %v, want %v", loaded.HasUpdate, info.HasUpdate)
	}
}

func TestLoadCache_ValidNonExpired(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v3.0.0",
		CurrentVersion: "v1.0.0",
		ReleaseURL:     "https://example.com/release/v3",
		CheckedAt:      time.Now().Unix(), // fresh
		HasUpdate:      true,
	}

	if err := saveCache(info); err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	loaded := loadCache()
	if loaded == nil {
		t.Fatal("loadCache() returned nil for valid non-expired cache")
	}

	if loaded.LatestVersion != info.LatestVersion {
		t.Errorf("loaded LatestVersion = %q, want %q", loaded.LatestVersion, info.LatestVersion)
	}
	if loaded.ReleaseURL != info.ReleaseURL {
		t.Errorf("loaded ReleaseURL = %q, want %q", loaded.ReleaseURL, info.ReleaseURL)
	}
}

func TestLoadCache_Expired(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v3.0.0",
		CurrentVersion: "v1.0.0",
		CheckedAt:      time.Now().Add(-25 * time.Hour).Unix(), // 25 hours ago => expired
		HasUpdate:      true,
	}

	if err := saveCache(info); err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	loaded := loadCache()
	if loaded != nil {
		t.Errorf("loadCache() = %+v, want nil for expired cache", loaded)
	}
}

func TestLoadCache_JustBeforeExpiry(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v3.0.0",
		CurrentVersion: "v1.0.0",
		CheckedAt:      time.Now().Add(-23 * time.Hour).Unix(), // 23 hours ago => still valid
		HasUpdate:      true,
	}

	if err := saveCache(info); err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	loaded := loadCache()
	if loaded == nil {
		t.Error("loadCache() = nil, want non-nil for cache that is 23h old (under 24h limit)")
	}
}

func TestLoadCache_FileNotExist(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	// Don't write any cache file
	loaded := loadCache()
	if loaded != nil {
		t.Errorf("loadCache() = %+v, want nil when no cache file exists", loaded)
	}
}

func TestLoadCache_CorruptedJSON(t *testing.T) {
	tmpDir, cleanup := setupTempCache(t)
	defer cleanup()

	cacheDir := filepath.Join(tmpDir, "dev")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}

	// Write garbage data
	cachePath := filepath.Join(cacheDir, cacheFileName)
	if err := os.WriteFile(cachePath, []byte("{invalid json!!!!"), 0644); err != nil {
		t.Fatalf("failed to write corrupted cache: %v", err)
	}

	loaded := loadCache()
	if loaded != nil {
		t.Errorf("loadCache() = %+v, want nil for corrupted JSON", loaded)
	}
}

func TestLoadCache_EmptyFile(t *testing.T) {
	tmpDir, cleanup := setupTempCache(t)
	defer cleanup()

	cacheDir := filepath.Join(tmpDir, "dev")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}

	cachePath := filepath.Join(cacheDir, cacheFileName)
	if err := os.WriteFile(cachePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty cache: %v", err)
	}

	loaded := loadCache()
	if loaded != nil {
		t.Errorf("loadCache() = %+v, want nil for empty file", loaded)
	}
}

func TestLoadCache_ReEvaluatesHasUpdate(t *testing.T) {
	// The cache may have been written when the binary was an older version.
	// loadCache() should re-evaluate HasUpdate against the current binary version.
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v0.1.0", // same as current version.Version
		CurrentVersion: "v0.0.1", // stale — from when the cache was written
		CheckedAt:      time.Now().Unix(),
		HasUpdate:      true, // was true when written
	}

	if err := saveCache(info); err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	loaded := loadCache()
	if loaded == nil {
		t.Fatal("loadCache() = nil, want non-nil")
	}

	// HasUpdate should be re-evaluated: v0.1.0 is NOT newer than v0.1.0 (current)
	// so HasUpdate should now be false
	if loaded.HasUpdate {
		t.Error("loadCache() should re-evaluate HasUpdate — v0.1.0 is not newer than current v0.1.0")
	}
}

func TestSaveCache_OverwritesExisting(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info1 := &UpdateInfo{
		LatestVersion: "v1.0.0",
		CheckedAt:     time.Now().Unix(),
		HasUpdate:     false,
	}
	info2 := &UpdateInfo{
		LatestVersion: "v5.0.0",
		CheckedAt:     time.Now().Unix(),
		HasUpdate:     true,
	}

	if err := saveCache(info1); err != nil {
		t.Fatalf("saveCache(info1) error: %v", err)
	}
	if err := saveCache(info2); err != nil {
		t.Fatalf("saveCache(info2) error: %v", err)
	}

	loaded := loadCache()
	if loaded == nil {
		t.Fatal("loadCache() = nil after second save")
	}
	if loaded.LatestVersion != "v5.0.0" {
		t.Errorf("loaded.LatestVersion = %q, want %q (second save should overwrite)", loaded.LatestVersion, "v5.0.0")
	}
}

// ---------------------------------------------------------------------------
// getCacheDir / getCachePath
// ---------------------------------------------------------------------------

func TestGetCacheDir_RespectsXDG(t *testing.T) {
	tmpDir := t.TempDir()
	origEnv := os.Getenv("XDG_CACHE_HOME")
	os.Setenv("XDG_CACHE_HOME", tmpDir)
	defer func() {
		if origEnv == "" {
			os.Unsetenv("XDG_CACHE_HOME")
		} else {
			os.Setenv("XDG_CACHE_HOME", origEnv)
		}
	}()

	dir, err := getCacheDir()
	if err != nil {
		t.Fatalf("getCacheDir() error: %v", err)
	}

	expected := filepath.Join(tmpDir, "dev")
	if dir != expected {
		t.Errorf("getCacheDir() = %q, want %q", dir, expected)
	}
}

func TestGetCacheDir_FallsBackToHome(t *testing.T) {
	origEnv := os.Getenv("XDG_CACHE_HOME")
	os.Unsetenv("XDG_CACHE_HOME")
	defer func() {
		if origEnv != "" {
			os.Setenv("XDG_CACHE_HOME", origEnv)
		}
	}()

	dir, err := getCacheDir()
	if err != nil {
		t.Fatalf("getCacheDir() error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".cache", "dev")
	if dir != expected {
		t.Errorf("getCacheDir() = %q, want %q (fallback to ~/.cache/dev)", dir, expected)
	}
}

func TestGetCachePath_ContainsCacheFileName(t *testing.T) {
	path, err := getCachePath()
	if err != nil {
		t.Fatalf("getCachePath() error: %v", err)
	}

	base := filepath.Base(path)
	if base != cacheFileName {
		t.Errorf("getCachePath() basename = %q, want %q", base, cacheFileName)
	}
}

// ---------------------------------------------------------------------------
// WaitForResult
// ---------------------------------------------------------------------------

func TestWaitForResult_NilChannel(t *testing.T) {
	// When checkerDone is nil, WaitForResult should return immediately
	oldDone := checkerDone
	checkerDone = nil
	defer func() { checkerDone = oldDone }()

	start := time.Now()
	WaitForResult(5 * time.Second)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("WaitForResult with nil channel took %v, want < 100ms", elapsed)
	}
}

func TestWaitForResult_AlreadyClosed(t *testing.T) {
	oldDone := checkerDone
	defer func() { checkerDone = oldDone }()

	ch := make(chan struct{})
	close(ch)
	checkerDone = ch

	start := time.Now()
	WaitForResult(5 * time.Second)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("WaitForResult with already-closed channel took %v, want < 100ms", elapsed)
	}
}

func TestWaitForResult_TimesOut(t *testing.T) {
	oldDone := checkerDone
	defer func() { checkerDone = oldDone }()

	// Channel that never closes
	checkerDone = make(chan struct{})

	start := time.Now()
	WaitForResult(100 * time.Millisecond)
	elapsed := time.Since(start)

	if elapsed < 80*time.Millisecond {
		t.Errorf("WaitForResult returned too early: %v, want >= ~100ms", elapsed)
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("WaitForResult took too long: %v, want ~100ms", elapsed)
	}
}

func TestWaitForResult_CompletesBeforeTimeout(t *testing.T) {
	oldDone := checkerDone
	defer func() { checkerDone = oldDone }()

	ch := make(chan struct{})
	checkerDone = ch

	// Close the channel after 50ms
	go func() {
		time.Sleep(50 * time.Millisecond)
		close(ch)
	}()

	start := time.Now()
	WaitForResult(5 * time.Second)
	elapsed := time.Since(start)

	if elapsed > 500*time.Millisecond {
		t.Errorf("WaitForResult should have returned quickly after channel closed, took %v", elapsed)
	}
}

// ---------------------------------------------------------------------------
// UpdateInfo struct serialization
// ---------------------------------------------------------------------------

func TestUpdateInfo_JSONRoundTrip(t *testing.T) {
	original := UpdateInfo{
		LatestVersion:  "v2.5.0",
		CurrentVersion: "v1.3.0",
		ReleaseURL:     "https://github.com/mttk2004/dev/releases/tag/v2.5.0",
		CheckedAt:      time.Now().Unix(),
		HasUpdate:      true,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded UpdateInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.LatestVersion != original.LatestVersion {
		t.Errorf("LatestVersion = %q, want %q", decoded.LatestVersion, original.LatestVersion)
	}
	if decoded.CurrentVersion != original.CurrentVersion {
		t.Errorf("CurrentVersion = %q, want %q", decoded.CurrentVersion, original.CurrentVersion)
	}
	if decoded.ReleaseURL != original.ReleaseURL {
		t.Errorf("ReleaseURL = %q, want %q", decoded.ReleaseURL, original.ReleaseURL)
	}
	if decoded.CheckedAt != original.CheckedAt {
		t.Errorf("CheckedAt = %d, want %d", decoded.CheckedAt, original.CheckedAt)
	}
	if decoded.HasUpdate != original.HasUpdate {
		t.Errorf("HasUpdate = %v, want %v", decoded.HasUpdate, original.HasUpdate)
	}
}

// ---------------------------------------------------------------------------
// Integration-style: FormatNotification after setting various states
// ---------------------------------------------------------------------------

func TestFormatNotification_AfterCacheLoad(t *testing.T) {
	_, cleanup := setupTempCache(t)
	defer cleanup()

	info := &UpdateInfo{
		LatestVersion:  "v9.9.9",
		CurrentVersion: "v0.1.0",
		CheckedAt:      time.Now().Unix(),
		HasUpdate:      true,
	}
	if err := saveCache(info); err != nil {
		t.Fatalf("saveCache() error: %v", err)
	}

	loaded := loadCache()
	if loaded == nil {
		t.Fatal("loadCache() = nil")
	}

	// Simulate what would happen at runtime
	checkerMu.Lock()
	checkerResult = loaded
	checkerMu.Unlock()

	got := FormatNotification()
	if got == "" {
		t.Error("FormatNotification() = \"\", want non-empty after loading cache with update")
	}
}

func TestFormatNotification_NoUpdateAfterUpgrade(t *testing.T) {
	// Simulate: cache says v0.1.0 is latest, current binary is v0.1.0 => no update
	checkerMu.Lock()
	checkerResult = &UpdateInfo{
		LatestVersion:  "v0.1.0",
		CurrentVersion: "v0.1.0",
		HasUpdate:      false,
	}
	checkerMu.Unlock()

	got := FormatNotification()
	if got != "" {
		t.Errorf("FormatNotification() = %q, want empty when versions match", got)
	}
}
