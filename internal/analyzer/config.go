package analyzer

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Config holds the analyzer configuration.
type Config struct {
	Mode            Mode
	Receiver        bool
	Param           bool
	Result          bool
	Field           bool
	InterfaceMethod bool
	FuncType        bool
	NamedType       bool
	SliceElem       bool
	MapValue        bool
	MapKey          bool
	ArrayElem       bool
	ChanElem        bool
	IgnoreGenerated bool
	IgnoreTests     bool
	HonorNolint     bool
	HonorNolintAll  bool
	AllowStdlib     bool
	AllowThirdParty bool
	AllowTypes      []string
	AllowPatterns   []string
	AllowPackages   []string
}

// DefaultConfig returns a Config with mode-specific performance-tuning defaults.
func DefaultConfig(mode Mode) *Config {
	cfg := &Config{
		Mode:            mode,
		Receiver:        true,
		Param:           false,
		Result:          false,
		Field:           false,
		InterfaceMethod: false,
		FuncType:        false,
		NamedType:       false,
		SliceElem:       false,
		MapValue:        false,
		MapKey:          false,
		ArrayElem:       false,
		ChanElem:        false,
		IgnoreGenerated: true,
		IgnoreTests:     false,
		HonorNolint:     true,
		HonorNolintAll:  true,
		AllowStdlib:     true,
		AllowThirdParty: false,
		AllowTypes:      nil,
		AllowPatterns:   nil,
		AllowPackages:   nil,
	}

	switch mode {
	case ModePointer:
		cfg.Param = true
		cfg.Field = true
		cfg.SliceElem = true
		cfg.MapValue = true
		cfg.ArrayElem = true
		cfg.ChanElem = true
	case ModeValue:
		cfg.Receiver = false
		cfg.Param = true
		cfg.Result = true
		cfg.Field = true
		cfg.SliceElem = true
		cfg.MapValue = true
		cfg.ArrayElem = true
		cfg.ChanElem = true
	}

	return cfg
}

func registerFlags(a *analysis.Analyzer, cfg *Config) {
	a.Flags.BoolVar(&cfg.Receiver, "receiver", cfg.Receiver, "check method receivers")
	a.Flags.BoolVar(&cfg.Param, "param", cfg.Param, "check function parameters")
	a.Flags.BoolVar(&cfg.Result, "result", cfg.Result, "check function results")
	a.Flags.BoolVar(&cfg.Field, "field", cfg.Field, "check struct fields")
	a.Flags.BoolVar(&cfg.InterfaceMethod, "interface-method", cfg.InterfaceMethod, "check interface methods")
	a.Flags.BoolVar(&cfg.FuncType, "func-type", cfg.FuncType, "check function types")
	a.Flags.BoolVar(&cfg.NamedType, "named-type", cfg.NamedType, "check named container types")
	a.Flags.BoolVar(&cfg.SliceElem, "slice-elem", cfg.SliceElem, "check slice element types")
	a.Flags.BoolVar(&cfg.MapValue, "map-value", cfg.MapValue, "check map value types")
	a.Flags.BoolVar(&cfg.MapKey, "map-key", cfg.MapKey, "check map key types")
	a.Flags.BoolVar(&cfg.ArrayElem, "array-elem", cfg.ArrayElem, "check array element types")
	a.Flags.BoolVar(&cfg.ChanElem, "chan-elem", cfg.ChanElem, "check channel element types")
	a.Flags.BoolVar(&cfg.IgnoreGenerated, "ignore-generated", cfg.IgnoreGenerated, "skip generated files")
	a.Flags.BoolVar(&cfg.IgnoreTests, "ignore-tests", cfg.IgnoreTests, "skip test files")
	a.Flags.BoolVar(
		&cfg.HonorNolint,
		"honor-nolint",
		cfg.HonorNolint,
		"honor //nolint:"+cfg.Mode.LinterName()+" comments",
	)
	a.Flags.BoolVar(&cfg.HonorNolintAll, "honor-nolint-all", cfg.HonorNolintAll, "honor //nolint:all comments")
	a.Flags.BoolVar(&cfg.AllowStdlib, "allow-stdlib", cfg.AllowStdlib, "exempt builtin and standard library packages")
	a.Flags.BoolVar(
		&cfg.AllowThirdParty,
		"allow-third-party",
		cfg.AllowThirdParty,
		"exempt non-stdlib packages outside the current Go module",
	)
	a.Flags.Var(&stringListFlag{values: &cfg.AllowTypes}, "allow-types", "comma-separated allowed type names")
	a.Flags.Var(&stringListFlag{values: &cfg.AllowPatterns}, "allow-patterns", "comma-separated allowed type patterns")
	a.Flags.Var(&stringListFlag{values: &cfg.AllowPackages}, "allow-packages", "comma-separated allowed package paths")
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
