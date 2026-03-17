# ADR-005: Performance-tuning default flag profiles

## Status

Accepted

## Context

The initial default behavior checked only method receivers. That was easy to explain, but it underfit the main use case for this repository: using analyzers to find internal pointer-vs-value refactor candidates.

For performance work, `ptrstruct` and `valuestruct` do not benefit from perfectly symmetric defaults:

- `ptrstruct` is most useful when it surfaces copy-heavy positions.
- `valuestruct` is most useful when it surfaces unnecessary pointers, allocations, and indirection.
- Receiver checks are not equally strong signals in both directions because standard-library value types can still require pointer receivers for mutating or unmarshal methods.

## Decision

Adopt mode-specific performance-tuning defaults.

- `ptrstruct` defaults to receiver, parameter, field, slice element, map value, array element, and channel element checks.
- `valuestruct` defaults to parameter, result, field, slice element, map value, array element, and channel element checks.
- `result` remains disabled by default for `ptrstruct`.
- `receiver` remains disabled by default for `valuestruct`.
- `map-key`, `interface-method`, `func-type`, and `named-type` remain opt-in for both analyzers.

## Consequences

- Default runs surface more actionable refactor candidates for performance work.
- `ptrstruct` and `valuestruct` now have intentionally different default profiles.
- Existing users who relied on receiver-only defaults must pass flags explicitly to preserve the old behavior.
