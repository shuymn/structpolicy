package allowthirdparty

import "thirdpartydep"

type Config struct {
	Name string
}

func ProcessDep(d thirdpartydep.Dep) {}

func Setup(c Config) {} // want `parameter c uses value struct Config; use \*Config`
