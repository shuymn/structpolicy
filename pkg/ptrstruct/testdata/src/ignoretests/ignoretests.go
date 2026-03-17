package ignoretests

type User struct {
	Name string
}

// This non-test file should always be checked.
func Save(u User) {} // want `parameter u uses value struct User; use \*User`
