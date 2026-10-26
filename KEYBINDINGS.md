# Keybindings

This file is the high-level keyboard contract for Razvan's desktop. Keep detailed app-specific cheatsheets near their config; use this file to prevent namespace conflicts.

## Modifier namespaces

- **Opt / Super / Mod4**: dwm, launchers, window-manager actions.
- **Ctrl-b**: tmux prefix; terminal multiplexing actions only.
- **Bare vim keys**: text/navigation inside shells, editors, browsers, and pickers.
- **Shift**: stronger or resizing variant of an existing action.
- **Ctrl**: alternate/fallback variant when a namespace needs one.

## macOS / Ghostty

- **Ctrl + grave (`)**: toggle Ghostty's global dropdown terminal; configured in `ghostty/config.ghostty`.
- A login agent in `launchd/` configures the section/grave key swap system-wide for the MacBook's built-in keyboard. On newer macOS, `hidutil` may need Input Monitoring approval; verify actual typing after granting it. External keyboards and the FileVault/pre-login keyboard are unchanged. See [permission, installation, and undo instructions](launchd/README.md).

## Global / dwm

Opt is dwm's `MODKEY`.

| Key | Action |
| --- | --- |
| `Opt+Space` | rofi launcher, remembers `drun`/`window` mode |
| `Opt+`` | dropdown terminal |
| `Opt+n` | dropdown notes |
| `Opt+b` | toggle dwm bar |
| `Opt+h/j/k/l` | directional window navigation |
| `Opt+Ctrl+j/k` | old dwm stack-order focus next/previous |
| `Opt+Shift+h/l` | shrink/grow master area; dropdowns resize through the same binding when active |
| `Opt+Shift+j/k` | reserved for possible future stack-pane height resizing; intentionally unbound for now |
| `Opt+i/d` | increase/decrease number of master clients |
| `Opt+Return` | zoom focused client into master |
| `Opt+Shift+Return` | open normal st/tmux terminal |
| `Opt+q` | close focused window |
| `Opt+t/f/m` | tile/floating/monocle layout |
| `Opt+Shift+Space` | toggle focused window floating |
| `Opt+1..9` | view tag |
| `Opt+Shift+1..9` | move focused window to tag |
| `Opt+,/.` | focus previous/next monitor |
| `Opt+Shift+,/.` | move focused window to previous/next monitor |
| `Opt+Shift+q` | restart/quit dwm through the configured restart flow |

Directional window navigation means:

- tiled layout: `h/l` moves between master and stack; `j/k` cycles within the current column
- non-tiled layouts: `j/k` falls back to old stack-order focus; `h/l` intentionally does nothing
- `Shift` is the resize/expand layer; `Opt+Shift+j/k` stays reserved until stack-height resizing is worth the extra dwm complexity

## tmux

Prefix is `Ctrl-b`.

| Key | Action |
| --- | --- |
| `C-b g` | curated tmux help popup |
| `C-b ?` | raw tmux key table |
| `C-b o` | SessionX fuzzy session/project switcher |
| `C-b p` | project picker popup |
| `C-b e` | scratch shell popup in current pane directory; close with `exit` or `Ctrl-d` |
| `C-b u` | choose URL from pane history and open it |
| `C-b f` | fuzzy-filter scrollback and copy selected line |
| `C-b /` | native scrollback search |
| `C-b r` | reload tmux config |
| `C-b C-s` / `C-b C-r` | save/restore tmux-resurrect state |
| `C-b d` | detach session |
| `C-b c` | new window in current directory |
| `C-b n` | next window |
| `C-b ,` | rename window |
| `C-b &` | close window immediately |
| `C-b \|` | vertical split, left/right panes |
| `C-b -` | horizontal split, top/bottom panes |
| `C-b h/j/k/l` | move between panes |
| `C-b H/J/K/L` | resize pane |
| `C-b T` | set pane title |
| `C-b x` | close pane immediately |
| `C-b z` | zoom/unzoom pane |
| `C-b [` | copy mode |

Avoid adding prefixless tmux bindings unless there is a very strong reason. Shell, Vim, Pi, Claude, and browser key handling need the bare keyspace.

## Terminal input policy

st owns a few modified-key translations so tmux/Pi/AI CLIs can distinguish them.

| Key | Intended behavior |
| --- | --- |
| `Ctrl+v` | paste from clipboard |
| `Shift+Enter` | newline / modified Enter for AI CLIs |
| `Opt+Left/Right` | word backward/forward |
| `Alt+Left/Right` | character backward/forward |
| `Ctrl+Left/Right` | start/end of input line |
| `Ctrl+Backspace` | delete from cursor to start of line |
| `Opt+Backspace` | delete previous word |
| `Alt+Backspace` | delete one character |

Here **Opt** means `Mod4Mask`; **Alt** means `Mod1Mask`.

## rofi

Launched by `Opt+Space` through `~/.config/rofi/scripts/launcher`.

| Key | Action |
| --- | --- |
| `Tab` / `Shift+Tab` | switch rofi mode (`drun` / `window`) |
| `Ctrl+j/k` | move down/up |
| `Ctrl+n/p` | move down/up fallback |
| `Return` | accept entry |

## zsh / browser / editor convention

- zsh uses vi mode.
- Browser Vimium and editor navigation should keep bare `h/j/k/l`, `/`, `n`, `p`, etc. available.
- Do not assign global/dwm/tmux actions to bare vim keys.

## Status ownership

- **dwm bar**: global/system state.
- **tmux bar**: session, window, pane, directory, project/git context.
- **zsh prompt**: input point only, especially minimal inside tmux.
