---
name: dwm-modifier
description: Modify Razvan's personalized suckless dwm setup safely, including config changes, native C patches, dropdown features, X11 behavior, build/install workflow, and patch hygiene. Use when changing /home/razvan/.config/dwm files such as config.h, config.def.h, dwm.c, dropdown.c, CONTEXT.md, ADRs, or when debugging dwm runtime behavior, WM_CLASS rules, keybindings, multi-monitor behavior, focus, tags, floating, or restart/install issues.
---

# DWM Modifier

## Quick start

Work in `/home/razvan/.config/dwm`.

```sh
git -C /home/razvan/.config status --short dwm
make -C /home/razvan/.config/dwm clean && make -C /home/razvan/.config/dwm
```

Do **not** install, kill X, restart dwm, or press restart keybindings unless the user explicitly approves.

## Local model

- Active build config: `config.h`.
- Template/default config: `config.def.h`; keep structurally compatible when changing public config types, but remember existing builds use `config.h`.
- Main source: `dwm.c`.
- Native dropdown implementation: `dropdown.c`, included from `dwm.c` after `config.h`.
- Domain docs: read `CONTEXT.md` and relevant `docs/adr/*.md` before changing terminology or behavior.
- Session startup: `dwm.desktop` -> `/home/razvan/.config/X11/xinitrc` -> `~/.local/bin/startdwm.sh`.
- Installed binary target: usually `/usr/local/bin/dwm` from `config.mk`.

## Workflow

1. **Inspect before editing**
   - Check git status under `/home/razvan/.config`.
   - Read relevant parts of `config.h`, `config.def.h`, `dwm.c`, included patch files, `CONTEXT.md`, and ADRs.
   - For runtime matching issues, use/ask for `xprop WM_CLASS WM_NAME` evidence.

2. **Classify the change**
   - Config-only: keybindings, rules, commands, layout constants.
   - Template-shape change: update both `config.h` and `config.def.h`.
   - Native feature/patch: keep `dwm.c` hooks minimal; prefer an included focused file like `dropdown.c` for feature logic.
   - Runtime/debug issue: reproduce or get concrete X11/window evidence before hypothesizing.

3. **Edit carefully**
   - Preserve suckless C style: C99, file-local `static`, tabs, declarations before statements, `/* */` comments, simple functions.
   - Avoid broad WM_CLASS rules. Prefer unique classes, e.g. `dropdown-terminal`.
   - Keep one logical change per edit batch when possible.
   - Never regenerate or overwrite `config.h` from `config.def.h` without asking.

4. **Validate**
   - Always run:
     ```sh
     make -C /home/razvan/.config/dwm clean && make -C /home/razvan/.config/dwm
     ```
   - For install preview, use:
     ```sh
     make -C /home/razvan/.config/dwm -n install
     ```
   - Report changed files, validation output, residual risks, and exact runtime test steps.

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

Only after explicit user approval: `cd /home/razvan/.config/dwm && sudo make install`. Then the user can press `MOD+Shift+q`; `startdwm.sh` relaunches only if the installed binary checksum changed.

## Common pitfalls

- Editing only `config.def.h` does not affect active `config.h`.
- External suckless patches often target `config.def.h`; port relevant parts to `config.h` manually.
- Duplicate keybindings in dwm can both execute because keypress iterates matching entries.
- `spawn` does not expand shell `~` unless using `SHCMD` or absolute paths.
- Missing `~/.local/state/dwm` can hide stderr log output.
- Multi-monitor features should use monitor work area (`wx/wy/ww/wh`) unless screen geometry is explicitly intended.
