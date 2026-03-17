package allowstdlib

import (
	"database/sql"
	"time"
)

type Config struct {
	Name string
}

func ProcessTime(t time.Time) {}

func ProcessNullString(ns sql.NullString) {}

func Setup(c Config) {} // want `parameter c uses value struct Config; use \*Config`
