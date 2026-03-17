package containers

type User struct {
	Name string
}

// Slice element by value: NG
func ProcessAll(users []User) {} // want `parameter users uses slice element User by value`

// Map value by value: NG
func Lookup(m map[string]User) {} // want `parameter m uses map value User by value`

// Pointer to slice of value struct: NG
func ProcessPtr(users *[]User) {} // want `parameter users uses pointer -> slice element User by value`

// Pointer to map of value struct: NG
func LookupPtr(m *map[string]User) {} // want `parameter m uses pointer -> map value User by value`

// Struct fields with containers
type Response struct {
	Items []User // want `field Items uses slice element User by value`
}

type Index struct {
	Data map[string]User // want `field Data uses map value User by value`
}

// Slice of pointers: OK
func ProcessOK(users []*User) {}

// Map of pointers: OK
func LookupOK(m map[string]*User) {}
