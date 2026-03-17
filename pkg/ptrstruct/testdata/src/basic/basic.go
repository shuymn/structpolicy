package basic

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// R1: receiver must be pointer
func (u User) Normalize() {} // want `receiver uses value struct User; use \*User`

// R2: parameter must not be value-struct
func Save(u User) {} // want `parameter u uses value struct User; use \*User`

// R3: result must not be value-struct
func Load() User { // want `result uses value struct User; use \*User`
	return User{}
}

// R4: field must not contain value-struct directly
type Response struct {
	Meta Profile // want `field Meta uses value struct Profile; use \*Profile`
}
