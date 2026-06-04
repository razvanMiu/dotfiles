---
name: rofi-modifier
description: Modify Razvan's personalized rofi launcher configuration and Catppuccin themes safely. Use when changing /home/razvan/.config/rofi files, debugging rofi theme loading, styling rofi widgets, adjusting drun/run behavior, validating .rasi syntax, or aligning rofi with dwm/Catppuccin visuals.
---

# Rofi Modifier

## Quick start

Work in `/home/razvan/.config/rofi`.

```sh
find /home/razvan/.config/rofi -maxdepth 2 -type f -printf '%P\n' | sort
rofi -dump-config >/tmp/rofi-config-check.out
rofi -theme /home/razvan/.config/rofi/catppuccin-config.rasi -dump-theme >/tmp/rofi-theme-check.out
```

Do **not** replace the active theme, delete theme files, or change dwm launcher commands without making the dependency explicit and validating rofi.

## Local model

- Active rofi config: `/home/razvan/.config/rofi/config.rasi`.
- Active local styled theme: `/home/razvan/.config/rofi/catppuccin-config.rasi`.
- Palette files kept under `themes/`:
  - `themes/catppuccin-mocha.rasi`
  - `themes/catppuccin-macchiato.rasi`
- Non-Catppuccin themes/templates were intentionally removed.
- dwm launches rofi with `rofi -show drun`; theme should come from `config.rasi`, not from dwm arguments.
- Current font: `JetBrainsMono Nerd Font 10`.

## Workflow

1. **Inspect before editing**
   - Read `config.rasi` and `catppuccin-config.rasi`.
   - If dwm integration matters, inspect `/home/razvan/.config/dwm/config.h` rofi command.
   - Check current files with `find /home/razvan/.config/rofi -maxdepth 2 -type f`.

2. **Classify the change**
   - Launcher behavior: edit `configuration {}` in `config.rasi` or the dwm command only if necessary.
   - Visual styling: edit `catppuccin-config.rasi`.
   - Palette values: edit/import `themes/catppuccin-*.rasi` only deliberately; prefer styling aliases in `catppuccin-config.rasi`.
   - Dwm integration: keep dwm command simple: `rofi -show drun`.

3. **Edit carefully**
   - Prefer `@import "catppuccin-config.rasi"` in `config.rasi`.
   - Avoid passing `-theme` from dwm unless testing or intentionally overriding config.
   - Use absolute paths only when command arrays require paths; rofi config imports are relative to the config file.
   - Keep no-results/list spacing tidy: `fixed-num-lines: false`, `listview dynamic: true`, and avoid unconditional `mainbox spacing`.
   - Validate after every `.rasi` edit.

4. **Validate**
   ```sh
   rofi -dump-config >/tmp/rofi-config-check.out 2>/tmp/rofi-config-check.err
   rofi -theme /home/razvan/.config/rofi/catppuccin-config.rasi -dump-theme >/tmp/rofi-theme-check.out 2>/tmp/rofi-theme-check.err
   ```
   Both commands must exit `0`; inspect stderr on failure.

## Current theme semantics

`catppuccin-config.rasi` imports Mocha and defines styled widgets directly, not through a rounded template.

Important widgets:

- `window`: centered 560px launcher with rounded border.
- `mainbox`: outer padding; keep `spacing: 0` to avoid empty bottom gap.
- `inputbar`: search field container.
- `listview`: result list; uses `dynamic: true`, `fixed-height: false`, no border.
- `element`: app rows.
- `element selected.*`: selected row accent.
- `message`: no border; used for rofi messages.

## Common pitfalls

- `rofi -show drun` opens desktop-app launcher mode; `run` opens PATH command mode.
- Dwm `spawn` does not shell-expand `~`; avoid `~` in dwm command arrays.
- `@import "themes/foo.rasi"` is relative to the importing rofi config file.
- Palette files alone only define colors; full widget styling lives in `catppuccin-config.rasi`.
- `mainbox spacing` can create awkward bottom space when search has no results.
- Old dashed separators usually come from `border` on `listview`, `message`, or inherited theme snippets.
- Always validate with `-dump-config`/`-dump-theme`; this catches syntax/import failures without opening the UI.
