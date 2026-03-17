package analyzer

import "testing"

func TestNolintMatcher_Ptrstruct(t *testing.T) {
	t.Parallel()

	m := nolintMatcher{linterName: "ptrstruct"}

	tests := []struct {
		name        string
		honorAll    bool
		text        string
		wantMatched bool
	}{
		{name: "exact ptrstruct", text: "//nolint:ptrstruct", wantMatched: true},
		{name: "ptrstruct with space", text: "// nolint:ptrstruct", wantMatched: true},
		{name: "ptrstruct with reason", text: "//nolint:ptrstruct // legacy API", wantMatched: true},
		{name: "nolint all honored", honorAll: true, text: "//nolint:all", wantMatched: true},
		{name: "nolint all not honored", text: "//nolint:all", wantMatched: false},
		{name: "nolint all with space", honorAll: true, text: "// nolint:all", wantMatched: true},
		{name: "other linter", text: "//nolint:govet", wantMatched: false},
		{name: "multiple linters including ptrstruct", text: "//nolint:govet,ptrstruct", wantMatched: true},
		{name: "regular comment", text: "// this is a comment", wantMatched: false},
		{name: "empty comment", text: "//", wantMatched: false},
		{name: "nolint bare", text: "//nolint", wantMatched: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := m
			m.honorAll = tt.honorAll
			got := m.match(tt.text)
			if got != tt.wantMatched {
				t.Errorf(
					"nolintMatcher{%q, honorAll=%v}.match(%q) = %v, want %v",
					m.linterName, m.honorAll, tt.text, got, tt.wantMatched,
				)
			}
		})
	}
}

func TestNolintMatcher_Valuestruct(t *testing.T) {
	t.Parallel()

	m := nolintMatcher{linterName: "valuestruct"}

	tests := []struct {
		name        string
		honorAll    bool
		text        string
		wantMatched bool
	}{
		{name: "exact valuestruct", text: "//nolint:valuestruct", wantMatched: true},
		{name: "valuestruct with space", text: "// nolint:valuestruct", wantMatched: true},
		{name: "valuestruct with reason", text: "//nolint:valuestruct // legacy API", wantMatched: true},
		{name: "nolint all honored", honorAll: true, text: "//nolint:all", wantMatched: true},
		{name: "nolint all not honored", text: "//nolint:all", wantMatched: false},
		{name: "ptrstruct not matched", text: "//nolint:ptrstruct", wantMatched: false},
		{name: "multiple linters including valuestruct", text: "//nolint:govet,valuestruct", wantMatched: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := m
			m.honorAll = tt.honorAll
			got := m.match(tt.text)
			if got != tt.wantMatched {
				t.Errorf(
					"nolintMatcher{%q, honorAll=%v}.match(%q) = %v, want %v",
					m.linterName, m.honorAll, tt.text, got, tt.wantMatched,
				)
			}
		})
	}
}
