package ptrstruct

import (
	"golang.org/x/tools/go/analysis"

	"github.com/shuymn/structpolicy/internal/analyzer"
)

// Analyzer reports struct types used by value where a pointer is expected.
var Analyzer = NewAnalyzer()

// NewAnalyzer creates a new ptrstruct analyzer with default configuration.
// Each call returns an independent analyzer with its own Config, safe for
// concurrent use in tests with different flag settings.
func NewAnalyzer() *analysis.Analyzer {
	return analyzer.NewAnalyzer(analyzer.ModePointer)
}
