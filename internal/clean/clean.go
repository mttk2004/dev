package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CacheItem represents a cleanable cache location or folder.
type CacheItem struct {
	ID          string
	Name        string
	Path        string
	Size        int64
	Description string
}

// FormatSize formats a byte size into human readable format.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGT"[exp])
}

// GetCacheItems scans well-known system caches and returns their sizes.
func GetCacheItems() []CacheItem {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	var items []CacheItem

	// 1. Pacman cache
	if size, err := getDirSize("/var/cache/pacman/pkg"); err == nil && size > 0 {
		items = append(items, CacheItem{
			ID:          "pacman",
			Name:        "Pacman Package Cache",
			Path:        "/var/cache/pacman/pkg",
			Size:        size,
			Description: "Downloaded pacman packages (.pkg.tar.zst)",
		})
	}

	// 2. AUR helper caches
	if home != "" {
		yayPath := filepath.Join(home, ".cache", "yay")
		if size, err := getDirSize(yayPath); err == nil && size > 0 {
			items = append(items, CacheItem{
				ID:          "yay",
				Name:        "Yay AUR Cache",
				Path:        yayPath,
				Size:        size,
				Description: "Build caches and packages for AUR (via yay)",
			})
		}

		paruPath := filepath.Join(home, ".cache", "paru")
		if size, err := getDirSize(paruPath); err == nil && size > 0 {
			items = append(items, CacheItem{
				ID:          "paru",
				Name:        "Paru AUR Cache",
				Path:        paruPath,
				Size:        size,
				Description: "Build caches and packages for AUR (via paru)",
			})
		}
	}

	// 3. Composer cache
	if home != "" {
		composerPath := filepath.Join(home, ".cache", "composer")
		if size, err := getDirSize(composerPath); err == nil && size > 0 {
			items = append(items, CacheItem{
				ID:          "composer",
				Name:        "Composer Cache",
				Path:        composerPath,
				Size:        size,
				Description: "Cached PHP dependency packages",
			})
		}
	}

	// 4. Gradle Cache
	if home != "" {
		gradlePath := filepath.Join(home, ".gradle", "caches")
		if size, err := getDirSize(gradlePath); err == nil && size > 0 {
			items = append(items, CacheItem{
				ID:          "gradle",
				Name:        "Gradle Cache",
				Path:        gradlePath,
				Size:        size,
				Description: "Cached Java/Kotlin libraries and build caches",
			})
		}
	}

	// 5. Maven Cache
	if home != "" {
		mavenPath := filepath.Join(home, ".m2", "repository")
		if size, err := getDirSize(mavenPath); err == nil && size > 0 {
			items = append(items, CacheItem{
				ID:          "maven",
				Name:        "Maven Repository",
				Path:        mavenPath,
				Size:        size,
				Description: "Cached Maven dependencies (~/.m2/repository)",
			})
		}
	}

	return items
}

// ScanNodeModules scans the current directory and its subdirectories up to a depth of 4 for node_modules.
func ScanNodeModules(startDir string) ([]CacheItem, error) {
	var items []CacheItem
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == "node_modules" {
			size, _ := getDirSize(path)
			rel, _ := filepath.Rel(startDir, path)
			items = append(items, CacheItem{
				ID:          "nodemodules_" + path,
				Name:        "node_modules (" + filepath.Base(filepath.Dir(path)) + ")",
				Path:        path,
				Size:        size,
				Description: "./" + rel,
			})
			return filepath.SkipDir
		}

		rel, errRel := filepath.Rel(startDir, path)
		if errRel == nil {
			depth := len(strings.Split(rel, string(filepath.Separator)))
			if depth > 4 {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return items, err
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
