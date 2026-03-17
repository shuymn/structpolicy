package allow

import "time"

type Config struct {
	Name string
}

// time.Time is in the allowlist for this test: should not trigger.
func Process(t time.Time) {}

// Config is not in the allowlist: should trigger.
func Setup(c Config) {} // want `parameter c uses value struct Config; use \*Config`
