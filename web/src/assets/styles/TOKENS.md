# CSS Token Guide

This template uses a 4-level token model:

1. `ref` - raw design values (`tokens.ref.css`)
2. `sys` - theme-aware semantic values (`tokens.sys.css`)
3. `cmp` - component contract values (`tokens.cmp.css`)
4. `state` - shared state colors (defined in `tokens.sys.css`)

## Files

- `tokens.ref.css`: hex/rgb/rgba source values only
- `tokens.sys.css`: application-level semantic tokens for light/dark theme
- `tokens.cmp.css`: component-level aliases built on system tokens
- `tokens.css`: import entry

## Naming Rules

- Reference: `--ref-*`
- System: `--sys-*`
- Component: `--cmp-*`
- State: `--state-*`
- Local component private variables: `--nav-*`, `--dropdown-*`

## Usage Rules

- Components and pages must use `sys`, `cmp`, or `state` tokens.
- Do not use `ref` tokens directly in component styles.
- Do not add hex colors in component/page styles.
- Keep chart business palettes inside chart components if they are domain-specific.
- Shared status meanings (success/warning/error/info) should use `state` tokens.

## Adding a New Token

1. Add raw value in `tokens.ref.css` if needed.
2. Map it into semantic meaning in `tokens.sys.css`.
3. Add component alias in `tokens.cmp.css` only when multiple components need the same contract.
4. Consume semantic token in components.

## Migration Principle

- Preserve current UI output.
- Prefer alias/mapping over ad-hoc component overrides.
- Remove dead tokens after migration.
