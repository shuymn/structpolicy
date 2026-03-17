# ptrstruct

`ptrstruct` reports declarations where a struct type is used by value instead of by pointer.

By default, `ptrstruct` uses the pointer-leaning side of the opposite performance-tuning profiles. It aims at surfacing `copy hotspot` candidates around call boundaries and declaration sites, so receivers, parameters, fields, interface method signatures, and function type declarations default to on. Results and container-heavy shapes stay opt-in on this side.

| | NG | OK | Flag | Default |
|---|---|---|---|---|
| Receiver | `func (u User) M()` | `func (u *User) M()` | `-receiver` | on |
| Parameter | `func Save(u User)` | `func Save(u *User)` | `-param` | on |
| Result | `func Load() User` | `func Load() *User` | `-result` | off |
| Field | `Meta Profile` | `Meta *Profile` | `-field` | on |
| Slice element | `[]User` | `[]*User` | `-slice-elem` | off |
| Map value | `map[string]User` | `map[string]*User` | `-map-value` | off |

Interface method and function type checks are also enabled by default. Named container types and container element checks remain opt-in.

Pointer wrapping a container does not exempt its contents: `*[]User` is still flagged because the slice element `User` is by value.

Empty structs (`struct{}`) are exempt.

## Installation

### Standalone

```bash
go install github.com/shuymn/structpolicy/cmd/ptrstruct@latest
```

How to use:

```bash
ptrstruct ./...
```

### golangci-lint

Use the [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) to add ptrstruct as a custom linter.

`.custom-gcl.yml`:

```yaml
version: v2.11.3

plugins:
  - module: github.com/shuymn/structpolicy
    path: /path/to/structpolicy
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
          allow_stdlib: true
          allow_third_party: true
          allow_types:
            - github.com/google/uuid.UUID
```

## Configuration

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-receiver` | `true` | Check method receivers |
| `-param` | `true` | Check function parameters |
| `-result` | `false` | Check function results |
| `-field` | `true` | Check struct fields |
| `-interface-method` | `true` | Check interface method signatures |
| `-func-type` | `true` | Check function type declarations |
| `-named-type` | `false` | Check named container types |
| `-slice-elem` | `false` | Check slice element types |
| `-map-value` | `false` | Check map value types |
| `-map-key` | `false` | Check map key types |
| `-array-elem` | `false` | Check array element types |
| `-chan-elem` | `false` | Check channel element types |
| `-allow-stdlib` | `true` | Exempt builtin and standard library packages |
| `-allow-third-party` | `false` | Exempt non-stdlib packages outside the current Go module |
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

`ptrstruct` enforces a stricter policy: struct receivers **must** be pointers, period. Its default profile also checks parameters, struct fields, interface methods, and function type declarations to surface pointer-leaning refactor candidates. The two tools are complementary: `recvcheck` catches mixed receiver sets, `ptrstruct` catches value-struct usage where pointer semantics are preferred.
