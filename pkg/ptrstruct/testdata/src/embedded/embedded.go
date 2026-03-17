package embedded

type Inner struct {
	Value string
}

// Embedded struct field by value: NG
type Outer struct {
	Inner // want `field Inner uses value struct Inner; use \*Inner`
}

// Embedded pointer: OK
type OuterOK struct {
	*Inner
}
