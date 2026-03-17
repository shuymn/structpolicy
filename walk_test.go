package ptrstruct

import (
	"go/types"
	"testing"
)

func newNamedStruct(pkgPath, pkgName, typeName string, fields ...*types.Var) *types.Named {
	pkg := types.NewPackage(pkgPath, pkgName)
	st := types.NewStruct(fields, nil)
	obj := types.NewTypeName(0, pkg, typeName, nil)
	named := types.NewNamed(obj, st, nil)
	return named
}

func TestFindViolation(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig()

	tests := []struct {
		name    string
		typ     types.Type
		wantNil bool
	}{
		{name: "basic type int", typ: types.Typ[types.Int], wantNil: true},
		{name: "basic type string", typ: types.Typ[types.String], wantNil: true},
		{name: "interface{}", typ: types.NewInterfaceType(nil, nil), wantNil: true},
		{name: "named struct by value", typ: user, wantNil: false},
		{name: "pointer to named struct", typ: types.NewPointer(user), wantNil: true},
		{name: "slice of named struct", typ: types.NewSlice(user), wantNil: false},
		{name: "slice of pointer to named struct", typ: types.NewSlice(types.NewPointer(user)), wantNil: true},
		{name: "pointer to slice of named struct", typ: types.NewPointer(types.NewSlice(user)), wantNil: false},
		{name: "map string to named struct", typ: types.NewMap(types.Typ[types.String], user), wantNil: false},
		{
			name:    "map string to pointer",
			typ:     types.NewMap(types.Typ[types.String], types.NewPointer(user)),
			wantNil: true,
		},
		{
			name:    "pointer to map string to named struct",
			typ:     types.NewPointer(types.NewMap(types.Typ[types.String], user)),
			wantNil: false,
		},
		{
			name:    "anonymous struct with fields",
			typ:     types.NewStruct([]*types.Var{types.NewVar(0, nil, "X", types.Typ[types.Int])}, nil),
			wantNil: false,
		},
		{name: "empty struct", typ: types.NewStruct(nil, nil), wantNil: true},
		{
			name: "pointer to anonymous struct",
			typ: types.NewPointer(
				types.NewStruct([]*types.Var{types.NewVar(0, nil, "X", types.Typ[types.Int])}, nil),
			),
			wantNil: true,
		},
		{name: "pointer to basic type", typ: types.NewPointer(types.Typ[types.Int]), wantNil: true},
		{
			name:    "nested slice of map of struct",
			typ:     types.NewSlice(types.NewMap(types.Typ[types.String], user)),
			wantNil: false,
		},
		{
			name: "empty named struct",
			typ:  newNamedStruct("example.com/app", "app", "Empty"),
			// empty struct has 0 fields but it IS a named struct — we check the
			// underlying struct field count to decide.  An empty named struct is
			// still treated as a value type (e.g., sentinel, token) so we exempt it.
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := FindViolation(tt.typ, cfg, nil)
			if tt.wantNil && v != nil {
				t.Errorf("expected no violation, got %+v", v)
			}
			if !tt.wantNil && v == nil {
				t.Error("expected violation, got nil")
			}
		})
	}
}

func TestFindViolation_MapKeyDefaultOff(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig()
	// MapKey is false by default
	m := types.NewMap(user, types.NewPointer(user))
	v := FindViolation(m, cfg, nil)
	if v != nil {
		t.Errorf("map key check should be off by default, got %+v", v)
	}
}

func TestFindViolation_MapKeyEnabled(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig()
	cfg.MapKey = true
	m := types.NewMap(user, types.NewPointer(user))
	v := FindViolation(m, cfg, nil)
	if v == nil {
		t.Error("expected violation for map key with MapKey enabled")
	}
}

func TestFindViolation_RecursiveTypes(t *testing.T) {
	t.Parallel()

	t.Run("recursive slice type A = []A", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig()
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "A", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewSlice(named))

		v := FindViolation(named, cfg, nil)
		if v != nil {
			t.Errorf("expected nil for recursive slice type, got %+v", v)
		}
	})

	t.Run("recursive map type M = map[string]M", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig()
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "M", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewMap(types.Typ[types.String], named))

		v := FindViolation(named, cfg, nil)
		if v != nil {
			t.Errorf("expected nil for recursive map type, got %+v", v)
		}
	})

	t.Run("recursive chan type C = chan C", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig()
		cfg.ChanElem = true
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "C", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewChan(types.SendRecv, named))

		v := FindViolation(named, cfg, nil)
		if v != nil {
			t.Errorf("expected nil for recursive chan type, got %+v", v)
		}
	})
}
