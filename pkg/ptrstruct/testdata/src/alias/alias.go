package alias

type User struct {
	Name string
}

// Type alias should be resolved via Unalias and detected.
type UserAlias = User

func Save(u UserAlias) {} // want `parameter u uses value struct User; use \*User`
