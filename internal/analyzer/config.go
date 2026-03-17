package analyzer

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// config holds the analyzer configuration.
type config struct {
	mode            mode
	receiver        bool
	param           bool
	result          bool
	field           bool
	interfaceMethod bool
	funcType        bool
	namedType       bool
	sliceElem       bool
	mapValue        bool
	mapKey          bool
	arrayElem       bool
	chanElem        bool
	ignoreGenerated bool
	ignoreTests     bool
	honorNolint     bool
	honorNolintAll  bool
	allowStdlib     bool
	allowThirdParty bool
	allowTypes      []string
	allowPatterns   []string
	allowPackages   []string
}

// defaultConfig returns a config with opposite mode-specific performance-tuning defaults.
func defaultConfig(mode mode) *config {
	cfg := &config{
		mode:            mode,
		receiver:        true,
		param:           false,
		result:          false,
		field:           false,
		interfaceMethod: false,
		funcType:        false,
		namedType:       false,
		sliceElem:       false,
		mapValue:        false,
		mapKey:          false,
		arrayElem:       false,
		chanElem:        false,
		ignoreGenerated: true,
		ignoreTests:     false,
		honorNolint:     true,
		honorNolintAll:  true,
		allowStdlib:     true,
		allowThirdParty: false,
		allowTypes:      nil,
		allowPatterns:   nil,
		allowPackages:   nil,
	}

	switch mode {
	case ModePointer:
		cfg.param = true
		cfg.field = true
		cfg.interfaceMethod = true
		cfg.funcType = true
	case ModeValue:
		cfg.receiver = false
		cfg.result = true
		cfg.namedType = true
		cfg.sliceElem = true
		cfg.mapValue = true
		cfg.mapKey = true
		cfg.arrayElem = true
		cfg.chanElem = true
	}

	return cfg
}

func registerFlags(a *analysis.Analyzer, cfg *config) {
	a.Flags.BoolVar(&cfg.receiver, "receiver", cfg.receiver, "check method receivers")
	a.Flags.BoolVar(&cfg.param, "param", cfg.param, "check function parameters")
	a.Flags.BoolVar(&cfg.result, "result", cfg.result, "check function results")
	a.Flags.BoolVar(&cfg.field, "field", cfg.field, "check struct fields")
	a.Flags.BoolVar(&cfg.interfaceMethod, "interface-method", cfg.interfaceMethod, "check interface methods")
	a.Flags.BoolVar(&cfg.funcType, "func-type", cfg.funcType, "check function types")
	a.Flags.BoolVar(&cfg.namedType, "named-type", cfg.namedType, "check named container types")
	a.Flags.BoolVar(&cfg.sliceElem, "slice-elem", cfg.sliceElem, "check slice element types")
	a.Flags.BoolVar(&cfg.mapValue, "map-value", cfg.mapValue, "check map value types")
	a.Flags.BoolVar(&cfg.mapKey, "map-key", cfg.mapKey, "check map key types")
	a.Flags.BoolVar(&cfg.arrayElem, "array-elem", cfg.arrayElem, "check array element types")
	a.Flags.BoolVar(&cfg.chanElem, "chan-elem", cfg.chanElem, "check channel element types")
	a.Flags.BoolVar(&cfg.ignoreGenerated, "ignore-generated", cfg.ignoreGenerated, "skip generated files")
	a.Flags.BoolVar(&cfg.ignoreTests, "ignore-tests", cfg.ignoreTests, "skip test files")
	a.Flags.BoolVar(
		&cfg.honorNolint,
		"honor-nolint",
		cfg.honorNolint,
		"honor //nolint:"+cfg.mode.linterName()+" comments",
	)
	a.Flags.BoolVar(&cfg.honorNolintAll, "honor-nolint-all", cfg.honorNolintAll, "honor //nolint:all comments")
	a.Flags.BoolVar(&cfg.allowStdlib, "allow-stdlib", cfg.allowStdlib, "exempt builtin and standard library packages")
	a.Flags.BoolVar(
		&cfg.allowThirdParty,
		"allow-third-party",
		cfg.allowThirdParty,
		"exempt non-stdlib packages outside the current Go module",
	)
	a.Flags.Var(&stringListFlag{values: &cfg.allowTypes}, "allow-types", "comma-separated allowed type names")
	a.Flags.Var(&stringListFlag{values: &cfg.allowPatterns}, "allow-patterns", "comma-separated allowed type patterns")
	a.Flags.Var(&stringListFlag{values: &cfg.allowPackages}, "allow-packages", "comma-separated allowed package paths")
}

// stringListFlag implements flag.Value for comma-separated string lists.
type stringListFlag struct {
	values *[]string
}

func (f *stringListFlag) String() string {
	if f.values == nil || *f.values == nil {
		return ""
	}
	return strings.Join(*f.values, ",")
}

func (f *stringListFlag) Set(s string) error {
	if s == "" {
		*f.values = nil
		return nil
	}
	parts := strings.Split(s, ",")
	n := 0
	for _, p := range parts {
		if p != "" {
			parts[n] = p
			n++
		}
	}
	*f.values = parts[:n]
	return nil
}
