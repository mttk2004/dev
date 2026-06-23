package registry

import (
	"testing"
)

// ---------------------------------------------------------------------------
// ID constants
// ---------------------------------------------------------------------------

func TestIDConstants_NonEmpty(t *testing.T) {
	ids := []struct {
		name string
		val  string
	}{
		{"IDNode", IDNode},
		{"IDBun", IDBun},
		{"IDComposer", IDComposer},
		{"IDJDK", IDJDK},
		{"IDGo", IDGo},
		{"IDPHP", IDPHP},
		{"IDDocker", IDDocker},
		{"IDPostgreSQL", IDPostgreSQL},
		{"IDRedis", IDRedis},
		{"IDNginx", IDNginx},
		{"IDPython", IDPython},
		{"IDMaven", IDMaven},
		{"IDMariaDB", IDMariaDB},
	}

	for _, tt := range ids {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val == "" {
				t.Errorf("%s is empty", tt.name)
			}
		})
	}
}

func TestIDConstants_Unique(t *testing.T) {
	all := []string{
		IDNode, IDBun, IDComposer, IDJDK, IDGo, IDPHP,
		IDDocker, IDPostgreSQL, IDRedis, IDNginx,
		IDPython, IDMaven, IDMariaDB,
	}

	seen := make(map[string]bool, len(all))
	for _, id := range all {
		if seen[id] {
			t.Errorf("duplicate ID constant: %q", id)
		}
		seen[id] = true
	}
}

// ---------------------------------------------------------------------------
// Packages slice integrity
// ---------------------------------------------------------------------------

func TestPackages_NotEmpty(t *testing.T) {
	if len(Packages) == 0 {
		t.Fatal("Packages slice is empty")
	}
}

func TestPackages_UniqueIDs(t *testing.T) {
	seen := make(map[string]bool, len(Packages))
	for _, p := range Packages {
		if seen[p.ID] {
			t.Errorf("duplicate package ID in Packages: %q", p.ID)
		}
		seen[p.ID] = true
	}
}

func TestPackages_RequiredFields(t *testing.T) {
	for _, p := range Packages {
		t.Run(p.ID, func(t *testing.T) {
			if p.ID == "" {
				t.Error("ID is empty")
			}
			if p.DisplayName == "" {
				t.Errorf("DisplayName is empty for %q", p.ID)
			}
			if p.CheckCmd == "" {
				t.Errorf("CheckCmd is empty for %q", p.ID)
			}
			if p.Install == nil {
				t.Errorf("Install func is nil for %q", p.ID)
			}
			if p.Remove == nil {
				t.Errorf("Remove func is nil for %q", p.ID)
			}
			if p.Update == nil {
				t.Errorf("Update func is nil for %q", p.ID)
			}
		})
	}
}

func TestPackages_IDsMatchConstants(t *testing.T) {
	constants := map[string]bool{
		IDNode: true, IDBun: true, IDComposer: true, IDJDK: true,
		IDGo: true, IDPHP: true, IDDocker: true, IDPostgreSQL: true,
		IDRedis: true, IDNginx: true, IDPython: true, IDMaven: true,
		IDMariaDB: true,
	}

	for _, p := range Packages {
		if !constants[p.ID] {
			t.Errorf("package %q has an ID that does not match any ID* constant", p.ID)
		}
	}
}

func TestPackages_NodeHasAltCheckCmd(t *testing.T) {
	p, err := LookupByID(IDNode)
	if err != nil {
		t.Fatalf("LookupByID(%q) failed: %v", IDNode, err)
	}
	if p.AltCheckCmd == "" {
		t.Error("Node package should have AltCheckCmd set (expected \"node\")")
	}
	if p.AltCheckCmd != "node" {
		t.Errorf("Node AltCheckCmd = %q, want %q", p.AltCheckCmd, "node")
	}
}

func TestPackages_ServicesOnlyOnRelevantPackages(t *testing.T) {
	servicePackages := map[string]bool{
		IDDocker: true, IDPostgreSQL: true, IDRedis: true,
		IDNginx: true, IDMariaDB: true,
	}

	for _, p := range Packages {
		hasServices := len(p.Services) > 0
		shouldHave := servicePackages[p.ID]

		if hasServices && !shouldHave {
			t.Errorf("package %q has Services %v but is not expected to", p.ID, p.Services)
		}
		if !hasServices && shouldHave {
			t.Errorf("package %q should have Services but has none", p.ID)
		}
	}
}

// ---------------------------------------------------------------------------
// LookupByID
// ---------------------------------------------------------------------------

func TestLookupByID_AllKnownIDs(t *testing.T) {
	for _, p := range Packages {
		t.Run(p.ID, func(t *testing.T) {
			found, err := LookupByID(p.ID)
			if err != nil {
				t.Fatalf("LookupByID(%q) returned error: %v", p.ID, err)
			}
			if found.ID != p.ID {
				t.Errorf("LookupByID(%q).ID = %q", p.ID, found.ID)
			}
			if found.DisplayName != p.DisplayName {
				t.Errorf("LookupByID(%q).DisplayName = %q, want %q", p.ID, found.DisplayName, p.DisplayName)
			}
		})
	}
}

func TestLookupByID_ReturnsPointerToOriginal(t *testing.T) {
	// Ensure we get a pointer into the Packages slice, not a copy
	p, err := LookupByID(IDNode)
	if err != nil {
		t.Fatal(err)
	}

	// Find the original by index
	var original *Package
	for i := range Packages {
		if Packages[i].ID == IDNode {
			original = &Packages[i]
			break
		}
	}

	if p != original {
		t.Error("LookupByID should return a pointer to the element in the Packages slice, not a copy")
	}
}

func TestLookupByID_NotFound(t *testing.T) {
	notFound := []string{
		"",
		"nonexistent",
		"Node",   // case-sensitive
		"NODE",
		" node",  // leading space
		"node ",  // trailing space
	}

	for _, id := range notFound {
		t.Run(id, func(t *testing.T) {
			p, err := LookupByID(id)
			if err == nil {
				t.Errorf("LookupByID(%q) = %v, want error", id, p)
			}
			if p != nil {
				t.Errorf("LookupByID(%q) returned non-nil package on error", id)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// IDs
// ---------------------------------------------------------------------------

func TestIDs_Length(t *testing.T) {
	ids := IDs()
	if len(ids) != len(Packages) {
		t.Errorf("IDs() returned %d items, want %d", len(ids), len(Packages))
	}
}

func TestIDs_Order(t *testing.T) {
	ids := IDs()
	for i, id := range ids {
		if id != Packages[i].ID {
			t.Errorf("IDs()[%d] = %q, want %q (must preserve Packages order)", i, id, Packages[i].ID)
		}
	}
}

func TestIDs_NoEmpty(t *testing.T) {
	for i, id := range IDs() {
		if id == "" {
			t.Errorf("IDs()[%d] is empty", i)
		}
	}
}

func TestIDs_ReturnsCopy(t *testing.T) {
	ids1 := IDs()
	ids2 := IDs()

	// Mutating one should not affect the other
	if len(ids1) > 0 {
		ids1[0] = "mutated"
		if ids2[0] == "mutated" {
			t.Error("IDs() should return a new slice each time, not a shared reference")
		}
	}
}

// ---------------------------------------------------------------------------
// DiagnosticMessage
// ---------------------------------------------------------------------------

func TestDiagnosticMessage_InstalledSimplePackage(t *testing.T) {
	// Use a package whose CheckCmd we know exists on any system: "bash"
	pkg := Package{
		ID:       "test-sh",
		CheckCmd: "bash",
	}

	msg := pkg.DiagnosticMessage()
	expected := "bash is installed and in PATH"
	if msg != expected {
		t.Errorf("DiagnosticMessage() = %q, want %q", msg, expected)
	}
}

func TestDiagnosticMessage_InstalledWithAltCmd(t *testing.T) {
	// Both primary and alt exist
	pkg := Package{
		ID:          "test-dual",
		CheckCmd:    "bash",
		AltCheckCmd: "ls",
	}

	msg := pkg.DiagnosticMessage()
	expected := "bash and ls are installed and in PATH"
	if msg != expected {
		t.Errorf("DiagnosticMessage() = %q, want %q", msg, expected)
	}
}

func TestDiagnosticMessage_InstalledButAltMissing(t *testing.T) {
	pkg := Package{
		ID:          "test-partial",
		CheckCmd:    "bash",
		AltCheckCmd: "this_does_not_exist_xyz_987",
	}

	msg := pkg.DiagnosticMessage()
	expected := "bash is installed, but this_does_not_exist_xyz_987 is not in PATH. Try restarting your terminal."
	if msg != expected {
		t.Errorf("DiagnosticMessage() = %q, want %q", msg, expected)
	}
}

func TestDiagnosticMessage_NotInstalled(t *testing.T) {
	pkg := Package{
		ID:       "test-missing",
		CheckCmd: "this_does_not_exist_xyz_123",
	}

	msg := pkg.DiagnosticMessage()
	expected := "this_does_not_exist_xyz_123 is missing or not in PATH"
	if msg != expected {
		t.Errorf("DiagnosticMessage() = %q, want %q", msg, expected)
	}
}

// ---------------------------------------------------------------------------
// IsInstalled
// ---------------------------------------------------------------------------

func TestIsInstalled_Exists(t *testing.T) {
	pkg := Package{CheckCmd: "bash"}
	if !pkg.IsInstalled() {
		t.Error("IsInstalled() = false for 'sh', want true")
	}
}

func TestIsInstalled_NotExists(t *testing.T) {
	pkg := Package{CheckCmd: "this_does_not_exist_xyz_456"}
	if pkg.IsInstalled() {
		t.Error("IsInstalled() = true for nonexistent command, want false")
	}
}

// ---------------------------------------------------------------------------
// IsFullyOperational
// ---------------------------------------------------------------------------

func TestIsFullyOperational_SimpleInstalled(t *testing.T) {
	pkg := Package{CheckCmd: "bash"}
	if !pkg.IsFullyOperational() {
		t.Error("IsFullyOperational() = false for 'sh' (no AltCheckCmd), want true")
	}
}

func TestIsFullyOperational_NotInstalled(t *testing.T) {
	pkg := Package{CheckCmd: "this_does_not_exist_xyz_789"}
	if pkg.IsFullyOperational() {
		t.Error("IsFullyOperational() = true for nonexistent command, want false")
	}
}

func TestIsFullyOperational_PrimaryOK_AltOK(t *testing.T) {
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "ls"}
	if !pkg.IsFullyOperational() {
		t.Error("IsFullyOperational() = false when both sh and ls exist, want true")
	}
}

func TestIsFullyOperational_PrimaryOK_AltMissing(t *testing.T) {
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "this_does_not_exist_xyz_alt"}
	if pkg.IsFullyOperational() {
		t.Error("IsFullyOperational() = true when alt command missing, want false")
	}
}

func TestIsFullyOperational_PrimaryMissing_AltExists(t *testing.T) {
	pkg := Package{CheckCmd: "this_does_not_exist_xyz_primary", AltCheckCmd: "ls"}
	if pkg.IsFullyOperational() {
		t.Error("IsFullyOperational() = true when primary missing, want false")
	}
}

// ---------------------------------------------------------------------------
// GetVersion / GetPath (smoke tests — actual output depends on system)
// ---------------------------------------------------------------------------

func TestGetVersion_KnownCommand(t *testing.T) {
	// "bash" should return *something* that isn't empty
	pkg := Package{CheckCmd: "bash"}
	ver := pkg.GetVersion()
	if ver == "" {
		t.Error("GetVersion() = \"\" for sh, want non-empty")
	}
}

func TestGetVersion_PrefersAltCmd(t *testing.T) {
	// When AltCheckCmd is set and exists, GetVersion should use it
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "ls"}
	ver := pkg.GetVersion()
	// We can't assert the exact version, but it should not be empty
	if ver == "" {
		t.Error("GetVersion() = \"\" when AltCheckCmd exists, want non-empty")
	}
}

func TestGetPath_KnownCommand(t *testing.T) {
	pkg := Package{CheckCmd: "bash"}
	path := pkg.GetPath()
	if path == "" {
		t.Error("GetPath() = \"\" for sh, want non-empty absolute path")
	}
	if len(path) > 0 && path[0] != '/' {
		t.Errorf("GetPath() = %q, want absolute path starting with /", path)
	}
}

func TestGetPath_PrefersAltCmd(t *testing.T) {
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "ls"}
	path := pkg.GetPath()
	if path == "" {
		t.Error("GetPath() = \"\" when AltCheckCmd exists, want non-empty")
	}
}

func TestGetPath_AltMissing_FallsToPrimary(t *testing.T) {
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "this_does_not_exist_fallback"}
	path := pkg.GetPath()
	// Should fall back to "bash" path
	if path == "" {
		t.Error("GetPath() = \"\" when AltCheckCmd missing, should fallback to primary")
	}
}

func TestGetVersion_AltMissing_FallsToPrimary(t *testing.T) {
	pkg := Package{CheckCmd: "bash", AltCheckCmd: "this_does_not_exist_fallback"}
	ver := pkg.GetVersion()
	if ver == "" {
		t.Error("GetVersion() = \"\" when AltCheckCmd missing, should fallback to primary")
	}
}
