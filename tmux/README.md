# tmux

XDG-friendly tmux setup for moving from kitty tabs to `st` + tmux.

## Layout

- Tracked config: `~/.config/tmux/tmux.conf`
- Plugin/vendor code: `~/.local/share/tmux/plugins`
- Saved session state: `~/.local/state/tmux/resurrect`
- Config is split under `~/.config/tmux/conf.d/`; `tmux.conf` is only the source order
- Scripts: `~/.config/tmux/scripts/git-branch`, `project-session`, `open-url`, `find-scrollback`, `sessionx-launch`
- Status helpers: `~/.config/tmux/status/git_branch.conf`

## Dependencies

Required for the configured workflow:

```sh
sudo pacman -S tmux fzf bat zoxide xclip
```

Notes:

- `tmux-sessionx` requires `fzf`, `fzf-tmux` from the fzf package, and `bat`; `zoxide` is optional but enabled.
- `xclip` is used for X11 clipboard copy from tmux copy mode and `tmux-yank`.
- Install/build `st` separately from the suckless source you want to use.

## Plugin bootstrap

Runtime plugin loading is explicit and deterministic: `conf.d/30-plugins.conf` and `conf.d/20-theme.conf` source the installed plugin files directly from `~/.local/share/tmux/plugins`. This avoids TPM hook-order races on reboot/reload.

Required plugin directories:

```text
~/.local/share/tmux/plugins/tmux
~/.local/share/tmux/plugins/tmux-sessionx
~/.local/share/tmux/plugins/tmux-sensible
~/.local/share/tmux/plugins/tmux-yank
~/.local/share/tmux/plugins/vim-tmux-navigator
~/.local/share/tmux/plugins/tmux-resurrect
~/.local/share/tmux/plugins/tmux-continuum
```

TPM can still live under `~/.local/share/tmux/plugins/tpm` as a bootstrap/update tool, but it is not part of tmux runtime startup. Current prefix is `C-b`.

## Pi compatibility

`tmux.conf` enables:

```tmux
set -g extended-keys on
```

This is not a tmux shortcut. It tells tmux to pass richer modified-key sequences through to terminal apps, so Pi can distinguish combinations such as Ctrl/Shift/Alt-modified keys when launched inside tmux. The local `st` build emits CSI-u for AI-CLI-sensitive keys such as Shift+Tab, Shift+Enter, Shift+Backspace, Alt+Space, and Ctrl+Space.

Reload after config edits with `C-b r`, or from a shell with:

```sh
tmux source ~/.config/tmux/tmux.conf
```

## Main bindings

The scheme mirrors the rest of the desktop: dwm owns `Super`, rofi owns `Super+Space`, browser Vimium/zsh vi mode own bare vim keys, and tmux actions require the tmux prefix.

- `C-b g`: workflow cheat sheet popup
- `C-b ?`: raw tmux key table
- `C-b o`: SessionX fuzzy session/window/project picker
- `C-b p`: project picker that attaches/creates a tmux session for a repo
- `C-b e`: scratch shell popup in current directory
- `C-b u`: pick/open URL from pane history
- `C-b f`: fuzzy-filter scrollback and copy selected line
- `C-b /`: native scrollback search in copy-mode
- `C-b C-s` / `C-b C-r`: explicit tmux-resurrect save/restore
- `C-b |`: horizontal split in current directory
- `C-b -`: vertical split in current directory
- `C-b c`: new window in current directory
- `C-b h/j/k/l`: pane navigation, matching Vim/Vimium directions
- `C-b H/J/K/L`: pane resize
- `C-b r`: reload config
- status-right: pane directory and Catppuccin-shaped current git branch/dirty marker when inside a repo
- pane borders show each pane's current directory, plus manual pane title when set with `C-b T`, so inactive panes keep visible context
- copy/search: `C-b /` native search, `C-b f` fuzzy-filter/copy a line, `C-u/C-d` animated 16-line scroll, `v` begins selection, `y` copies to clipboard
