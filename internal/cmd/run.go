package cmd

import (
	"fmt"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

// exitDiagnostics indicates that analysis succeeded but reported diagnostics.
// This follows the convention established by singlechecker/multichecker.
const exitDiagnostics = 3

// Run executes the analyzer as a standalone command-line tool.
// It parses flags from os.Args, loads packages, runs the analysis,
// and calls os.Exit with the appropriate exit code.
func Run(a *analysis.Analyzer) {
	a.Flags.Usage = func() { printUsage(a) }

	if err := a.Flags.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	patterns := a.Flags.Args()
	if len(patterns) == 0 {
		patterns = []string{"."}
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: true,
	}, patterns...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	graph, err := checker.Analyze([]*analysis.Analyzer{a}, pkgs, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := graph.PrintText(os.Stderr, -1); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitCode := 0
	for _, root := range graph.Roots {
		if root.Err != nil {
			exitCode = 1
			continue
		}
		if len(root.Diagnostics) > 0 {
			exitCode = exitDiagnostics
		}
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func printUsage(a *analysis.Analyzer) {
	fmt.Fprintf(os.Stderr, "%s: %s\n\n", a.Name, a.Doc)
	fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n", a.Name)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	a.Flags.PrintDefaults()
}
