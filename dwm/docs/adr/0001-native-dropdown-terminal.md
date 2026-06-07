# Native dropdown terminal

We replaced the external `tdrop`-managed Kitty dropdown with a dwm-native dropdown terminal because the desired behavior depends on dwm's own monitor work area, tag visibility, and stacking rules. Stock scratchpad patches and hidden scratchpad tags were considered, but the chosen behavior is more specific: a single sticky-while-open st window running tmux that can be resized, remembers width/height per monitor while dwm is running, moves to the selected monitor when toggled there, and stays above other windows. st stays internally opaque for dropdowns; picom owns dropdown translucency and blur through dropdown WM_CLASS rules.

## Considered Options

- Keep `tdrop`: already works, but requires geometry hacks for bar offset, border overflow, and dwm-specific restore behavior.
- Apply a stock scratchpad patch: closer to dwm-native behavior, but still needs custom multi-monitor placement and sticky-while-open semantics.
- Implement a custom native dropdown: more local C code, but directly models the desired behavior using dwm's `Monitor` work-area geometry.

## Consequences

The implementation will live mostly in an included `features/dropdown.c` file to keep the main `dwm.c` readable while still sharing dwm's static internals. It is configured through `dropdowns[]`, rule dropdown indexes, and keybindings that pass a dropdown index, so additional dropdown-like windows can be added from `config.h` without duplicating core logic. Small changes to `dwm.c` remain necessary for client/monitor state, key dispatch, protected operations, and z-order integration.
