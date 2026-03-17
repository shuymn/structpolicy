package analyzer

import (
	"fmt"
	"go/build"
	"go/types"
	"regexp"
	"strings"
)

// classifier checks whether a named type is exempted by the allowlist.
type classifier struct {
	allowStdlib     bool
	allowThirdParty bool
	modulePath      string
	types           map[string]bool
	patterns        []*regexp.Regexp
	packages        map[string]bool
	stdlib          map[string]bool
}

func newClassifier(cfg *config, modulePath string) (*classifier, error) {
	c := &classifier{
		allowStdlib:     cfg.allowStdlib,
		allowThirdParty: cfg.allowThirdParty,
		modulePath:      modulePath,
		types:           make(map[string]bool, len(cfg.allowTypes)),
		packages:        make(map[string]bool, len(cfg.allowPackages)),
		stdlib:          make(map[string]bool),
	}
	for _, t := range cfg.allowTypes {
		c.types[t] = true
	}
	for _, p := range cfg.allowPackages {
		c.packages[p] = true
	}
	for _, pat := range cfg.allowPatterns {
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

// isAllowed reports whether the named type is exempted by any allowlist entry.
func (c *classifier) isAllowed(named *types.Named) bool {
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

func (c *classifier) isStdlib(pkg *types.Package) bool {
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

func (c *classifier) isThirdParty(pkg *types.Package) bool {
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
