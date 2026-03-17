package analyzer

// mode determines the direction of the analysis check.
type mode int

const (
	// ModePointer checks that struct types are used by pointer, not by value (ptrstruct).
	ModePointer mode = iota
	// ModeValue checks that struct types are used by value, not by pointer (valuestruct).
	ModeValue
)

// linterName returns the linter name for this mode.
func (m mode) linterName() string {
	if m == ModeValue {
		return "valuestruct"
	}
	return "ptrstruct"
}

// doc returns the analyzer documentation string for this mode.
func (m mode) doc() string {
	if m == ModeValue {
		return "enforce value usage for struct-bearing declaration types"
	}
	return "enforce pointer usage for struct-bearing declaration types"
}
