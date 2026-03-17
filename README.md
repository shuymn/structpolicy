# ptrstruct

A Go static analyzer that enforces pointer usage for struct-bearing declaration types. Built on [`golang.org/x/tools/go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis).

## What it does

`ptrstruct` reports declarations where a struct type is used by value instead of by pointer.

| | NG | OK |
|---|---|---|
| Receiver | `func (u User) M()` | `func (u *User) M()` |
| Parameter | `func Save(u User)` | `func Save(u *User)` |
| Result | `func Load() User` | `func Load() *User` |
| Field | `Meta Profile` | `Meta *Profile` |
| Slice element | `[]User` | `[]*User` |
| Map value | `map[string]User` | `map[string]*User` |

Pointer wrapping a container does not exempt its contents: `*[]User` is still flagged because the slice element `User` is by value.

Empty structs (`struct{}`) are exempt.

## Installation

### Standalone

```bash
go install github.com/shuymn/ptrstruct/cmd/ptrstruct@latest
```

```bash
ptrstruct ./...
```

### golangci-lint

Use the [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) to add ptrstruct as a custom linter.

`.custom-gcl.yml`:

```yaml
version: v2.11.3

plugins:
  - module: github.com/shuymn/ptrstruct
    path: /path/to/ptrstruct
```

Build a custom binary with `golangci-lint custom`, then configure `.golangci.yml`:

```yaml
linters:
  enable:
    - ptrstruct

  settings:
    custom:
      ptrstruct:
        type: module
        settings:
          allow_types:
            - time.Time
            - github.com/google/uuid.UUID
```

## Configuration

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-receiver` | `true` | Check method receivers |
| `-param` | `true` | Check function parameters |
| `-result` | `true` | Check function results |
| `-field` | `true` | Check struct fields |
| `-slice-elem` | `true` | Check slice element types |
| `-map-value` | `true` | Check map value types |
| `-map-key` | `false` | Check map key types |
| `-array-elem` | `false` | Check array element types |
| `-chan-elem` | `false` | Check channel element types |
| `-allow-types` | | Comma-separated fully qualified type names to exempt (e.g. `time.Time`) |
| `-allow-patterns` | | Comma-separated regex patterns for type names to exempt |
| `-allow-packages` | | Comma-separated package paths to exempt |
| `-ignore-generated` | `true` | Skip generated files |
| `-ignore-tests` | `false` | Skip test files |
| `-honor-nolint` | `true` | Honor `//nolint:ptrstruct` comments |
| `-honor-nolint-all` | `true` | Honor `//nolint:all` comments |

### Suppression

Use `//nolint:ptrstruct` to suppress diagnostics:

```go
func LoadLegacy() User {} //nolint:ptrstruct // public legacy API

//nolint:ptrstruct // compatibility layer
type LegacyResponse struct {
    Meta Meta
}
```

File-level suppression before the package clause:

```go
//nolint:ptrstruct // frozen legacy transport package
package legacytransport
```

## Difference from recvcheck

[recvcheck](https://github.com/raeperd/recvcheck) enforces receiver type consistency — if a type has both value and pointer receivers, it flags the inconsistency. It does not care whether receivers are pointers or values, only that they are uniform.

`ptrstruct` enforces a stricter policy: struct receivers **must** be pointers, period. It also goes beyond receivers to check parameters, results, struct fields, and container elements. The two tools are complementary: `recvcheck` catches mixed receiver sets, `ptrstruct` prevents value structs from appearing in API surfaces.

## Local Development

Requires [Task](https://taskfile.dev/) as the build interface.

```bash
task          # list available tasks
task build    # build the binary
task test     # run tests with race detection, shuffle, count=10
task lint     # run golangci-lint
task check    # full verification (lint + build + test + tidy)
```

## License

[MIT](LICENSE)
