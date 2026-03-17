// Package typesigflags verifies that func-type and interface-method checks
// respect the param and result flags.
//
// ptrstruct default profile: func-type=true, interface-method=true,
// param=true, result=false.
//
// Only parameter-position violations should be reported; result-position
// violations must be suppressed because result=false.
package typesigflags

type User struct {
	Name string
}

// Function type with a value-struct result: OK (result=false).
type LoadFunc func() User

// Function type with a value-struct param: NG (param=true).
type SaveFunc func(u User) // want `function type SaveFunc parameter u uses value struct User; use \*User`

// Interface with a value-struct result: OK (result=false).
type Loader interface {
	Load() User
}

// Interface with a value-struct param: NG (param=true).
type Saver interface {
	Save(u User) // want `interface method Save parameter u uses value struct User; use \*User`
}
