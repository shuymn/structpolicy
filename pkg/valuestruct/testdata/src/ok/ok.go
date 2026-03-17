package ok

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Value receiver: OK
func (u User) String() string { return u.Name }

// Value parameter: OK
func Process(u User) {}

// Value result: OK
func Create() User { return User{} }

// Value field: OK
type Wrapper struct {
	Inner User
}

// Basic types: OK
func Handle(name string, count int) (bool, error) {
	return true, nil
}

// Interface parameter: OK
func HandleAny(v interface{}) {}

// Slice of values: OK
func SaveAll(users []User) {}

// Map of values: OK
func Index(m map[string]User) {}

// Empty struct: OK (signal type)
type Token struct{}

func Send(t *Token) {}

type Container struct {
	Signal *Token
}
