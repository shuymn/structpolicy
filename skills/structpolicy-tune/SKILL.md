---
name: structpolicy-tune
description: >
  Run ptrstruct and valuestruct on a Go codebase to enforce consistent struct
  pointer/value semantics and reduce unnecessary heap allocations. Trigger when
  the user asks to run structpolicy, ptrstruct, or valuestruct, or wants to tune
  struct pointer/value usage, reduce Go struct allocations, or check struct usage
  consistency.
---

# structpolicy-tune

Analyze a Go codebase with `ptrstruct` and `valuestruct`, then apply targeted
fixes to improve performance (fewer heap allocations, fewer unnecessary copies)
and enforce consistent struct usage patterns.

## Overview

`ptrstruct` flags struct types used **by value** where a **pointer** is expected
(e.g., large structs passed as parameters, method receivers copying state).

`valuestruct` flags struct types used **by pointer** where a **value** is
expected (e.g., small read-only structs returned as `*T` that force heap
allocation).

Running both analyzers together surfaces the full picture: places where you copy
too much *and* places where you allocate unnecessarily.

## Workflow

### 1. Capture baseline benchmark

Before making any changes, write and run benchmarks for the functions that will
be modified. This gives a concrete before/after comparison.

If the project already has benchmark tests (`*_bench_test.go`), run them:

```bash
go test -bench=. -benchmem -count=6 ./... | tee bench-before.txt
```

If not, write targeted benchmarks for the hot-path functions flagged by the
analyzers. Focus on functions that are called frequently (e.g., type walkers,
validators, per-field checkers) rather than one-time setup. A minimal benchmark:

```go
func BenchmarkFindViolation(b *testing.B) {
    // setup: create realistic inputs
    for b.Loop() {
        FindViolation(input, cfg, cls)
    }
}
```

Use `-benchmem` to capture allocs/op — this is the primary metric that
pointer/value changes affect.

### 2. Run both analyzers

```bash
ptrstruct ./... 2>&1
valuestruct ./... 2>&1
```

Collect all findings. A zero exit code means no violations.

### 4. Triage findings

Not every finding is actionable. Read the source of each flagged type before
deciding what to change. Classify each finding into one of three buckets:

#### Fix — internal types you control

These are the wins. Typical patterns:

| Analyzer     | Finding pattern                        | Fix                                                      |
|--------------|----------------------------------------|----------------------------------------------------------|
| ptrstruct    | receiver uses value struct `T`         | Change receiver to `*T`                                  |
| ptrstruct    | parameter uses value struct `T`        | Change parameter to `*T`                                 |
| valuestruct  | result uses pointer to struct `T`      | See "Choosing a return style" below                      |
| ptrstruct    | field uses value struct `T`            | See "Field value → pointer considerations" below         |

**Choosing a return style for valuestruct findings**

When valuestruct flags `*T` in a return position, there are several options.
Pick the one that fits the context:

| Option | When to use |
|--------|-------------|
| Return `T` directly | The function always succeeds (no "not found" case) |
| Return `(T, bool)` | Small, immutable struct with a clear found/not-found semantic — similar to map lookup or type assertion. **Only when the function does not also return `error`.** |
| Keep `*T` | Large struct (> ~128 bytes), mutable struct, when `*T` is the dominant convention in the surrounding code, or **when the function already returns `error`** |

`(T, bool)` eliminates a heap allocation and makes the "absent" case explicit,
but it is less common than `*T` + nil in user-defined Go functions. The Go
standard library uses `(T, bool)` extensively for built-in operations (map
lookup, type assertion, channel receive) and `sql.Null*` types adopt the same
idea as a struct field. For user-defined code, `*T` with nil is the more
familiar pattern for most Go developers.

**`(T, bool, error)` is not idiomatic Go.** When a function already returns
`error`, the `*T` + nil pattern is the standard way to express "not found" or
"absent". Do not introduce a `bool` alongside `error` — keep `*T` in this
case.

Default to **keeping `*T`** unless the function is on a proven hot path and
`-benchmem` shows significant allocs/op. The style departure of `(T, bool)` is
only worth it with benchmark evidence.

When converting `*T` returns to `(T, bool)`:
- Update **every** caller, including test files
- Replace `if v == nil` with `if !ok`
- Replace `*v` dereferences with direct value usage

When a function takes a struct parameter that callers always create locally
(stack-allocated from a value return), prefer `*T` for the parameter — the
caller passes `&v` and the pointer stays on the stack because the function only
reads it. This avoids copying while keeping the value off the heap.

**Field value → pointer considerations**

When ptrstruct flags a struct field as value where pointer is expected, evaluate
the nil-safety cost before changing:

| Situation | Action |
|-----------|--------|
| Field is always set at construction time and never observed as zero-value | Fix — change to `*T` |
| Field requires lazy initialization or nil checks to avoid panic (e.g., `if f.x == nil { f.x = &X{} }`) | **Skip** — the nil-safety burden outweighs the allocation benefit |
| Field is in a mutable struct where nil/non-nil carries semantic meaning | Investigate — weigh the API clarity vs. allocation trade-off |

The general principle: if changing a field to pointer forces defensive nil
checks at usage sites, the safety cost is too high unless benchmarks show a
clear hot-path benefit.

#### Skip — external / framework-constrained types

Do not change types when:

- **External API contract**: The type comes from a third-party package that
  requires pointer semantics. Common examples: `*analysis.Analyzer`
  (go/analysis framework), `*http.Request`, `*sql.DB`.

- **Mutable state**: The struct contains a cache, counter, or other state that
  mutates after construction. Returning by value would silently discard
  mutations. Look for map fields that get written to after `New*()` returns,
  `sync.Mutex` fields, or any method that modifies `self`.

- **Registered field addresses**: If something like `flag.BoolVar(&cfg.Field,
  ...)` or `a.Flags.BoolVar(&cfg.Receiver, ...)` stores a pointer to a field,
  the struct *must* have a stable address (heap-allocated, not copied). Returning
  by value would invalidate those pointers.

- **Slice element types (`[]*T` → `[]T`)**: If valuestruct flags slice elements,
  converting `[]*T` to `[]T` forces callers into index-based access patterns
  (`&slice[i]`, `d := &decls[i]`) to obtain pointers for mutation or passing to
  functions expecting `*T`. This is not idiomatic Go — the standard pattern is
  `for _, v := range slice` with direct use. Skip unless benchmarks prove a
  significant allocation win on a hot path.

- **Interface satisfaction**: If a pointer receiver is needed to satisfy an
  interface and the struct is returned by value, callers would need `&v` just to
  pass it around. If this creates widespread ergonomic harm, skip.

When skipping, note the reason briefly so the user understands why.

#### Investigate — ambiguous cases

Some findings need deeper analysis:

- **Large structs returned by value**: If valuestruct says "return `T`" but the
  struct is > ~128 bytes with no mutable state, returning by value is debatable.
  Check how often the function is called — hot paths benefit more from avoiding
  allocation; cold paths can tolerate a copy.

- **Competing findings**: ptrstruct and valuestruct may disagree on the same
  type in different positions (e.g., ptrstruct wants pointer params, valuestruct
  wants value returns). This is fine — the two checks cover different positions.
  Apply each fix to the position it targets.

### 5. Apply fixes

Work through the actionable findings methodically:

1. **Group by type** — all findings about the same struct type should be fixed
   together to maintain consistency.

2. **Fix the type's own methods first** (receivers, return types), then fix
   callers.

3. **Update tests** — test files that call modified functions need the same
   signature updates. Search with:
   ```
   grep -rn 'FunctionName' --include='*_test.go'
   ```

### 6. Verify

Run the project's standard check suite (lint, build, test). Use whatever the
project already has — a Makefile target, Taskfile, `go test ./...`, etc. Check
CLAUDE.md or the project README for the canonical command.

If linting fails (e.g., cognitive complexity), refactor the offending function —
typically by extracting a helper.

### 7. Re-run analyzers

After fixing, re-run both analyzers to confirm the actionable findings are
resolved and no new issues appeared:

```bash
ptrstruct ./... 2>&1
valuestruct ./... 2>&1
```

Report what remains and why each remaining finding is intentionally skipped.

### 8. Measure improvement

Re-run the same benchmarks from step 1:

```bash
go test -bench=. -benchmem -count=6 ./... | tee bench-after.txt
```

Compare with `benchstat`:

```bash
benchstat bench-before.txt bench-after.txt
```

The key metrics to look for:

- **allocs/op** — should decrease (fewer heap allocations from `*T` → `(T, bool)`)
- **B/op** — should decrease (less memory allocated per operation)
- **ns/op** — may decrease (less GC pressure, fewer copies of large structs)

If `benchstat` is not installed, `go install golang.org/x/perf/cmd/benchstat@latest`.

## Output

When done, provide a summary table:

```
| Type            | Analyzer    | Change                              | Status  |
|-----------------|-------------|-------------------------------------|---------|
| nolintMatcher   | ptrstruct   | receiver/params → pointer           | Fixed   |
| Violation       | valuestruct | *Violation return → (Violation,bool)| Fixed   |
| Config          | valuestruct | *Config return                      | Skipped |
| ...             | ...         | ...                                 | ...     |
```

Include the reason for each "Skipped" entry.

Also include the `benchstat` output showing the before/after comparison.
