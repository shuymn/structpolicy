# valuestruct

`valuestruct` reports declarations where a struct type is used by pointer instead of by value.

This is the inverse of [`ptrstruct`](../ptrstruct/). Use it to enforce value semantics for lightweight or immutable struct types.

By default, `valuestruct` uses a performance-tuning profile aimed at surfacing `allocation / indirection hotspot` candidates. It checks parameters, results, fields, and container element positions. Receivers stay opt-in because standard-library value types such as `time.Time` still use pointer receivers for mutation and unmarshal helpers.

| | NG | OK | Flag | Default |
|---|---|---|---|---|
| Receiver | `func (u *User) M()` | `func (u User) M()` | `-receiver` | off |
| Parameter | `func Save(u *User)` | `func Save(u User)` | `-param` | on |
| Result | `func Load() *User` | `func Load() User` | `-result` | on |
| Field | `Meta *Profile` | `Meta Profile` | `-field` | on |
| Slice element | `[]*User` | `[]User` | `-slice-elem` | on |
| Map value | `map[string]*User` | `map[string]User` | `-map-value` | on |

Empty structs (`struct{}`) are exempt.

Array and channel element checks are also enabled by default. Map key checks remain opt-in.

## Installation

### Standalone

```bash
go install github.com/shuymn/structpolicy/cmd/valuestruct@latest
```

How to use:

```bash
valuestruct ./...
```

### golangci-lint

Use the [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) to add valuestruct as a custom linter.

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
    - valuestruct

  settings:
    custom:
      valuestruct:
        type: module
        settings:
          allow_stdlib: true
          allow_third_party: true
          allow_types:
            - github.com/google/uuid.UUID
```

## Configuration

### Flags

`valuestruct` shares the same flag set as `ptrstruct`.

| Flag | Default | Description |
|------|---------|-------------|
| `-receiver` | `false` | Check method receivers |
| `-param` | `true` | Check function parameters |
| `-result` | `true` | Check function results |
| `-field` | `true` | Check struct fields |
| `-interface-method` | `false` | Check interface method signatures |
| `-func-type` | `false` | Check function type declarations |
| `-named-type` | `false` | Check named container types |
| `-slice-elem` | `true` | Check slice element types |
| `-map-value` | `true` | Check map value types |
| `-map-key` | `false` | Check map key types |
| `-array-elem` | `true` | Check array element types |
| `-chan-elem` | `true` | Check channel element types |
| `-allow-stdlib` | `true` | Exempt builtin and standard library packages |
| `-allow-third-party` | `false` | Exempt non-stdlib packages outside the current Go module |
| `-allow-types` | | Comma-separated fully qualified type names to exempt (e.g. `time.Time`) |
| `-allow-patterns` | | Comma-separated regex patterns for type names to exempt |
| `-allow-packages` | | Comma-separated package paths to exempt |
| `-ignore-generated` | `true` | Skip generated files |
| `-ignore-tests` | `false` | Skip test files |
| `-honor-nolint` | `true` | Honor `//nolint:valuestruct` comments |
| `-honor-nolint-all` | `true` | Honor `//nolint:all` comments |

### Suppression

Use `//nolint:valuestruct` to suppress diagnostics:

```go
func LoadLegacy() *User {} //nolint:valuestruct // public legacy API

//nolint:valuestruct // compatibility layer
type LegacyResponse struct {
    Meta *Meta
}
```

File-level suppression before the package clause:

```go
//nolint:valuestruct // frozen legacy transport package
package legacytransport
```
