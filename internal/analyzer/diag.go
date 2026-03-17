package analyzer

// FormatDiagnostic produces the diagnostic message for a violation.
// position is a human label such as "receiver", "parameter req", "field Meta".
func FormatDiagnostic(position string, v *Violation, mode Mode) string {
	if mode == ModeValue {
		if v.Path == "" {
			return position + " uses pointer to struct " + v.TypeName + "; use " + v.TypeName
		}
		return position + " uses " + v.Path + " " + v.TypeName + " by pointer"
	}
	if v.Path == "" {
		return position + " uses value struct " + v.TypeName + "; use *" + v.TypeName
	}
	return position + " uses " + v.Path + " " + v.TypeName + " by value"
}
