package nested

type User struct {
	Name string
}

// Nested containers: slice of map of struct
func DeepNest(data []map[string]User) {} // want `parameter data uses slice element -> map value User by value`

// Pointer wrapping a nested container
func PtrDeepNest(
	data *[]map[string]User, // want `parameter data uses pointer -> slice element -> map value User by value`
) {
}
