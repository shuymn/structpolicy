package analyzer

import "go/types"

func newNamedStruct(pkgPath, pkgName, typeName string, fields ...*types.Var) *types.Named {
	pkg := types.NewPackage(pkgPath, pkgName)
	st := types.NewStruct(fields, nil)
	obj := types.NewTypeName(0, pkg, typeName, nil)
	named := types.NewNamed(obj, st, nil)
	return named
}
