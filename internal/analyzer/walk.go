package analyzer

import "go/types"

// Violation describes a struct usage occurrence found by the type walker.
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

// FindViolation walks t recursively and returns the first violation.
// ok is false if the type is clean.
// cls may be nil if no allowlist is configured.
func FindViolation(t types.Type, cfg *Config, cls *Classifier) (Violation, bool) {
	w := &walker{cfg: cfg, cls: cls, seen: make(map[*types.Named]bool)}
	return w.walk(t, "")
}

func (w *walker) walk(t types.Type, path string) (Violation, bool) {
	t = types.Unalias(t)

	switch tt := t.(type) {
	case *types.Pointer:
		return w.walkPointer(tt, path)
	case *types.Named:
		return w.walkNamed(tt, path)
	case *types.Struct:
		return w.walkStruct(tt, path)
	case *types.Slice:
		return w.walkSlice(tt, path)
	case *types.Map:
		return w.walkMap(tt, path)
	case *types.Array:
		return w.walkArray(tt, path)
	case *types.Chan:
		return w.walkChan(tt, path)
	default:
		return Violation{}, false
	}
}

func (w *walker) walkPointer(tt *types.Pointer, path string) (Violation, bool) {
	elem := types.Unalias(tt.Elem())

	if w.cfg.Mode == ModePointer {
		// ptrstruct: pointer to struct is OK, stop walking.
		if isStructType(elem) {
			return Violation{}, false
		}
	} else {
		// valuestruct: pointer to struct is a violation.
		if v, ok := w.checkStructViolation(elem, path); ok {
			return v, true
		}
	}

	return w.walk(elem, appendPath(path, pathPointer))
}

func (w *walker) walkNamed(tt *types.Named, path string) (Violation, bool) {
	if w.seen[tt] {
		return Violation{}, false
	}
	w.seen[tt] = true

	if w.cls != nil && w.cls.IsAllowed(tt) {
		return Violation{}, false
	}

	under := tt.Underlying()
	if st, ok := under.(*types.Struct); ok {
		if w.cfg.Mode == ModePointer {
			// ptrstruct: bare named struct is a violation.
			if st.NumFields() == 0 {
				return Violation{}, false
			}
			return Violation{
				Path:     path,
				TypeName: tt.Obj().Name(),
				Named:    tt,
			}, true
		}
		// valuestruct: bare named struct is OK.
		return Violation{}, false
	}

	return w.walk(under, path)
}

func (w *walker) walkStruct(tt *types.Struct, path string) (Violation, bool) {
	if w.cfg.Mode == ModeValue {
		// valuestruct: bare anonymous struct is OK.
		return Violation{}, false
	}
	// ptrstruct: anonymous struct with fields is a violation.
	if tt.NumFields() == 0 {
		return Violation{}, false
	}
	return Violation{
		Path:     path,
		TypeName: "struct{...}",
		Named:    nil,
	}, true
}

// checkStructViolation checks whether elem is a struct type and returns a
// violation if so. Used by ModeValue to detect pointer-to-struct.
//
// For *types.Named, it marks seen to avoid redundant work when walkPointer
// falls through to walk → walkNamed. The exception is non-struct underlying
// types, which must NOT be marked so walkNamed can recurse into them.
func (w *walker) checkStructViolation(elem types.Type, path string) (Violation, bool) {
	switch et := elem.(type) {
	case *types.Named:
		if w.seen[et] {
			return Violation{}, false
		}
		if w.cls != nil && w.cls.IsAllowed(et) {
			w.seen[et] = true // walkNamed would also return nil
			return Violation{}, false
		}
		st, ok := et.Underlying().(*types.Struct)
		if !ok {
			return Violation{}, false // don't mark seen — walkNamed must recurse
		}
		w.seen[et] = true
		if st.NumFields() == 0 {
			return Violation{}, false
		}
		return Violation{
			Path:     path,
			TypeName: et.Obj().Name(),
			Named:    et,
		}, true
	case *types.Struct:
		if et.NumFields() == 0 {
			return Violation{}, false
		}
		return Violation{
			Path:     path,
			TypeName: "struct{...}",
			Named:    nil,
		}, true
	}
	return Violation{}, false
}

func (w *walker) walkSlice(tt *types.Slice, path string) (Violation, bool) {
	if !w.cfg.SliceElem {
		return Violation{}, false
	}
	return w.walk(tt.Elem(), appendPath(path, pathSliceElement))
}

func (w *walker) walkMap(tt *types.Map, path string) (Violation, bool) {
	if w.cfg.MapKey {
		if v, ok := w.walk(tt.Key(), appendPath(path, pathMapKey)); ok {
			return v, true
		}
	}
	if !w.cfg.MapValue {
		return Violation{}, false
	}
	return w.walk(tt.Elem(), appendPath(path, pathMapValue))
}

func (w *walker) walkArray(tt *types.Array, path string) (Violation, bool) {
	if !w.cfg.ArrayElem {
		return Violation{}, false
	}
	return w.walk(tt.Elem(), appendPath(path, pathArrayElement))
}

func (w *walker) walkChan(tt *types.Chan, path string) (Violation, bool) {
	if !w.cfg.ChanElem {
		return Violation{}, false
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
