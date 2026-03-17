package skiptests

type User struct {
	Name string
}

// Non-test file: always checked regardless of IgnoreTests.
func Save(u User) {} // want `parameter u uses value struct User; use \*User`
