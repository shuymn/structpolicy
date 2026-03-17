package ignoretests

// Helper in test file: not a test function but still in _test.go.
// When IgnoreTests is false (default), this should flag.
func helperSave(u User) {} // want `parameter u uses value struct User; use \*User`
