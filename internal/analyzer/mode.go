package analyzer

// Mode determines the direction of the analysis check.
type Mode int

const (
	// ModePointer checks that struct types are used by pointer, not by value (ptrstruct).
	ModePointer Mode = iota
	// ModeValue checks that struct types are used by value, not by pointer (valuestruct).
	ModeValue
)

// LinterName returns the linter name for this mode.
func (m Mode) LinterName() string {
	if m == ModeValue {
		return "valuestruct"
	}
	return "ptrstruct"
}

// Doc returns the analyzer documentation string for this mode.
func (m Mode) Doc() string {
	if m == ModeValue {
		return "enforce value usage for struct-bearing declaration types"
	}
	return "enforce pointer usage for struct-bearing declaration types"
}
