package generics

type User struct {
	Name string
}

// Type parameters should pass through without error.
func First[T any](items []T) T {
	return items[0]
}

// Constrained type parameter: OK
func Get[T comparable](m map[T]string, key T) string {
	return m[key]
}
