package suppress

type User struct {
	Name string
}

func SaveInline(u User) {} //nolint:ptrstruct // legacy API

func SaveAll(u User) {} //nolint:all // suppress everything
