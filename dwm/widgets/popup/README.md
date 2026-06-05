# dwm popup widget prototype

Option C spike: a dedicated lightweight Go/X11 popup window for pixel-based widgets.

This intentionally does not use Kitty, Bubble Tea, or terminal rows/cols.

## Run manually

```sh
go run . -panel calendar
# or
go run . -panel audio -w 420 -h 320
```

By default the window uses `override_redirect` so it can be tested without changing dwm. Use `-managed` later when testing dwm rules/placement.

Keys:

- `q` / `Esc`: hide
- `h` / `l`: previous/next calendar month

## IPC prototype

Start a hidden long-lived popup process:

```sh
go run . -hidden
```

Control it from another shell:

```sh
go run ./cmd/dwm-popupctl toggle calendar
go run ./cmd/dwm-popupctl toggle audio
go run ./cmd/dwm-popupctl show calendar
go run ./cmd/dwm-popupctl hide
go run ./cmd/dwm-popupctl quit
```

Socket path:

```txt
${XDG_STATE_HOME:-$HOME/.local/state}/dwm/popup.sock
```

Current status: visual prototype only. Audio is placeholder data and performs no system changes.
