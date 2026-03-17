package oneper

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Only the first violation should be reported per declaration.
// Two value-struct params, but only 1 diagnostic.
func SaveBoth(u User, p Profile) {} // want `parameter u uses value struct User; use \*User`

// Struct with multiple value-struct fields: only first flagged.
type Response struct {
	A User    // want `field A uses value struct User; use \*User`
	B Profile // no second diagnostic
}
