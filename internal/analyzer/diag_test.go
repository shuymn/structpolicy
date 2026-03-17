package analyzer

import "testing"

func TestFormatDiagnostic_ModePointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		position string
		v        Violation
		want     string
	}{
		{
			name:     "receiver by value",
			position: "receiver",
			v:        Violation{TypeName: "User"},
			want:     "receiver uses value struct User; use *User",
		},
		{
			name:     "parameter by value",
			position: "parameter req",
			v:        Violation{TypeName: "User"},
			want:     "parameter req uses value struct User; use *User",
		},
		{
			name:     "result by value",
			position: "result",
			v:        Violation{TypeName: "User"},
			want:     "result uses value struct User; use *User",
		},
		{
			name:     "field by value",
			position: "field Meta",
			v:        Violation{TypeName: "Meta"},
			want:     "field Meta uses value struct Meta; use *Meta",
		},
		{
			name:     "slice element",
			position: "parameter users",
			v:        Violation{Path: "slice element", TypeName: "User"},
			want:     "parameter users uses slice element User by value",
		},
		{
			name:     "map value",
			position: "field Index",
			v:        Violation{Path: "map value", TypeName: "User"},
			want:     "field Index uses map value User by value",
		},
		{
			name:     "pointer to slice element",
			position: "field Items",
			v:        Violation{Path: "pointer -> slice element", TypeName: "User"},
			want:     "field Items uses pointer -> slice element User by value",
		},
		{
			name:     "anonymous struct",
			position: "field Inner",
			v:        Violation{TypeName: "struct{...}"},
			want:     "field Inner uses value struct struct{...}; use *struct{...}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FormatDiagnostic(tt.position, &tt.v, ModePointer)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestFormatDiagnostic_ModeValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		position string
		v        Violation
		want     string
	}{
		{
			name:     "receiver by pointer",
			position: "receiver",
			v:        Violation{TypeName: "User"},
			want:     "receiver uses pointer to struct User; use User",
		},
		{
			name:     "parameter by pointer",
			position: "parameter req",
			v:        Violation{TypeName: "User"},
			want:     "parameter req uses pointer to struct User; use User",
		},
		{
			name:     "result by pointer",
			position: "result",
			v:        Violation{TypeName: "User"},
			want:     "result uses pointer to struct User; use User",
		},
		{
			name:     "field by pointer",
			position: "field Meta",
			v:        Violation{TypeName: "Meta"},
			want:     "field Meta uses pointer to struct Meta; use Meta",
		},
		{
			name:     "slice element",
			position: "parameter users",
			v:        Violation{Path: "slice element", TypeName: "User"},
			want:     "parameter users uses slice element User by pointer",
		},
		{
			name:     "map value",
			position: "field Index",
			v:        Violation{Path: "map value", TypeName: "User"},
			want:     "field Index uses map value User by pointer",
		},
		{
			name:     "anonymous struct pointer",
			position: "field Inner",
			v:        Violation{TypeName: "struct{...}"},
			want:     "field Inner uses pointer to struct struct{...}; use struct{...}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FormatDiagnostic(tt.position, &tt.v, ModeValue)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}
