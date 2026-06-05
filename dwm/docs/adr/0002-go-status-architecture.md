# ADR 0002: Go-owned status with native dwm rendering

## Status

Accepted.

## Context

The original shell status path treated the dwm bar as one plain text string. That made targeted controller updates fragile:

- reading root status with `xprop` corrupts Nerd Font UTF-8 glyphs;
- caching individual shell segments can write stale CPU/GPU/RAM/time values back into the bar;
- every new controller risks duplicating parsing or refresh logic;
- a separate Go/X11 overlay bar was visually hard to integrate with dwm and introduced too many rendering dependencies.

The desktop direction remains a polished minimal dwm setup, not a heavy external bar stack.

## Decision

Use a split architecture:

- `dwm-status` is a Go daemon/controller that owns status state, segment refresh scheduling, GPU sampling, and Unix-socket IPC.
- dwm owns the top-bar pixels and renders the status natively inside the existing bar.
- `dwm-status` emits ordered structured segments by joining segment text with ASCII unit separator (`0x1f`) in the root window name.
- dwm treats that separator as structure, not visible text: it measures, draws, and hit-tests each segment separately.
- Controllers refresh only their corresponding segment, e.g. `dwm-status --refresh volume`.
- The old shell status implementation is removed.

## Consequences

- Go remains the place to add, remove, reorder, sample, or refresh status elements.
- dwm stays small: it only understands segment boundaries and native click actions.
- The bar is visually integrated because dwm renders it using the same font, colors, height, and monitor geometry as the rest of the top bar.
- `xprop` must not be used as status state.
- External overlay bars are avoided unless future requirements exceed what native dwm rendering can reasonably support.
