package ok

type User struct {
	Name string
}

type Profile struct {
	Bio string
}

// Pointer receiver: OK
func (u *User) Normalize() {}

// Pointer parameter: OK
func Save(u *User) {}

// Pointer result: OK
func Load() *User {
	return &User{}
}

// Pointer field: OK
type Response struct {
	Meta *Profile
}

// Basic types: OK
func Process(name string, count int) (bool, error) {
	return true, nil
}

// Interface parameter: OK
func Handle(v interface{}) {}

// Slice of pointers: OK
func SaveAll(users []*User) {}

// Map of pointers: OK
func Index(m map[string]*User) {}

// Empty struct: OK (signal type)
type Token struct{}

func Send(t Token) {}

type Container struct {
	Signal Token
}
