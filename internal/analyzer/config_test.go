package analyzer

import (
	"slices"
	"testing"
)

func TestDefaultConfig_ModePointer(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(ModePointer)

	if cfg.mode != ModePointer {
		t.Error("mode should be ModePointer")
	}

	// Pointer-leaning defaults for an initial refactor pass.
	if !cfg.receiver {
		t.Error("receiver should default to true")
	}
	if cfg.result {
		t.Error("result should default to false")
	}
	if !cfg.param {
		t.Error("param should default to true")
	}
	if !cfg.field {
		t.Error("field should default to true")
	}
	if !cfg.interfaceMethod {
		t.Error("interfaceMethod should default to true")
	}
	if !cfg.funcType {
		t.Error("funcType should default to true")
	}
	if cfg.namedType {
		t.Error("namedType should default to false")
	}

	if cfg.mapKey {
		t.Error("mapKey should default to false")
	}
	if cfg.sliceElem {
		t.Error("sliceElem should default to false")
	}
	if cfg.mapValue {
		t.Error("mapValue should default to false")
	}
	if cfg.arrayElem {
		t.Error("arrayElem should default to false")
	}
	if cfg.chanElem {
		t.Error("chanElem should default to false")
	}

	// File filtering
	if !cfg.ignoreGenerated {
		t.Error("ignoreGenerated should default to true")
	}
	if cfg.ignoreTests {
		t.Error("ignoreTests should default to false")
	}

	// Suppression toggles
	if !cfg.honorNolint {
		t.Error("honorNolint should default to true")
	}
	if !cfg.honorNolintAll {
		t.Error("honorNolintAll should default to true")
	}
	if !cfg.allowStdlib {
		t.Error("allowStdlib should default to true")
	}
	if cfg.allowThirdParty {
		t.Error("allowThirdParty should default to false")
	}

	// Allowlists
	if cfg.allowTypes != nil {
		t.Error("allowTypes should default to nil")
	}
	if cfg.allowPatterns != nil {
		t.Error("allowPatterns should default to nil")
	}
	if cfg.allowPackages != nil {
		t.Error("allowPackages should default to nil")
	}
}

func TestDefaultConfig_ModeValue(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(ModeValue)

	if cfg.mode != ModeValue {
		t.Error("mode should be ModeValue")
	}

	// Value-leaning defaults for an initial refactor pass.
	if cfg.receiver {
		t.Error("receiver should default to false")
	}
	if cfg.param {
		t.Error("param should default to false")
	}
	if !cfg.result {
		t.Error("result should default to true")
	}
	if cfg.field {
		t.Error("field should default to false")
	}
	if cfg.interfaceMethod {
		t.Error("interfaceMethod should default to false")
	}
	if cfg.funcType {
		t.Error("funcType should default to false")
	}
	if !cfg.namedType {
		t.Error("namedType should default to true")
	}
	if !cfg.sliceElem {
		t.Error("sliceElem should default to true")
	}
	if !cfg.mapValue {
		t.Error("mapValue should default to true")
	}
	if !cfg.mapKey {
		t.Error("mapKey should default to true")
	}
	if !cfg.arrayElem {
		t.Error("arrayElem should default to true")
	}
	if !cfg.chanElem {
		t.Error("chanElem should default to true")
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
