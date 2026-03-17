package ptrstruct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/shuymn/ptrstruct"
)

// allChecksAnalyzer returns an analyzer with all declaration and container
// checks enabled, matching the old Phase-1 defaults.
func allChecksAnalyzer(t *testing.T) *analysis.Analyzer {
	t.Helper()
	a := ptrstruct.NewAnalyzer()
	for _, flag := range []string{"param", "result", "field", "slice-elem", "map-value"} {
		if err := a.Flags.Set(flag, "true"); err != nil {
			t.Fatal(err)
		}
	}
	return a
}

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "basic", "ok", "containers", "generics")
}

func TestAnalyzer_Suppress(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "suppress")
}

func TestAnalyzer_FileNolint(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "filenolint")
}

func TestAnalyzer_TypeBlock(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "typeblock")
}

func TestAnalyzer_Allow(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	for _, flag := range []string{"param", "result", "field", "slice-elem", "map-value"} {
		if err := a.Flags.Set(flag, "true"); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.Flags.Set("allow-types", "time.Time"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "allow")
}

func TestAnalyzer_AllowStdlib(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	for _, flag := range []string{"param", "result", "field", "slice-elem", "map-value"} {
		if err := a.Flags.Set(flag, "true"); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.Flags.Set("allow-stdlib", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "allowstdlib")
}

func TestAnalyzer_AllowThirdParty(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	for _, flag := range []string{"param", "result", "field", "slice-elem", "map-value"} {
		if err := a.Flags.Set(flag, "true"); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.Flags.Set("allow-stdlib", "false"); err != nil {
		t.Fatal(err)
	}
	if err := a.Flags.Set("allow-third-party", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "allowthirdparty")
}

func TestAnalyzer_Alias(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "alias")
}

func TestAnalyzer_Nested(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "nested")
}

func TestAnalyzer_OnePer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "oneper")
}

func TestAnalyzer_Generated(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "generated")
}

func TestAnalyzer_Embedded(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "embedded")
}

func TestAnalyzer_IgnoreTestsDefault(t *testing.T) {
	t.Parallel()

	// Default: IgnoreTests=false, so test files ARE checked.
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allChecksAnalyzer(t), "ignoretests")
}

func TestAnalyzer_IgnoreTestsEnabled(t *testing.T) {
	t.Parallel()

	// With IgnoreTests=true, _test.go files should be skipped.
	// skiptests/_test.go has a violation but no // want comment;
	// if the analyzer checked it, analysistest would fail.
	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	for _, flag := range []string{"param", "result", "field", "slice-elem", "map-value"} {
		if err := a.Flags.Set(flag, "true"); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.Flags.Set("ignore-tests", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "skiptests")
}
