package ptrstruct

import "go/types"

// Violation describes a struct-by-value occurrence found by the type walker.
type Violation struct {
	Path     string       // Human-readable path, e.g. "slice element", "map value".
	TypeName string       // Short name of the struct type found, e.g. "User".
	Named    *types.Named // The named type if the violation is on a named struct; nil for anonymous structs.
}

// walker holds traversal state for a single FindViolation call.
// Cycle detection uses seen keyed on *types.Named because Go's type system
// only permits recursion through named types.
type walker struct {
	cfg  *Config
	cls  *Classifier
	seen map[*types.Named]bool
}

// FindViolation walks t recursively and returns the first value-struct
// violation, or nil if the type is clean.
// cls may be nil if no allowlist is configured.
func FindViolation(t types.Type, cfg *Config, cls *Classifier) *Violation {
	w := &walker{cfg: cfg, cls: cls, seen: make(map[*types.Named]bool)}
	return w.walk(t, "")
}

func (w *walker) walk(t types.Type, path string) *Violation {
	t = types.Unalias(t)

	switch tt := t.(type) {
	case *types.Pointer:
		return w.walkPointer(tt, path)
	case *types.Named:
		return w.walkNamed(tt, path)
	case *types.Struct:
		return walkStruct(tt, path)
	case *types.Slice:
		return w.walkSlice(tt, path)
	case *types.Map:
		return w.walkMap(tt, path)
	case *types.Array:
		return w.walkArray(tt, path)
	case *types.Chan:
		return w.walkChan(tt, path)
	default:
		return nil
	}
}

func (w *walker) walkPointer(tt *types.Pointer, path string) *Violation {
	elem := types.Unalias(tt.Elem())

	if isStructType(elem) {
		return nil
	}

	return w.walk(elem, appendPath(path, pathPointer))
}

func (w *walker) walkNamed(tt *types.Named, path string) *Violation {
	if w.seen[tt] {
		return nil
	}
	w.seen[tt] = true

	if w.cls != nil && w.cls.IsAllowed(tt) {
		return nil
	}

	under := tt.Underlying()
	if st, ok := under.(*types.Struct); ok {
		if st.NumFields() == 0 {
			return nil
		}
		return &Violation{
			Path:     path,
			TypeName: tt.Obj().Name(),
			Named:    tt,
		}
	}

	return w.walk(under, path)
}

func walkStruct(tt *types.Struct, path string) *Violation {
	if tt.NumFields() == 0 {
		return nil
	}
	return &Violation{
		Path:     path,
		TypeName: "struct{...}",
		Named:    nil,
	}
}

func (w *walker) walkSlice(tt *types.Slice, path string) *Violation {
	if !w.cfg.SliceElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathSliceElement))
}

func (w *walker) walkMap(tt *types.Map, path string) *Violation {
	if w.cfg.MapKey {
		if v := w.walk(tt.Key(), appendPath(path, pathMapKey)); v != nil {
			return v
		}
	}
	if !w.cfg.MapValue {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathMapValue))
}

func (w *walker) walkArray(tt *types.Array, path string) *Violation {
	if !w.cfg.ArrayElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathArrayElement))
}

func (w *walker) walkChan(tt *types.Chan, path string) *Violation {
	if !w.cfg.ChanElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathChanElement))
}

// Path segment constants for violation paths.
const (
	pathPointer      = "pointer"
	pathSliceElement = "slice element"
	pathMapKey       = "map key"
	pathMapValue     = "map value"
	pathArrayElement = "array element"
	pathChanElement  = "chan element"
)

// isStructType reports whether t is a struct (named or anonymous).
func isStructType(t types.Type) bool {
	switch tt := t.(type) {
	case *types.Named:
		_, ok := tt.Underlying().(*types.Struct)
		return ok
	case *types.Struct:
		return true
	default:
		return false
	}
}

func appendPath(base, segment string) string {
	if base == "" {
		return segment
	}
	return base + " -> " + segment
}
