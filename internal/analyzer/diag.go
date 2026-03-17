package analyzer

// formatDiagnostic produces the diagnostic message for a violation.
// position is a human label such as "receiver", "parameter req", "field Meta".
func formatDiagnostic(position string, v *violation, mode mode) string {
	if mode == ModeValue {
		if v.path == "" {
			return position + " uses pointer to struct " + v.typeName + "; use " + v.typeName
		}
		return position + " uses " + v.path + " " + v.typeName + " by pointer"
	}
	if v.path == "" {
		return position + " uses value struct " + v.typeName + "; use *" + v.typeName
	}
	return position + " uses " + v.path + " " + v.typeName + " by value"
}
