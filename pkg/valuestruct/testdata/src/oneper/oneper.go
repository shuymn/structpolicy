package oneper

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Only the first violation should be reported per declaration.
// Two pointer-struct params, but only 1 diagnostic.
func SaveBoth(u *User, p *Profile) {} // want `parameter u uses pointer to struct User; use User`

// Struct with multiple pointer-struct fields: only first flagged.
type Response struct {
	A *User    // want `field A uses pointer to struct User; use User`
	B *Profile // no second diagnostic
}
