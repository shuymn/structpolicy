package analyzer

import (
	"slices"
	"testing"
)

func TestDefaultConfig_ModePointer(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModePointer)

	if cfg.Mode != ModePointer {
		t.Error("Mode should be ModePointer")
	}

	// Copy-reduction defaults.
	if !cfg.Receiver {
		t.Error("Receiver should default to true")
	}
	if cfg.Result {
		t.Error("Result should default to false")
	}
	if !cfg.Param {
		t.Error("Param should default to true")
	}
	if !cfg.Field {
		t.Error("Field should default to true")
	}
	if cfg.InterfaceMethod {
		t.Error("InterfaceMethod should default to false")
	}
	if cfg.FuncType {
		t.Error("FuncType should default to false")
	}
	if cfg.NamedType {
		t.Error("NamedType should default to false")
	}

	if cfg.MapKey {
		t.Error("MapKey should default to false")
	}
	if !cfg.SliceElem {
		t.Error("SliceElem should default to true")
	}
	if !cfg.MapValue {
		t.Error("MapValue should default to true")
	}
	if !cfg.ArrayElem {
		t.Error("ArrayElem should default to true")
	}
	if !cfg.ChanElem {
		t.Error("ChanElem should default to true")
	}

	// File filtering
	if !cfg.IgnoreGenerated {
		t.Error("IgnoreGenerated should default to true")
	}
	if cfg.IgnoreTests {
		t.Error("IgnoreTests should default to false")
	}

	// Suppression toggles
	if !cfg.HonorNolint {
		t.Error("HonorNolint should default to true")
	}
	if !cfg.HonorNolintAll {
		t.Error("HonorNolintAll should default to true")
	}
	if !cfg.AllowStdlib {
		t.Error("AllowStdlib should default to true")
	}
	if cfg.AllowThirdParty {
		t.Error("AllowThirdParty should default to false")
	}

	// Allowlists
	if cfg.AllowTypes != nil {
		t.Error("AllowTypes should default to nil")
	}
	if cfg.AllowPatterns != nil {
		t.Error("AllowPatterns should default to nil")
	}
	if cfg.AllowPackages != nil {
		t.Error("AllowPackages should default to nil")
	}
}

func TestDefaultConfig_ModeValue(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig(ModeValue)

	if cfg.Mode != ModeValue {
		t.Error("Mode should be ModeValue")
	}

	// Allocation- and indirection-reduction defaults.
	if cfg.Receiver {
		t.Error("Receiver should default to false")
	}
	if !cfg.Param {
		t.Error("Param should default to true")
	}
	if !cfg.Result {
		t.Error("Result should default to true")
	}
	if !cfg.Field {
		t.Error("Field should default to true")
	}
	if cfg.InterfaceMethod {
		t.Error("InterfaceMethod should default to false")
	}
	if cfg.FuncType {
		t.Error("FuncType should default to false")
	}
	if cfg.NamedType {
		t.Error("NamedType should default to false")
	}
	if !cfg.SliceElem {
		t.Error("SliceElem should default to true")
	}
	if !cfg.MapValue {
		t.Error("MapValue should default to true")
	}
	if cfg.MapKey {
		t.Error("MapKey should default to false")
	}
	if !cfg.ArrayElem {
		t.Error("ArrayElem should default to true")
	}
	if !cfg.ChanElem {
		t.Error("ChanElem should default to true")
	}
}

func TestStringListFlag_Set(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "normal", input: "foo,bar", want: []string{"foo", "bar"}},
		{name: "empty string", input: "", want: nil},
		{name: "double comma", input: "foo,,bar", want: []string{"foo", "bar"}},
		{name: "trailing comma", input: "foo,", want: []string{"foo"}},
		{name: "leading comma", input: ",foo", want: []string{"foo"}},
		{name: "all commas", input: ",,,", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var values []string
			f := stringListFlag{values: &values}
			if err := f.Set(tt.input); err != nil {
				t.Fatal(err)
			}
			if tt.want == nil {
				if values != nil {
					t.Errorf("got %v, want nil", values)
				}
				return
			}
			if !slices.Equal(values, tt.want) {
				t.Errorf("got %v, want %v", values, tt.want)
			}
		})
	}
}
