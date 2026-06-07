# dwm Configuration

This context defines the user-facing language for Razvan's personal dwm window-manager behavior.

## Language

**Dropdown terminal**:
A single reusable st terminal running a tmux session that can be toggled into view above the current workspace on the active monitor. While open, it remains visible across normal tag switches until explicitly hidden. Its manually resized dimensions are remembered separately per monitor. st stays internally opaque; picom owns dropdown translucency and blur through dropdown WM_CLASS rules.
_Avoid_: tdrop terminal, quake terminal, scratchpad terminal

**Scratchpad**:
A persistent window that can be hidden and shown without being tied to one normal tag.
_Avoid_: normal terminal, tagged terminal

**Polished minimal desktop**:
A restrained dwm setup that improves daily-use clarity and visual cohesion without adding a heavy rice stack or fragile visual dependencies.
_Avoid_: full rice, heavy desktop environment, effect-heavy setup

**Go status controller**:
The Go application that owns status-bar state, segment refresh scheduling, structured segment output, GPU sampling, and controller IPC for compact system telemetry. Controllers ask it to refresh named segments such as `volume`; they do not parse the X root window name or rebuild unrelated status values.
_Avoid_: shell status loop, xprop status parser, stale segment cache

**Native status bar**:
The dwm top-bar area used for compact system telemetry and time. Go owns the status data, while dwm renders structured status segments natively in the existing bar. Native status parsing, drawing, and click dispatch live in `features/status.c` to keep `dwm.c` focused on core hooks.
_Avoid_: polybar, eww bar, heavy status bar, external overlay bar

**Local feature modules**:
Focused C files under `features/` included by `dwm.c` after `config.h`. They keep local dwm extensions distinguishable from core dwm code while still sharing dwm internals.
_Avoid_: mixing feature implementation blocks into `dwm.c`, separate object build complexity

**Dropdown notes**:
A native dropdown st window running Neovim against the notes inbox under `$XDG_DOCUMENTS_DIR/notes`, currently toggled with Mod+n.
_Avoid_: dropdown-test, temporary test dropdown

**Dropdown widgets**:
A single persistent native dropdown st window running the Catppuccin-styled Go widget panel. Status-bar segments switch it to focused tabs such as Calendar or Audio instead of spawning separate widget terminals.
_Avoid_: one terminal process per status widget, external overlay widget windows

**Calendar event cache**:
A local provider-neutral cache of read-only calendar events used by the Calendar tab in the dropdown widget panel. It contains display-ready event fields such as title, provider, body, location, meeting URL, all-day flag, date, start, and end.
_Avoid_: Google calendar store, Outlook calendar store, live calendar API

**Audio widget**:
A tab in the persistent dropdown widget panel for PipeWire/PulseAudio outputs, inputs, and application streams. Audio actions refresh only the `volume` status segment through `dwm-status`.
_Avoid_: separate audio terminal, pavucontrol clone, recomputing unrelated status segments
