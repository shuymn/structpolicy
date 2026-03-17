# Default Performance Profiles

## Status

Approved

## Goal

Define default flag profiles that help internal performance-tuning refactors surface actionable pointer-vs-value candidates without forcing a single global style rule.

## Decision

- `ptrstruct` default profile targets copy reduction.
- `valuestruct` default profile targets allocation and indirection reduction.
- The defaults are intentionally not mirror images.

## Default Profiles

### ptrstruct

- Enable `receiver`, `param`, `field`, `slice-elem`, `map-value`, `array-elem`, `chan-elem`.
- Leave `result`, `map-key`, `interface-method`, `func-type`, and `named-type` disabled.

### valuestruct

- Enable `param`, `result`, `field`, `slice-elem`, `map-value`, `array-elem`, `chan-elem`.
- Leave `receiver`, `map-key`, `interface-method`, `func-type`, and `named-type` disabled.

## Rationale

- Pointer-oriented types often benefit from receiver, parameter, field, and container checks because copy costs compound across those positions.
- Value-oriented types often benefit from parameter, result, field, and container checks because unnecessary pointers add indirection and may increase heap pressure.
- `valuestruct` leaves receiver checks off by default because standard-library value types still use pointer receivers for mutating or unmarshal methods.
- Both analyzers leave `map-key` off because changing key representation can alter equality semantics and introduces more risk than typical performance wins.

## Acceptance Criteria

- When `ModePointer` is used, the analyzer shall enable the copy-reduction default profile.
- When `ModeValue` is used, the analyzer shall enable the allocation-reduction default profile.
- The README and command documentation shall describe the new defaults and their intent.
- Tests shall assert the mode-specific defaults.
