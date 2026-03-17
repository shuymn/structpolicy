package typeblock

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Type block: per-TypeSpec nolint on A, but B should still flag.
type (
	//nolint:ptrstruct // suppress A
	A struct {
		X User
	}

	B struct {
		Y User // want `field Y uses value struct User; use \*User`
	}
)
