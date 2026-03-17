package analyzer

import (
	"fmt"
	"go/build"
	"go/types"
	"regexp"
	"strings"
)

// Classifier checks whether a named type is exempted by the allowlist.
type Classifier struct {
	allowStdlib     bool
	allowThirdParty bool
	modulePath      string
	types           map[string]bool
	patterns        []*regexp.Regexp
	packages        map[string]bool
	stdlib          map[string]bool
}

// NewClassifier creates a Classifier from the allowlist fields of cfg.
// It returns an error if any AllowPatterns entry is not a valid regexp.
func NewClassifier(cfg *Config) (*Classifier, error) {
	return newClassifier(cfg, "")
}

func newClassifier(cfg *Config, modulePath string) (*Classifier, error) {
	c := &Classifier{
		allowStdlib:     cfg.AllowStdlib,
		allowThirdParty: cfg.AllowThirdParty,
		modulePath:      modulePath,
		types:           make(map[string]bool, len(cfg.AllowTypes)),
		packages:        make(map[string]bool, len(cfg.AllowPackages)),
		stdlib:          make(map[string]bool),
	}
	for _, t := range cfg.AllowTypes {
		c.types[t] = true
	}
	for _, p := range cfg.AllowPackages {
		c.packages[p] = true
	}
	for _, pat := range cfg.AllowPatterns {
		if pat == "" {
			continue
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("invalid allow-pattern %q: %w", pat, err)
		}
		c.patterns = append(c.patterns, re)
	}
	return c, nil
}

// IsAllowed reports whether the named type is exempted by any allowlist entry.
func (c *Classifier) IsAllowed(named *types.Named) bool {
	pkg := named.Obj().Pkg()

	if c.allowStdlib && c.isStdlib(pkg) {
		return true
	}
	if c.allowThirdParty && c.isThirdParty(pkg) {
		return true
	}

	fqn := qualifiedName(named)

	if c.types[fqn] {
		return true
	}

	if pkg != nil && c.packages[pkg.Path()] {
		return true
	}

	for _, re := range c.patterns {
		if re.MatchString(fqn) {
			return true
		}
	}

	return false
}

func (c *Classifier) isStdlib(pkg *types.Package) bool {
	if pkg == nil {
		return true
	}

	path := pkg.Path()
	if ok, found := c.stdlib[path]; found {
		return ok
	}

	info, err := build.Default.Import(path, "", build.FindOnly)
	ok := err == nil && info.Goroot
	c.stdlib[path] = ok
	return ok
}

func (c *Classifier) isThirdParty(pkg *types.Package) bool {
	if pkg == nil || c.modulePath == "" {
		return false
	}
	if c.isStdlib(pkg) {
		return false
	}

	path := pkg.Path()
	return path != c.modulePath && !strings.HasPrefix(path, c.modulePath+"/")
}

func qualifiedName(named *types.Named) string {
	obj := named.Obj()
	if obj.Pkg() == nil {
		return obj.Name()
	}
	return obj.Pkg().Path() + "." + obj.Name()
}
