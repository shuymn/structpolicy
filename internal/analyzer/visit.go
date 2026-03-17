package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type funcDeclChecker func(*analysis.Pass, *ast.FuncDecl, *Config, *Classifier) (analysis.Diagnostic, bool)

var funcDeclChecks = [...]funcDeclChecker{checkReceiver, checkParams, checkResults}

type typeSpecChecker func(*analysis.Pass, *ast.TypeSpec, *Config, *Classifier) (analysis.Diagnostic, bool)

var typeSpecChecks = [...]typeSpecChecker{
	checkStructFields,
	checkInterfaceMethods,
	checkFuncType,
	checkNamedType,
}

func visitFuncDecl(
	pass *analysis.Pass,
	file *ast.File,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
	fileSupp *fileSuppression,
) {
	for _, check := range funcDeclChecks {
		diag, ok := check(pass, decl, cfg, cls)
		if !ok {
			continue
		}
		if !isSuppressed(pass, diag.Pos, decl, file, cfg, fileSupp) {
			pass.Report(diag)
		}
		return // 1 violation per declaration
	}
}

func visitGenDecl(
	pass *analysis.Pass,
	file *ast.File,
	genDecl *ast.GenDecl,
	cfg *Config,
	cls *Classifier,
	fileSupp *fileSuppression,
) {
	for _, spec := range genDecl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		visitTypeSpec(pass, file, ts, genDecl, cfg, cls, fileSupp)
	}
}

func visitTypeSpec(
	pass *analysis.Pass,
	file *ast.File,
	spec *ast.TypeSpec,
	genDecl *ast.GenDecl,
	cfg *Config,
	cls *Classifier,
	fileSupp *fileSuppression,
) {
	for _, check := range typeSpecChecks {
		diag, ok := check(pass, spec, cfg, cls)
		if !ok {
			continue
		}
		declNode := blockOrSpec(genDecl, spec)
		if !isSuppressed(pass, diag.Pos, declNode, file, cfg, fileSupp) {
			pass.Report(diag)
		}
		return // 1 violation per declaration
	}
}

func checkStructFields(
	pass *analysis.Pass,
	spec *ast.TypeSpec,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.Field {
		return analysis.Diagnostic{}, false
	}

	obj := pass.TypesInfo.Defs[spec.Name]
	if obj == nil {
		return analysis.Diagnostic{}, false
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return analysis.Diagnostic{}, false
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return analysis.Diagnostic{}, false
	}

	for i := range st.NumFields() {
		field := st.Field(i)
		v, ok := FindViolation(field.Type(), cfg, cls)
		if !ok {
			continue
		}

		msg := FormatDiagnostic("field "+field.Name(), &v, cfg.Mode)
		pos := spec.Pos()
		if astField := structField(spec, i); astField != nil {
			pos = astField.Pos()
		}
		return analysis.Diagnostic{Pos: pos, Message: msg}, true
	}

	return analysis.Diagnostic{}, false
}

func checkInterfaceMethods(
	pass *analysis.Pass,
	spec *ast.TypeSpec,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.InterfaceMethod {
		return analysis.Diagnostic{}, false
	}

	iface, ok := spec.Type.(*ast.InterfaceType)
	if !ok || iface.Methods == nil {
		return analysis.Diagnostic{}, false
	}

	for _, field := range iface.Methods.List {
		if len(field.Names) == 0 {
			continue // embedded interface
		}

		ft, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		if diag, ok := checkInterfaceMethod(pass, field.Names[0].Name, ft, cfg, cls); ok {
			return diag, true
		}
	}

	return analysis.Diagnostic{}, false
}

func checkInterfaceMethod(
	pass *analysis.Pass,
	methodName string,
	ft *ast.FuncType,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if cfg.Param && ft.Params != nil {
		if diag, ok := checkFieldList(
			pass,
			ft.Params,
			cfg,
			cls,
			func(f *ast.Field) string { return "interface method " + methodName + " " + fieldLabeler(f) },
		); ok {
			return diag, true
		}
	}
	if cfg.Result && ft.Results != nil {
		return checkFieldList(
			pass,
			ft.Results,
			cfg,
			cls,
			func(*ast.Field) string { return "interface method " + methodName + " result" },
		)
	}
	return analysis.Diagnostic{}, false
}

func checkFuncType(
	pass *analysis.Pass,
	spec *ast.TypeSpec,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.FuncType {
		return analysis.Diagnostic{}, false
	}

	ft, ok := spec.Type.(*ast.FuncType)
	if !ok {
		return analysis.Diagnostic{}, false
	}

	typeName := spec.Name.Name
	if cfg.Param && ft.Params != nil {
		if diag, ok := checkFieldList(
			pass,
			ft.Params,
			cfg,
			cls,
			func(f *ast.Field) string { return "function type " + typeName + " " + fieldLabeler(f) },
		); ok {
			return diag, true
		}
	}
	if cfg.Result && ft.Results != nil {
		return checkFieldList(
			pass,
			ft.Results,
			cfg,
			cls,
			func(*ast.Field) string { return "function type " + typeName + " result" },
		)
	}
	return analysis.Diagnostic{}, false
}

func checkNamedType(
	pass *analysis.Pass,
	spec *ast.TypeSpec,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.NamedType {
		return analysis.Diagnostic{}, false
	}

	obj := pass.TypesInfo.Defs[spec.Name]
	if obj == nil {
		return analysis.Diagnostic{}, false
	}

	named, ok := obj.Type().(*types.Named)
	if !ok || !isNamedContainerType(named.Underlying()) {
		return analysis.Diagnostic{}, false
	}

	v, ok := FindViolation(named, cfg, cls)
	if !ok {
		return analysis.Diagnostic{}, false
	}

	msg := FormatDiagnostic("named type "+spec.Name.Name, &v, cfg.Mode)
	return analysis.Diagnostic{Pos: spec.Pos(), Message: msg}, true
}

func checkReceiver(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.Receiver || decl.Recv == nil || len(decl.Recv.List) == 0 {
		return analysis.Diagnostic{}, false
	}

	field := decl.Recv.List[0]
	t := pass.TypesInfo.TypeOf(field.Type)
	if t == nil {
		return analysis.Diagnostic{}, false
	}

	v, ok := FindViolation(t, cfg, cls)
	if !ok {
		return analysis.Diagnostic{}, false
	}

	msg := FormatDiagnostic("receiver", &v, cfg.Mode)
	return analysis.Diagnostic{Pos: field.Pos(), Message: msg}, true
}

func checkParams(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.Param || decl.Type.Params == nil {
		return analysis.Diagnostic{}, false
	}
	return checkFieldList(pass, decl.Type.Params, cfg, cls, fieldLabeler)
}

func checkResults(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) (analysis.Diagnostic, bool) {
	if !cfg.Result || decl.Type.Results == nil {
		return analysis.Diagnostic{}, false
	}
	return checkFieldList(pass, decl.Type.Results, cfg, cls, func(*ast.Field) string { return "result" })
}

func checkFieldList(
	pass *analysis.Pass,
	fields *ast.FieldList,
	cfg *Config,
	cls *Classifier,
	label func(*ast.Field) string,
) (analysis.Diagnostic, bool) {
	for _, field := range fields.List {
		t := pass.TypesInfo.TypeOf(field.Type)
		if t == nil {
			continue
		}
		v, ok := FindViolation(t, cfg, cls)
		if !ok {
			continue
		}
		msg := FormatDiagnostic(label(field), &v, cfg.Mode)
		return analysis.Diagnostic{Pos: field.Pos(), Message: msg}, true
	}
	return analysis.Diagnostic{}, false
}

func fieldLabeler(field *ast.Field) string {
	if len(field.Names) > 0 {
		return "parameter " + field.Names[0].Name
	}
	return "parameter unnamed"
}

func structField(spec *ast.TypeSpec, i int) *ast.Field {
	st, ok := spec.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return nil
	}

	idx := 0
	for _, f := range st.Fields.List {
		n := len(f.Names)
		if n == 0 {
			n = 1 // embedded field
		}
		if i < idx+n {
			return f
		}
		idx += n
	}
	return nil
}

func isNamedContainerType(t types.Type) bool {
	switch tt := types.Unalias(t).(type) {
	case *types.Slice, *types.Map, *types.Array, *types.Chan:
		return true
	case *types.Pointer:
		return isNamedContainerType(tt.Elem())
	default:
		return false
	}
}

// blockOrSpec returns the suppression target node for nolint checking.
// For standalone type declarations (single spec), the GenDecl carries the doc.
// For grouped type blocks (multiple specs), each TypeSpec has its own doc.
func blockOrSpec(genDecl *ast.GenDecl, spec *ast.TypeSpec) ast.Node {
	if len(genDecl.Specs) > 1 {
		return spec
	}
	return genDecl
}
