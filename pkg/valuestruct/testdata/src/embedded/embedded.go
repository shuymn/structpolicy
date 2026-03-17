package embedded

type Inner struct {
	Value string
}

// Embedded struct field by pointer: NG
type Outer struct {
	*Inner // want `field Inner uses pointer to struct Inner; use Inner`
}

// Embedded value: OK
type OuterOK struct {
	Inner
}
