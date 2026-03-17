package analyzer

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer creates a new analyzer with default configuration for the given mode.
// Each call returns an independent analyzer with its own Config, safe for
// concurrent use in tests with different flag settings.
func NewAnalyzer(mode Mode) *analysis.Analyzer {
	cfg := DefaultConfig(mode)
	a := &analysis.Analyzer{
		Name: mode.LinterName(),
		Doc:  mode.Doc(),
		Run:  func(pass *analysis.Pass) (any, error) { return run(pass, cfg) },
	}
	registerFlags(a, cfg)
	return a
}

func run(pass *analysis.Pass, cfg *Config) (any, error) {
	modulePath := ""
	if cfg.AllowThirdParty {
		modulePath = modulePathForPass(pass)
	}

	cls, err := newClassifier(cfg, modulePath)
	if err != nil {
		return nil, err
	}

	for _, file := range pass.Files {
		if cfg.IgnoreGenerated && ast.IsGenerated(file) {
			continue
		}
		if cfg.IgnoreTests && isTestFile(pass, file) {
			continue
		}
		visitFile(pass, file, cfg, cls)
	}
	return nil, nil
}

func visitFile(pass *analysis.Pass, file *ast.File, cfg *Config, cls *Classifier) {
	var fileSupp fileSuppression
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			visitFuncDecl(pass, file, d, cfg, cls, &fileSupp)
		case *ast.GenDecl:
			visitGenDecl(pass, file, d, cfg, cls, &fileSupp)
		}
	}
}

func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	name := pass.Fset.Position(file.Package).Filename
	return strings.HasSuffix(name, "_test.go")
}

func modulePathForPass(pass *analysis.Pass) string {
	if pass.Module != nil && pass.Module.Path != "" {
		return pass.Module.Path
	}
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Package).Filename
		if filename == "" {
			continue
		}
		if modulePath := modulePathForFile(filename); modulePath != "" {
			return modulePath
		}
	}
	return ""
}

func modulePathForFile(filename string) string {
	dir := filepath.Dir(filename)
	for {
		gomod := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(gomod)
		if err == nil {
			return modfile.ModulePath(data)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
