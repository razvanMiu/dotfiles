# st

Razvan's XDG-tracked suckless `st` build.

## Current choices

- Font: `JetBrainsMono Nerd Font:size=12`
- Theme: Catppuccin Macchiato
- TERM: `st-256color`
- Opacity: `1.0`; dropdown translucency/blur is owned by picom rules, not st
- Clipboard: selections are copied to X11 CLIPBOARD; `Ctrl+v`, middle-click, and Shift+Insert paste CLIPBOARD
- AI CLI modified keys emit CSI-u for tmux/Pi/Claude-style input handling: Shift+Tab (`\033[9;2u`), Shift+Enter (`\033[13;2u`), Shift+Backspace (`\033[127;2u`), Alt+Space (`\033[32;3u`), and Ctrl+Space (`\033[32;5u`)
- Navigation policy: Opt uses `Mod4Mask`, Alt uses `Mod1Mask`. Opt+Left/Right sends Alt+b/f for word movement; Ctrl+Left/Right sends Ctrl+A/E for start/end of the input line; Alt+Left/Right sends plain left/right arrow movement.
- Backspace policy: Ctrl+Backspace sends Ctrl+U (`\025`) to kill from cursor to start of line; Opt+Backspace uses Mod4Mask and sends Alt+Backspace (`\033\177`) for backward-word deletion; Alt+Backspace uses Mod1Mask and sends plain Backspace (`\177`) for one-character deletion.
- Scrollback: intentionally delegated to tmux, not patched into st

## tmux compatibility

Install this build's terminfo after build/install:

```sh
make -C ~/.config/st install
```

The install target uses `tic -sx -o $XDG_DATA_HOME/terminfo st.info`, which installs the `st-256color` terminfo entry under `~/.local/share/terminfo` instead of legacy `~/.terminfo`. tmux is configured separately with `default-terminal "tmux-256color"` and RGB terminal features for `st-256color`.

Recommended launch command:

```sh
st -A 1.0 -e tmux new-session -A -s main
```

## Build

```sh
make -C ~/.config/st clean && make -C ~/.config/st
make -C ~/.config/st install
```

No sudo is needed: `config.mk` installs to `~/.local`.

## Patch policy

Keep the build lean. Good additions later:

- `boxdraw`, if tmux borders/glyphs show gaps with the chosen font/size
- `font2`, only if Nerd Font fallback proves insufficient

Avoid st scrollback unless there is a specific reason; tmux already provides scrollback/copy mode and avoids maintaining a larger patch stack.
