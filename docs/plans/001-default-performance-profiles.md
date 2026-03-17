# Default Performance Profiles

## Status

Approved

## Goal

Define opposite default flag profiles that help internal performance-tuning refactors start by pushing code toward one side, then adjust based on measurement.

## Decision

- `ptrstruct` default profile is the pointer-leaning side.
- `valuestruct` default profile is the value-leaning side.
- The defaults are intentionally exact opposites.

## Default Profiles

### ptrstruct

- Enable `receiver`, `param`, `field`, `interface-method`, and `func-type`.
- Leave `result`, `named-type`, `slice-elem`, `map-value`, `map-key`, `array-elem`, and `chan-elem` disabled.

### valuestruct

- Enable `result`, `named-type`, `slice-elem`, `map-value`, `map-key`, `array-elem`, and `chan-elem`.
- Leave `receiver`, `param`, `field`, `interface-method`, and `func-type` disabled.

## Rationale

- The pointer-leaning side emphasizes call-boundary and declaration-site hotspots first.
- The value-leaning side emphasizes returns and container-heavy hotspots first.
- Exact opposites make it easy to do a first pass in either direction without designing a custom profile up front.
- The trade-off is that some defaults are less risk-weighted than a hand-tuned profile, so users are expected to adjust after measuring.

## Acceptance Criteria

- When `ModePointer` is used, the analyzer shall enable the pointer-leaning default profile.
- When `ModeValue` is used, the analyzer shall enable the value-leaning default profile.
- The README and command documentation shall describe the new defaults and their intent.
- Tests shall assert the mode-specific defaults.
