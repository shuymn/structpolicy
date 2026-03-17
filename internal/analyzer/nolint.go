package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// fileSuppression caches file-level nolint status so it is computed once per file.
type fileSuppression struct {
	checked    bool
	suppressed bool
}

// nolintMatcher bundles the parameters needed to match nolint directives.
type nolintMatcher struct {
	honorAll   bool
	linterName string
}

func newNolintMatcher(cfg *Config) nolintMatcher {
	return nolintMatcher{
		honorAll:   cfg.HonorNolintAll,
		linterName: cfg.Mode.LinterName(),
	}
}

// match reports whether a comment text matches //nolint:<linterName> or
// //nolint:all (when honorAll is true).
func (m nolintMatcher) match(text string) bool {
	// Strip leading "//" and whitespace.
	s := strings.TrimPrefix(text, "//")
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, "nolint:") {
		return false
	}

	// Extract linter list: "nolint:a,b // reason" -> "a,b"
	s = strings.TrimPrefix(s, "nolint:")
	if idx := strings.Index(s, "//"); idx >= 0 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)

	for name := range strings.SplitSeq(s, ",") {
		name = strings.TrimSpace(name)
		if name == m.linterName {
			return true
		}
		if m.honorAll && name == "all" {
			return true
		}
	}
	return false
}

// isSuppressed reports whether the diagnostic at pos is suppressed by a nolint
// comment. It checks inline (same line), block (doc comment on decl), and
// file-level (comment before package clause) comments.
func isSuppressed(
	pass *analysis.Pass,
	pos token.Pos,
	decl ast.Node,
	file *ast.File,
	cfg *Config,
	fileSupp *fileSuppression,
) bool {
	if !cfg.HonorNolint {
		return false
	}

	m := newNolintMatcher(cfg)

	if fileSuppressed(file, pass.Fset, m, fileSupp) {
		return true
	}

	line := pass.Fset.Position(pos).Line

	return isSuppressedInline(file, pass.Fset, line, m) ||
		isSuppressedBlock(decl, m)
}

// fileSuppressed checks the cached file-level suppression, computing it once.
func fileSuppressed(file *ast.File, fset *token.FileSet, m nolintMatcher, fs *fileSuppression) bool {
	if !fs.checked {
		fs.checked = true
		fs.suppressed = checkFileSuppression(file, fset, m)
	}
	return fs.suppressed
}

// isSuppressedInline checks whether any comment on the same line contains a
// matching nolint directive. Comments are position-sorted, so we break early
// once past the target line.
func isSuppressedInline(file *ast.File, fset *token.FileSet, line int, m nolintMatcher) bool {
	for _, cg := range file.Comments {
		cgLine := fset.Position(cg.Pos()).Line
		if cgLine > line {
			break
		}
		for _, c := range cg.List {
			if fset.Position(c.Pos()).Line == line && m.match(c.Text) {
				return true
			}
		}
	}
	return false
}

// isSuppressedBlock checks whether the doc comment group attached to decl
// contains a nolint directive.
func isSuppressedBlock(decl ast.Node, m nolintMatcher) bool {
	var doc *ast.CommentGroup
	switch d := decl.(type) {
	case *ast.FuncDecl:
		doc = d.Doc
	case *ast.GenDecl:
		doc = d.Doc
	case *ast.TypeSpec:
		doc = d.Doc
	}
	if doc == nil {
		return false
	}
	for _, c := range doc.List {
		if m.match(c.Text) {
			return true
		}
	}
	return false
}

// checkFileSuppression checks whether a comment before the package clause
// contains a nolint directive.
func checkFileSuppression(file *ast.File, fset *token.FileSet, m nolintMatcher) bool {
	pkgLine := fset.Position(file.Package).Line
	for _, cg := range file.Comments {
		if fset.Position(cg.End()).Line >= pkgLine {
			break
		}
		for _, c := range cg.List {
			if m.match(c.Text) {
				return true
			}
		}
	}
	return false
}
