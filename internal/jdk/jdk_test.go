package jdk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseVersionNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"8", 8},
		{"11", 11},
		{"17", 17},
		{"21", 21},
		{"1.8", 8},
		{"1.8.0_292", 8},
		{"17.0.1", 17},
		{"21-openjdk", 21},
		{"openjdk-21", 21},
		{"random string", 0},
	}

	for _, test := range tests {
		result := parseVersionNumber(test.input)
		if result != test.expected {
			t.Errorf("parseVersionNumber(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestParseGradleContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedVer int
		expectedDet string
	}{
		{
			"toolchain syntax",
			`java {
				toolchain {
					languageVersion = JavaLanguageVersion.of(17)
				}
			}`,
			17,
			"toolchain",
		},
		{
			"sourceCompatibility string single quotes",
			`sourceCompatibility = '11'`,
			11,
			"sourceCompatibility",
		},
		{
			"sourceCompatibility string double quotes",
			`sourceCompatibility = "21"`,
			21,
			"sourceCompatibility",
		},
		{
			"sourceCompatibility JavaVersion constant",
			`sourceCompatibility = JavaVersion.VERSION_17`,
			17,
			"sourceCompatibility",
		},
		{
			"targetCompatibility string",
			`targetCompatibility = '8'`,
			8,
			"targetCompatibility",
		},
		{
			"no match",
			`println "hello world"`,
			0,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ver, det := parseGradleContent(test.content)
			if ver != test.expectedVer || det != test.expectedDet {
				t.Errorf("parseGradleContent() = (%d, %q); want (%d, %q)", ver, det, test.expectedVer, test.expectedDet)
			}
		})
	}
}

func TestDetectProjectJavaVersion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "jdk-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test 1: No files
	ver, file := DetectProjectJavaVersion(tempDir)
	if ver != 0 || file != "" {
		t.Errorf("expected 0, \"\" when no files exist; got %d, %q", ver, file)
	}

	// Test 2: .java-version file
	err = os.WriteFile(filepath.Join(tempDir, ".java-version"), []byte("21\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write .java-version: %v", err)
	}
	ver, file = DetectProjectJavaVersion(tempDir)
	if ver != 21 || file != ".java-version" {
		t.Errorf("expected 21, \".java-version\"; got %d, %q", ver, file)
	}

	// Remove .java-version for next test
	os.Remove(filepath.Join(tempDir, ".java-version"))

	// Test 3: pom.xml with <java.version>
	pomContent := `
	<project>
		<properties>
			<java.version>17</java.version>
		</properties>
	</project>
	`
	err = os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte(pomContent), 0644)
	if err != nil {
		t.Fatalf("failed to write pom.xml: %v", err)
	}
	ver, file = DetectProjectJavaVersion(tempDir)
	if ver != 17 || file != "pom.xml (<java.version>)" {
		t.Errorf("expected 17, \"pom.xml (<java.version>)\"; got %d, %q", ver, file)
	}

	// Remove pom.xml
	os.Remove(filepath.Join(tempDir, "pom.xml"))

	// Test 4: build.gradle
	gradleContent := `
	plugins {
		id 'java'
	}
	sourceCompatibility = '11'
	`
	err = os.WriteFile(filepath.Join(tempDir, "build.gradle"), []byte(gradleContent), 0644)
	if err != nil {
		t.Fatalf("failed to write build.gradle: %v", err)
	}
	ver, file = DetectProjectJavaVersion(tempDir)
	if ver != 11 || file != "build.gradle (sourceCompatibility)" {
		t.Errorf("expected 11, \"build.gradle (sourceCompatibility)\"; got %d, %q", ver, file)
	}
}
