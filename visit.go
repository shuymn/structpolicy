package ptrstruct

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type funcDeclChecker func(*analysis.Pass, *ast.FuncDecl, *Config, *Classifier) *analysis.Diagnostic

func visitFuncDecl(
	pass *analysis.Pass,
	file *ast.File,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
	fileSupp *fileSuppression,
) {
	checks := []funcDeclChecker{checkReceiver, checkParams, checkResults}
	for _, check := range checks {
		diag := check(pass, decl, cfg, cls)
		if diag == nil {
			continue
		}
		if !isSuppressed(pass, diag.Pos, decl, file, cfg, fileSupp) {
			pass.Report(*diag)
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
	if !cfg.Field {
		return
	}

	obj := pass.TypesInfo.Defs[spec.Name]
	if obj == nil {
		return
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return
	}

	for i := range st.NumFields() {
		field := st.Field(i)
		v := FindViolation(field.Type(), cfg, cls)
		if v == nil {
			continue
		}

		msg := FormatDiagnostic("field "+field.Name(), v)
		pos := spec.Pos()
		if astField := structField(spec, i); astField != nil {
			pos = astField.Pos()
		}
		diag := analysis.Diagnostic{Pos: pos, Message: msg}

		declNode := blockOrSpec(genDecl, spec)
		if !isSuppressed(pass, diag.Pos, declNode, file, cfg, fileSupp) {
			pass.Report(diag)
		}
		return // 1 violation per declaration
	}
}

func checkReceiver(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) *analysis.Diagnostic {
	if !cfg.Receiver || decl.Recv == nil || len(decl.Recv.List) == 0 {
		return nil
	}

	field := decl.Recv.List[0]
	t := pass.TypesInfo.TypeOf(field.Type)
	if t == nil {
		return nil
	}

	v := FindViolation(t, cfg, cls)
	if v == nil {
		return nil
	}

	msg := FormatDiagnostic("receiver", v)
	return &analysis.Diagnostic{Pos: field.Pos(), Message: msg}
}

func checkParams(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) *analysis.Diagnostic {
	if !cfg.Param || decl.Type.Params == nil {
		return nil
	}
	return checkFieldList(pass, decl.Type.Params, cfg, cls, fieldLabeler)
}

func checkResults(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	cfg *Config,
	cls *Classifier,
) *analysis.Diagnostic {
	if !cfg.Result || decl.Type.Results == nil {
		return nil
	}
	return checkFieldList(pass, decl.Type.Results, cfg, cls, func(*ast.Field) string { return "result" })
}

func checkFieldList(
	pass *analysis.Pass,
	fields *ast.FieldList,
	cfg *Config,
	cls *Classifier,
	label func(*ast.Field) string,
) *analysis.Diagnostic {
	for _, field := range fields.List {
		t := pass.TypesInfo.TypeOf(field.Type)
		if t == nil {
			continue
		}
		v := FindViolation(t, cfg, cls)
		if v == nil {
			continue
		}
		msg := FormatDiagnostic(label(field), v)
		return &analysis.Diagnostic{Pos: field.Pos(), Message: msg}
	}
	return nil
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

// blockOrSpec returns the suppression target node for nolint checking.
// For standalone type declarations (single spec), the GenDecl carries the doc.
// For grouped type blocks (multiple specs), each TypeSpec has its own doc.
func blockOrSpec(genDecl *ast.GenDecl, spec *ast.TypeSpec) ast.Node {
	if len(genDecl.Specs) > 1 {
		return spec
	}
	return genDecl
}
