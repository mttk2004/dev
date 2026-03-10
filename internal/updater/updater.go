package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"dev/internal/version"
)

const (
	// GitHubRepo is the owner/repo for checking releases.
	// Update this to match your actual GitHub repository.
	GitHubRepo = "kietdev/dev"

	// githubAPIURL is the endpoint for fetching the latest release.
	githubAPIURL = "https://api.github.com/repos/%s/releases/latest"

	// cacheDuration defines how long the cached release info is valid.
	cacheDuration = 24 * time.Hour

	// httpTimeout is the timeout for the GitHub API request.
	httpTimeout = 5 * time.Second

	// cacheFileName is the name of the cache file stored under ~/.cache/dev/
	cacheFileName = "latest-release.json"
)

// UpdateInfo holds the result of an update check.
type UpdateInfo struct {
	LatestVersion  string `json:"latest_version"`
	CurrentVersion string `json:"current_version"`
	ReleaseURL     string `json:"release_url"`
	CheckedAt      int64  `json:"checked_at"`
	HasUpdate      bool   `json:"has_update"`
}

// checker is the singleton that holds the async result.
var (
	checkerOnce   sync.Once
	checkerResult *UpdateInfo
	checkerMu     sync.RWMutex
	checkerDone   chan struct{}
)

// StartAsyncCheck kicks off a background goroutine that checks for a newer
// release on GitHub. It only runs the check once per process lifetime.
// Call GetUpdateInfo() later to retrieve the result (non-blocking if not ready).
func StartAsyncCheck() {
	checkerDone = make(chan struct{})
	checkerOnce.Do(func() {
		go func() {
			defer close(checkerDone)
			info := checkForUpdate()
			if info != nil {
				checkerMu.Lock()
				checkerResult = info
				checkerMu.Unlock()
			}
		}()
	})
}

// WaitForResult blocks until the async check is complete (or was never started).
// Use this right before rendering the doctor report so the result is available.
// It has a maximum wait time to avoid blocking the UI indefinitely.
func WaitForResult(maxWait time.Duration) {
	if checkerDone == nil {
		return
	}
	select {
	case <-checkerDone:
	case <-time.After(maxWait):
	}
}

// GetUpdateInfo returns the update info if available, or nil if the check
// hasn't completed yet or no update was found.
func GetUpdateInfo() *UpdateInfo {
	checkerMu.RLock()
	defer checkerMu.RUnlock()
	return checkerResult
}

// FormatNotification returns a styled notification string if an update is
// available, or an empty string if not.
func FormatNotification() string {
	info := GetUpdateInfo()
	if info == nil || !info.HasUpdate {
		return ""
	}
	return fmt.Sprintf("💡 New version available: %s (you have %s)",
		info.LatestVersion, info.CurrentVersion)
}

// checkForUpdate performs the actual update check, using the cache if valid.
func checkForUpdate() *UpdateInfo {
	// Try loading from cache first
	if cached := loadCache(); cached != nil {
		return cached
	}

	// Fetch from GitHub
	latestVersion, releaseURL, err := fetchLatestRelease()
	if err != nil {
		// Silently fail — this is a background check, not critical
		return nil
	}

	currentVersion := normalizeVersion(version.Version)
	latest := normalizeVersion(latestVersion)

	info := &UpdateInfo{
		LatestVersion:  "v" + latest,
		CurrentVersion: "v" + currentVersion,
		ReleaseURL:     releaseURL,
		CheckedAt:      time.Now().Unix(),
		HasUpdate:      isNewer(latest, currentVersion),
	}

	// Save to cache regardless of whether there's an update
	_ = saveCache(info)

	return info
}

// githubRelease represents the relevant fields from the GitHub releases API.
type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// fetchLatestRelease calls the GitHub API to get the latest release tag.
func fetchLatestRelease() (string, string, error) {
	url := fmt.Sprintf(githubAPIURL, GitHubRepo)

	client := &http.Client{Timeout: httpTimeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "dev-cli-update-checker")

	// Respect GITHUB_TOKEN if available to avoid rate limiting
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", fmt.Errorf("failed to parse release JSON: %w", err)
	}

	if release.TagName == "" {
		return "", "", fmt.Errorf("no tag_name found in release response")
	}

	return release.TagName, release.HTMLURL, nil
}

// getCacheDir returns the path to the cache directory (~/.cache/dev/).
func getCacheDir() (string, error) {
	// Prefer XDG_CACHE_HOME if set
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	dir := filepath.Join(cacheHome, "dev")
	return dir, nil
}

// getCachePath returns the full path to the cache file.
func getCachePath() (string, error) {
	dir, err := getCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, cacheFileName), nil
}

// loadCache tries to load a valid (non-expired) cached update info from disk.
func loadCache() *UpdateInfo {
	cachePath, err := getCachePath()
	if err != nil {
		return nil
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil
	}

	var info UpdateInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil
	}

	// Check if cache is still valid
	checkedAt := time.Unix(info.CheckedAt, 0)
	if time.Since(checkedAt) > cacheDuration {
		return nil // Cache expired
	}

	// Re-evaluate HasUpdate against current version in case the binary was
	// upgraded since the cache was written
	currentVersion := normalizeVersion(version.Version)
	info.CurrentVersion = "v" + currentVersion
	info.HasUpdate = isNewer(normalizeVersion(info.LatestVersion), currentVersion)

	return &info
}

// saveCache writes the update info to the cache file.
func saveCache(info *UpdateInfo) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// normalizeVersion strips a leading "v" and trims whitespace from a version string.
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	return v
}

// isNewer returns true if latestVersion is strictly newer than currentVersion.
// Both inputs should be normalized (no "v" prefix).
// Compares semver-like strings: major.minor.patch
func isNewer(latest, current string) bool {
	latestParts := parseVersionParts(latest)
	currentParts := parseVersionParts(current)

	// Compare each part
	maxLen := len(latestParts)
	if len(currentParts) > maxLen {
		maxLen = len(currentParts)
	}

	for i := 0; i < maxLen; i++ {
		var l, c int
		if i < len(latestParts) {
			l = latestParts[i]
		}
		if i < len(currentParts) {
			c = currentParts[i]
		}
		if l > c {
			return true
		}
		if l < c {
			return false
		}
	}

	return false // They are equal
}

// parseVersionParts splits a version string like "1.2.3" into []int{1, 2, 3}.
func parseVersionParts(v string) []int {
	// Extract just the numeric version part (e.g., from "1.2.3-beta" get "1.2.3")
	re := regexp.MustCompile(`^(\d+(?:\.\d+)*)`)
	match := re.FindString(v)
	if match == "" {
		return []int{0}
	}

	parts := strings.Split(match, ".")
	result := make([]int, len(parts))
	for i, p := range parts {
		var n int
		fmt.Sscanf(p, "%d", &n)
		result[i] = n
	}
	return result
}
