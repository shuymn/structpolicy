package analyzer

import (
	"go/types"
	"strings"
	"testing"
)

func TestNewClassifier_InvalidPattern(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowPatterns = []string{"[invalid"}
	_, err := NewClassifier(cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestClassifier_IsAllowed_ByType(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowTypes = []string{"example.com/app.User"}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if !cls.IsAllowed(user) {
		t.Error("User should be allowed by type")
	}

	profile := newNamedStruct("example.com/app", "app", "Profile",
		types.NewVar(0, nil, "Bio", types.Typ[types.String]),
	)
	if cls.IsAllowed(profile) {
		t.Error("Profile should not be allowed")
	}
}

func TestClassifier_IsAllowed_ByPackage(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowPackages = []string{"example.com/external"}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	ext := newNamedStruct("example.com/external", "external", "Foo",
		types.NewVar(0, nil, "V", types.Typ[types.Int]),
	)
	if !cls.IsAllowed(ext) {
		t.Error("Foo from allowed package should be allowed")
	}

	internal := newNamedStruct("example.com/internal", "internal", "Bar",
		types.NewVar(0, nil, "V", types.Typ[types.Int]),
	)
	if cls.IsAllowed(internal) {
		t.Error("Bar from non-allowed package should not be allowed")
	}
}

func TestClassifier_IsAllowed_ByPattern(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowPatterns = []string{`\.Null[A-Z]\w*$`}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	nullStr := newNamedStruct("database/sql", "sql", "NullString",
		types.NewVar(0, nil, "String", types.Typ[types.String]),
		types.NewVar(0, nil, "Valid", types.Typ[types.Bool]),
	)
	if !cls.IsAllowed(nullStr) {
		t.Error("NullString should match pattern")
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(user) {
		t.Error("User should not match pattern")
	}
}

func TestClassifier_IsAllowed_ByStdlib(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowStdlib = true
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	timeValue := newNamedStruct("time", "time", "Time",
		types.NewVar(0, nil, "wall", types.Typ[types.Uint64]),
	)
	if !cls.IsAllowed(timeValue) {
		t.Error("time.Time should be allowed when stdlib exemption is enabled")
	}

	nullStr := newNamedStruct("database/sql", "sql", "NullString",
		types.NewVar(0, nil, "String", types.Typ[types.String]),
		types.NewVar(0, nil, "Valid", types.Typ[types.Bool]),
	)
	if !cls.IsAllowed(nullStr) {
		t.Error("database/sql.NullString should be allowed when stdlib exemption is enabled")
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(user) {
		t.Error("non-stdlib packages should not be allowed by stdlib exemption")
	}
}

func TestClassifier_IsAllowed_ByThirdParty(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowStdlib = false
	cfg.AllowThirdParty = true
	cls, err := newClassifier(cfg, "example.com/app")
	if err != nil {
		t.Fatal(err)
	}

	uuid := newNamedStruct("github.com/google/uuid", "uuid", "UUID",
		types.NewVar(0, nil, "Bytes", types.NewArray(types.Typ[types.Byte], 16)),
	)
	if !cls.IsAllowed(uuid) {
		t.Error("third-party type should be allowed when third-party exemption is enabled")
	}

	internal := newNamedStruct("example.com/app/internal/model", "model", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(internal) {
		t.Error("current-module package should not be allowed by third-party exemption")
	}

	timeValue := newNamedStruct("time", "time", "Time",
		types.NewVar(0, nil, "wall", types.Typ[types.Uint64]),
	)
	if cls.IsAllowed(timeValue) {
		t.Error("stdlib package should not be allowed by third-party exemption alone")
	}
}

func TestClassifier_EmptyConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(user) {
		t.Error("nothing should be allowed with empty config")
	}
}

func TestNewClassifier_EmptyPatternIgnored(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowPatterns = []string{""}

	cls, err := newClassifier(cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	named := newNamedStruct("example.com/foo", "foo", "ShouldNotBeAllowed",
		types.NewVar(0, nil, "X", types.Typ[types.Int]),
	)
	if cls.IsAllowed(named) {
		t.Error("empty pattern should be ignored, not match everything")
	}
}

func TestClassifier_IsAllowed_ReDoSPattern(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)
	cfg.AllowPatterns = []string{`(a+)+`}

	cls, err := newClassifier(cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	// Go's regexp uses RE2 (linear time), so this completes instantly
	// even with a pathological input designed to trigger backtracking engines.
	named := newNamedStruct(strings.Repeat("a", 100)+"b", "pkg", "B",
		types.NewVar(0, nil, "X", types.Typ[types.Int]),
	)
	cls.IsAllowed(named)
}
