package basic

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Pointer receiver is a violation in valuestruct.
func (u *User) Normalize() {} // want `receiver uses pointer to struct User; use User`

// Value receiver is OK.
func (u User) String() string { return u.Name }

// Pointer parameter is a violation.
func Save(u *User) {} // want `parameter u uses pointer to struct User; use User`

// Value parameter is OK.
func Process(u User) {}

// Pointer result is a violation.
func Load() *User { // want `result uses pointer to struct User; use User`
	return nil
}

// Value result is OK.
func Create() User {
	return User{}
}

// Field with pointer is a violation.
type Response struct {
	Meta *Profile // want `field Meta uses pointer to struct Profile; use Profile`
}

// Field with value is OK.
type Request struct {
	Meta Profile
}
