// Package typesigflags verifies that func-type and interface-method checks
// respect the param and result flags.
//
// This fixture is run with func-type=true, interface-method=true,
// param=false, result=true.
//
// Only result-position violations should be reported; parameter-position
// violations must be suppressed because param=false.
package typesigflags

type User struct {
	Name string
}

// Function type with a pointer-struct param: OK (param=false).
type SaveFunc func(u *User)

// Function type with a pointer-struct result: NG (result=true).
type LoadFunc func() *User // want `function type LoadFunc result uses pointer to struct User; use User`

// Interface with a pointer-struct param: OK (param=false).
type Saver interface {
	Save(u *User)
}

// Interface with a pointer-struct result: NG (result=true).
type Loader interface {
	Load() *User // want `interface method Load result uses pointer to struct User; use User`
}
