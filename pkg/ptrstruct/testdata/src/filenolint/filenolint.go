//nolint:ptrstruct // entire file is legacy transport
package filenolint

type User struct {
	Name string
}

// Should not report: file-level nolint suppresses everything.
func Save(u User) {}

type Response struct {
	Meta User
}
