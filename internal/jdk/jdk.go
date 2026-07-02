package jdk

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// JDKStatus represents an installed JDK environment.
type JDKStatus struct {
	EnvName   string
	IsDefault bool
}

// JDKInfo represents a defined JDK version available for installation.
type JDKInfo struct {
	Version      int
	PackageName  string
	EnvName      string
	IsLTS        bool
	Description  string
}

// RecommendedJDKs defines the list of JDKs that can be installed/managed.
var RecommendedJDKs = []JDKInfo{
	{Version: 8, PackageName: "jdk8-openjdk", EnvName: "java-8-openjdk", IsLTS: true, Description: "Java 8 - Legacy LTS"},
	{Version: 11, PackageName: "jdk11-openjdk", EnvName: "java-11-openjdk", IsLTS: true, Description: "Java 11 - Previous LTS"},
	{Version: 17, PackageName: "jdk17-openjdk", EnvName: "java-17-openjdk", IsLTS: true, Description: "Java 17 - Very popular LTS"},
	{Version: 21, PackageName: "jdk21-openjdk", EnvName: "java-21-openjdk", IsLTS: true, Description: "Java 21 - Latest LTS"},
	{Version: 25, PackageName: "jdk25-openjdk", EnvName: "java-25-openjdk", IsLTS: true, Description: "Java 25 - Upcoming/New LTS"},
	{Version: 26, PackageName: "jdk-openjdk", EnvName: "java-26-openjdk", IsLTS: false, Description: "Java 26 - Current stable release"},
}

// GetInstalledJDKs runs archlinux-java status and parses the output.
func GetInstalledJDKs() ([]JDKStatus, error) {
	cmd := exec.Command("archlinux-java", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		if _, errPath := exec.LookPath("archlinux-java"); errPath != nil {
			return nil, errPath
		}
		// It might return non-zero when no Java environments are configured
	}

	var jdks []JDKStatus
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Available Java environments:") || strings.HasPrefix(trimmed, "No Java environments installed") {
			continue
		}

		isDefault := false
		name := trimmed
		if strings.Contains(trimmed, "(default)") {
			isDefault = true
			name = strings.TrimSpace(strings.Replace(trimmed, "(default)", "", 1))
		}

		if name != "" {
			jdks = append(jdks, JDKStatus{
				EnvName:   name,
				IsDefault: isDefault,
			})
		}
	}

	return jdks, nil
}

// GetPackageOwningJDK returns the pacman package name that owns the JDK environment.
func GetPackageOwningJDK(envName string) (string, error) {
	dirPath := "/usr/lib/jvm/" + envName
	cmd := exec.Command("pacman", "-Qo", dirPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return getFallbackPackageName(envName), nil
	}

	outStr := string(output)
	if strings.Contains(outStr, "is owned by") {
		parts := strings.Split(outStr, "is owned by")
		if len(parts) > 1 {
			pkgInfo := strings.TrimSpace(parts[1])
			pkgParts := strings.Fields(pkgInfo)
			if len(pkgParts) > 0 {
				return pkgParts[0], nil
			}
		}
	}

	return getFallbackPackageName(envName), nil
}

func getFallbackPackageName(envName string) string {
	switch envName {
	case "java-8-openjdk":
		return "jdk8-openjdk"
	case "java-11-openjdk":
		return "jdk11-openjdk"
	case "java-17-openjdk":
		return "jdk17-openjdk"
	case "java-21-openjdk":
		return "jdk21-openjdk"
	case "java-25-openjdk":
		return "jdk25-openjdk"
	case "java-26-openjdk":
		return "jdk-openjdk"
	}

	trimmed := strings.Replace(envName, "java-", "jdk", 1)
	return trimmed
}

// DetectProjectJavaVersion scans directory files for Java version requirements.
// Returns the detected major version and the filename where it was found, or 0, "".
func DetectProjectJavaVersion(dir string) (int, string) {
	// Look for .java-version
	javaVerFile := dir + "/.java-version"
	if data, err := os.ReadFile(javaVerFile); err == nil {
		verStr := strings.TrimSpace(string(data))
		if ver := parseVersionNumber(verStr); ver > 0 {
			return ver, ".java-version"
		}
	}

	// Look for pom.xml
	pomFile := dir + "/pom.xml"
	if data, err := os.ReadFile(pomFile); err == nil {
		content := string(data)
		// Try <java.version>
		reJavaVer := regexp.MustCompile(`<java\.version>(?:1\.)?(\d+)</java\.version>`)
		if m := reJavaVer.FindStringSubmatch(content); len(m) > 1 {
			if ver, _ := strconv.Atoi(m[1]); ver > 0 {
				return ver, "pom.xml (<java.version>)"
			}
		}
		// Try <maven.compiler.source>
		reSource := regexp.MustCompile(`<maven\.compiler\.source>(?:1\.)?(\d+)</maven\.compiler\.source>`)
		if m := reSource.FindStringSubmatch(content); len(m) > 1 {
			if ver, _ := strconv.Atoi(m[1]); ver > 0 {
				return ver, "pom.xml (<maven.compiler.source>)"
			}
		}
		// Try <maven.compiler.target>
		reTarget := regexp.MustCompile(`<maven\.compiler\.target>(?:1\.)?(\d+)</maven\.compiler\.target>`)
		if m := reTarget.FindStringSubmatch(content); len(m) > 1 {
			if ver, _ := strconv.Atoi(m[1]); ver > 0 {
				return ver, "pom.xml (<maven.compiler.target>)"
			}
		}
		// Try <release>
		reRelease := regexp.MustCompile(`<release>(\d+)</release>`)
		if m := reRelease.FindStringSubmatch(content); len(m) > 1 {
			if ver, _ := strconv.Atoi(m[1]); ver > 0 {
				return ver, "pom.xml (<release>)"
			}
		}
	}

	// Look for build.gradle
	gradleFile := dir + "/build.gradle"
	if data, err := os.ReadFile(gradleFile); err == nil {
		if ver, detail := parseGradleContent(string(data)); ver > 0 {
			return ver, "build.gradle (" + detail + ")"
		}
	}

	// Look for build.gradle.kts
	gradleKtsFile := dir + "/build.gradle.kts"
	if data, err := os.ReadFile(gradleKtsFile); err == nil {
		if ver, detail := parseGradleContent(string(data)); ver > 0 {
			return ver, "build.gradle.kts (" + detail + ")"
		}
	}

	// Look for system.properties
	sysPropsFile := dir + "/system.properties"
	if data, err := os.ReadFile(sysPropsFile); err == nil {
		content := string(data)
		reSys := regexp.MustCompile(`java\.runtime\.version\s*=\s*(?:1\.)?(\d+)`)
		if m := reSys.FindStringSubmatch(content); len(m) > 1 {
			if ver, _ := strconv.Atoi(m[1]); ver > 0 {
				return ver, "system.properties"
			}
		}
	}

	return 0, ""
}

func parseVersionNumber(s string) int {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`(?:1\.)?(\d+)`)
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		ver, _ := strconv.Atoi(m[1])
		return ver
	}
	return 0
}

func parseGradleContent(content string) (int, string) {
	reToolchain := regexp.MustCompile(`languageVersion\s*=\s*JavaLanguageVersion\.of\((\d+)\)`)
	if m := reToolchain.FindStringSubmatch(content); len(m) > 1 {
		if ver, _ := strconv.Atoi(m[1]); ver > 0 {
			return ver, "toolchain"
		}
	}

	reSource := regexp.MustCompile(`sourceCompatibility\s*=\s*(?:JavaVersion\.VERSION_)?['"]?(?:1\.)?(\d+)['"]?`)
	if m := reSource.FindStringSubmatch(content); len(m) > 1 {
		if ver, _ := strconv.Atoi(m[1]); ver > 0 {
			return ver, "sourceCompatibility"
		}
	}

	reTarget := regexp.MustCompile(`targetCompatibility\s*=\s*(?:JavaVersion\.VERSION_)?['"]?(?:1\.)?(\d+)['"]?`)
	if m := reTarget.FindStringSubmatch(content); len(m) > 1 {
		if ver, _ := strconv.Atoi(m[1]); ver > 0 {
			return ver, "targetCompatibility"
		}
	}

	return 0, ""
}
