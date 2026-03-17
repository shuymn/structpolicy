package analyzer

import (
	"go/types"
	"testing"
)

func TestFindViolation(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig(ModePointer)
	cfg.Param = true
	cfg.Result = true
	cfg.Field = true
	cfg.SliceElem = true
	cfg.MapValue = true

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
			name:    "empty named struct",
			typ:     newNamedStruct("example.com/app", "app", "Empty"),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, ok := FindViolation(tt.typ, cfg, nil)
			if tt.wantNil && ok {
				t.Errorf("expected no violation, got %+v", v)
			}
			if !tt.wantNil && !ok {
				t.Error("expected violation, got none")
			}
		})
	}
}

func TestFindViolation_ModeValue(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig(ModeValue)
	cfg.Param = true
	cfg.Result = true
	cfg.Field = true
	cfg.SliceElem = true
	cfg.MapValue = true

	tests := []struct {
		name    string
		typ     types.Type
		wantNil bool
	}{
		{name: "basic type int", typ: types.Typ[types.Int], wantNil: true},
		{name: "basic type string", typ: types.Typ[types.String], wantNil: true},
		{name: "interface{}", typ: types.NewInterfaceType(nil, nil), wantNil: true},
		{name: "named struct by value", typ: user, wantNil: true},
		{name: "pointer to named struct", typ: types.NewPointer(user), wantNil: false},
		{name: "slice of named struct", typ: types.NewSlice(user), wantNil: true},
		{name: "slice of pointer to named struct", typ: types.NewSlice(types.NewPointer(user)), wantNil: false},
		{name: "pointer to slice of named struct", typ: types.NewPointer(types.NewSlice(user)), wantNil: true},
		{name: "map string to named struct", typ: types.NewMap(types.Typ[types.String], user), wantNil: true},
		{
			name:    "map string to pointer to named struct",
			typ:     types.NewMap(types.Typ[types.String], types.NewPointer(user)),
			wantNil: false,
		},
		{
			name:    "pointer to map string to named struct",
			typ:     types.NewPointer(types.NewMap(types.Typ[types.String], user)),
			wantNil: true,
		},
		{
			name:    "anonymous struct with fields",
			typ:     types.NewStruct([]*types.Var{types.NewVar(0, nil, "X", types.Typ[types.Int])}, nil),
			wantNil: true,
		},
		{name: "empty struct", typ: types.NewStruct(nil, nil), wantNil: true},
		{
			name: "pointer to anonymous struct",
			typ: types.NewPointer(
				types.NewStruct([]*types.Var{types.NewVar(0, nil, "X", types.Typ[types.Int])}, nil),
			),
			wantNil: false,
		},
		{name: "pointer to basic type", typ: types.NewPointer(types.Typ[types.Int]), wantNil: true},
		{
			name:    "empty named struct",
			typ:     newNamedStruct("example.com/app", "app", "Empty"),
			wantNil: true,
		},
		{
			name:    "pointer to empty named struct",
			typ:     types.NewPointer(newNamedStruct("example.com/app", "app", "Empty")),
			wantNil: true,
		},
		{
			name:    "pointer to pointer to named struct",
			typ:     types.NewPointer(types.NewPointer(user)),
			wantNil: false,
		},
		{
			name:    "pointer to slice of pointer to named struct",
			typ:     types.NewPointer(types.NewSlice(types.NewPointer(user))),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, ok := FindViolation(tt.typ, cfg, nil)
			if tt.wantNil && ok {
				t.Errorf("expected no violation, got %+v", v)
			}
			if !tt.wantNil && !ok {
				t.Error("expected violation, got none")
			}
		})
	}
}

func TestFindViolation_ModeValue_ArrayAndChan(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig(ModeValue)
	cfg.ArrayElem = true
	cfg.ChanElem = true

	tests := []struct {
		name    string
		typ     types.Type
		wantNil bool
	}{
		{
			name:    "array of pointer to named struct",
			typ:     types.NewArray(types.NewPointer(user), 3),
			wantNil: false,
		},
		{
			name:    "array of named struct",
			typ:     types.NewArray(user, 3),
			wantNil: true,
		},
		{
			name:    "chan of pointer to named struct",
			typ:     types.NewChan(types.SendRecv, types.NewPointer(user)),
			wantNil: false,
		},
		{
			name:    "chan of named struct",
			typ:     types.NewChan(types.SendRecv, user),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, ok := FindViolation(tt.typ, cfg, nil)
			if tt.wantNil && ok {
				t.Errorf("expected no violation, got %+v", v)
			}
			if !tt.wantNil && !ok {
				t.Error("expected violation, got none")
			}
		})
	}
}

func TestFindViolation_MapKeyDefaultOff(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig(ModePointer)
	// MapKey is false by default
	m := types.NewMap(user, types.NewPointer(user))
	v, ok := FindViolation(m, cfg, nil)
	if ok {
		t.Errorf("map key check should be off by default, got %+v", v)
	}
}

func TestFindViolation_MapKeyEnabled(t *testing.T) {
	t.Parallel()

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	cfg := DefaultConfig(ModePointer)
	cfg.MapKey = true
	m := types.NewMap(user, types.NewPointer(user))
	_, ok := FindViolation(m, cfg, nil)
	if !ok {
		t.Error("expected violation for map key with MapKey enabled")
	}
}

func TestFindViolation_RecursiveTypes(t *testing.T) {
	t.Parallel()

	t.Run("recursive slice type A = []A", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModePointer)
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "A", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewSlice(named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive slice type, got %+v", v)
		}
	})

	t.Run("recursive map type M = map[string]M", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModePointer)
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "M", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewMap(types.Typ[types.String], named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive map type, got %+v", v)
		}
	})

	t.Run("recursive chan type C = chan C", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModePointer)
		cfg.ChanElem = true
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "C", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewChan(types.SendRecv, named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive chan type, got %+v", v)
		}
	})
}

func TestFindViolation_RecursiveTypes_ModeValue(t *testing.T) {
	t.Parallel()

	t.Run("recursive slice type A = []A", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModeValue)
		cfg.SliceElem = true
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "A", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewSlice(named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive slice type, got %+v", v)
		}
	})

	t.Run("recursive map type M = map[string]M", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModeValue)
		cfg.MapValue = true
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "M", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewMap(types.Typ[types.String], named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive map type, got %+v", v)
		}
	})

	t.Run("recursive chan type C = chan C", func(t *testing.T) {
		t.Parallel()
		cfg := DefaultConfig(ModeValue)
		cfg.ChanElem = true
		pkg := types.NewPackage("example.com/foo", "foo")
		tn := types.NewTypeName(0, pkg, "C", nil)
		named := types.NewNamed(tn, nil, nil)
		named.SetUnderlying(types.NewChan(types.SendRecv, named))

		v, ok := FindViolation(named, cfg, nil)
		if ok {
			t.Errorf("expected no violation for recursive chan type, got %+v", v)
		}
	})
}
