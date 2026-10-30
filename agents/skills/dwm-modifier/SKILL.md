---
name: dwm-modifier
description: Suckless dwm window manager configuration, native C patches, and X11 rules in ~/.config/dwm. Use when modifying dwm config.h, patches, dropdown/status features, or debugging window rules and keybindings.
---

# DWM Modifier

## Quick start

Work in `~/.config/dwm`.

```sh
git -C ~/.config status --short dwm
make -C ~/.config/dwm clean && make -C ~/.config/dwm
```

Do **not** install, kill X, restart dwm, or press restart keybindings unless the user explicitly approves.

## Local model

- Active build config: `config.h`.
- Template/default config: `config.def.h`; keep structurally compatible when changing public config types, but remember existing builds use `config.h`.
- Main source: `dwm.c`.
- Local feature modules live under `features/`, included from `dwm.c` after `config.h`: `features/dropdown.c` and `features/status.c`.
- Domain docs: read `CONTEXT.md` and relevant `docs/adr/*.md` before changing terminology or behavior.
- Session startup: `dwm.desktop` -> `~/.config/X11/xinitrc` -> `~/.local/bin/startdwm.sh`.
- Installed binary target: usually `/usr/local/bin/dwm` from `config.mk`.

## Workflow

1. **Inspect before editing**
   - Check git status under `~/.config`.
   - Read relevant parts of `config.h`, `config.def.h`, `dwm.c`, included patch files, `CONTEXT.md`, and ADRs.
   - For runtime matching issues, inspect `xprop WM_CLASS WM_NAME` evidence.

2. **Classify the change**
   - Config-only: keybindings, rules, commands, layout constants.
   - Template-shape change: update both `config.h` and `config.def.h`.
   - Native feature/patch: keep `dwm.c` hooks minimal; prefer an included focused file under `features/` for feature logic.
   - Runtime/debug issue: reproduce or get concrete X11/window evidence before hypothesizing.

3. **Edit carefully**
   - Preserve suckless C style: C99, file-local `static`, tabs, declarations before statements, `/* */` comments, simple functions.
   - Avoid broad WM_CLASS rules. Prefer unique classes, e.g. `dropdown-terminal`.
   - Keep one logical change per edit batch when possible.
   - Never regenerate or overwrite `config.h` from `config.def.h` without asking.

4. **Validate & Completion criteria**
   - Clean compilation:
     ```sh
     make -C ~/.config/dwm clean && make -C ~/.config/dwm
     ```
   - Install preview:
     ```sh
     make -C ~/.config/dwm -n install
     ```
   - Done criteria:
     1. `make clean && make` exits with status `0` and zero compiler warnings.
     2. `git diff -C ~/.config dwm` confirms only intended changes are touched.
     3. Changed files, validation output, residual risks, and exact runtime test steps are reported.

## Dropdown-specific rules

Current native dropdown semantics:

- Dropdowns are configured by `Dropdown dropdowns[]` in `config.h`.
- `Rule.dropdown` maps a matched client to a dropdown index.
- Keybindings call `toggledropdown` with `{.i = index}`.
- Only one dropdown should be visible at a time.
- Visible dropdowns are sticky via `TAGMASK`; hidden dropdowns use `tags = 0`.
- Dropdown size is remembered in memory per monitor and per dropdown index.
- Dropdown placement uses monitor work area (`wx/wy/ww/wh`), not global screen dimensions.
- The command class must match the rule class, e.g. Kitty `--class dropdown-test` with rule class `dropdown-test`.

To add a dropdown, modify `config.h` like:

```c
static const char *namecmd[] = { "kitty", "--class", "dropdown-name", NULL };

static const Dropdown dropdowns[] = {
	{ droptermcmd, 1.0, 0.5 },
	{ namecmd,    0.8, 0.4 },
};

{ "dropdown-name", NULL, NULL, 0, 1, 1, -1 },
{ MODKEY, XK_n, toggledropdown, {.i = 1 } },
```

## Install/restart protocol

Only after explicit user approval: `cd ~/.config/dwm && sudo make install`. Then the user can press `MOD+Shift+q`; `startdwm.sh` relaunches only if the installed binary checksum changed.

## Common pitfalls

- Editing only `config.def.h` does not affect active `config.h`.
- External suckless patches often target `config.def.h`; port relevant parts to `config.h` manually.
- Duplicate keybindings in dwm can both execute because keypress iterates matching entries.
- `spawn` does not expand shell `~` unless using `SHCMD` or absolute paths.
- Missing `~/.local/state/dwm` can hide stderr log output.
- Multi-monitor features should use monitor work area (`wx/wy/ww/wh`) unless screen geometry is explicitly intended.
