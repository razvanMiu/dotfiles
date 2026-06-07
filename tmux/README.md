# tmux

XDG-friendly tmux setup for moving from kitty tabs to `st` + tmux.

## Layout

- Tracked config: `~/.config/tmux/tmux.conf`
- Plugin/vendor code: `~/.local/share/tmux/plugins`
- Saved session state: `~/.local/state/tmux/resurrect`
- Config is split under `~/.config/tmux/conf.d/`; `tmux.conf` is only the source order
- Scripts: `~/.config/tmux/scripts/git-branch`, `project-session`, `open-url`, `find-scrollback`
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

TPM itself is intentionally not tracked in this repo:

```sh
mkdir -p ~/.local/share/tmux/plugins ~/.local/state/tmux/resurrect
git clone https://github.com/tmux-plugins/tpm ~/.local/share/tmux/plugins/tpm
```

Inside tmux, install/update plugins with:

```text
prefix + I  # install
prefix + U  # update
```

Current prefix is `C-b`.

## Pi compatibility

`tmux.conf` enables:

```tmux
set -g extended-keys on
```

This is not a tmux shortcut. It tells tmux to pass richer modified-key sequences through to terminal apps, so Pi can distinguish combinations such as Ctrl/Shift/Alt-modified keys when launched inside tmux.

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
