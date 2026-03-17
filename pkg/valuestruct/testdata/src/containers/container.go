package containers

type User struct {
	Name string
}

// Slice element by pointer: NG
func ProcessAll(users []*User) {} // want `parameter users uses slice element User by pointer`

// Map value by pointer: NG
func Lookup(m map[string]*User) {} // want `parameter m uses map value User by pointer`

// Pointer to slice of pointer struct: NG
func ProcessPtr(users *[]*User) {} // want `parameter users uses pointer -> slice element User by pointer`

// Pointer to map of pointer struct: NG
func LookupPtr(m *map[string]*User) {} // want `parameter m uses pointer -> map value User by pointer`

// Struct fields with containers
type Response struct {
	Items []*User // want `field Items uses slice element User by pointer`
}

type Index struct {
	Data map[string]*User // want `field Data uses map value User by pointer`
}

// Slice of values: OK
func ProcessOK(users []User) {}

// Map of values: OK
func LookupOK(m map[string]User) {}
