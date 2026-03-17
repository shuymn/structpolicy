# ADR-005: Performance-tuning default flag profiles

## Status

Accepted

## Context

The initial default behavior checked only method receivers. That was easy to explain, but it underfit the main use case for this repository: using analyzers to find internal pointer-vs-value refactor candidates.

For the "push in one direction first, then adjust while tuning" workflow, `ptrstruct` and `valuestruct` benefit from clearly opposite defaults:

- `ptrstruct` should bias toward pointer-leaning call boundaries and declaration sites.
- `valuestruct` should bias toward value-leaning returns and container-heavy shapes.
- Exact opposites make the starting policy obvious and easy to reason about.

## Decision

Adopt mode-specific performance-tuning defaults.

- `ptrstruct` defaults to receiver, parameter, field, interface method, and function type checks.
- `valuestruct` defaults to result, named type, slice element, map value, map key, array element, and channel element checks.
- The two profiles are exact opposites across these flags:
  `receiver`, `param`, `result`, `field`, `interface-method`, `func-type`, `named-type`, `slice-elem`, `map-value`, `map-key`, `array-elem`, `chan-elem`.
- `interface-method`, `func-type`, and `named-type` are implemented as real checks so their defaults are behaviorally meaningful.

## Consequences

- Default runs support a clear first-pass push toward pointer-heavy or value-heavy code.
- `ptrstruct` and `valuestruct` now have intentionally opposite default profiles.
- Existing users who relied on receiver-only defaults must pass flags explicitly to preserve the old behavior.
