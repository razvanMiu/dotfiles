# Code Context

## Files Retrieved
1. `/home/razvan/.config/dwm/config.h` (lines 3-110) - active dwm appearance, tags/rules/layouts, launcher/dropdown commands, key/button bindings.
2. `/home/razvan/.config/dwm/dwm.c` (lines 49-55, 69-145, 177-200, 236-297, 1104, 1688-1723, 1771-1788, 2038-2054) - confirms local patch hooks, dropdown client state, status mechanism, protected operations.
3. `/home/razvan/.config/dwm/features/dropdown.c` - native dropdown implementation and behavior.
4. `/home/razvan/.config/dwm/startdwm.sh` (lines 1-18) - session startup/autostart-equivalent commands and dwm restart loop.
5. `/home/razvan/.config/X11/xinitrc` (lines 1-36) - X session entry point used by display manager and startx path.
6. `/home/razvan/.config/dwm/dwm.desktop` (lines 1-8) - XSession desktop file points to X11 xinitrc.
7. `/home/razvan/.config/rofi/config.rasi` (lines 1-5) - rofi default config imports Catppuccin Macchiato but is not what dwm keybinding uses.
8. `/home/razvan/.config/rofi/catppuccin-default.rasi` (lines 1-120) - installed Catppuccin-based rofi theme variant.
9. `/home/razvan/.config/rofi/themes/rounded-nord-dark.rasi` (lines 1-18) - actual rofi theme passed by dwm launcher binding.
10. `/home/razvan/.config/dwm/CONTEXT.md` (lines 1-17) and `docs/adr/0001-native-dropdown-terminal.md` (lines 1-18) - local terminology/design intent for dropdown terminal.

## Key Code

Active dwm visual setup (`config.h:3-18`):
```c
static const unsigned int borderpx  = 1;
static const int showbar            = 1;
static const int topbar             = 1;
static const char *fonts[]          = { "JetBrainsMono Nerd Font:size=10" };
static const char col_base[]        = "#1e1e2e"; /* Catppuccin Mocha */
static const char col_surface0[]    = "#313244";
static const char col_text[]        = "#cdd6f4";
static const char col_mauve[]       = "#cba6f7";
```
- Bar enabled, top aligned, stock dwm bar drawing.
- 1px borders.
- Catppuccin Mocha-ish palette for dwm (`Norm`: text/base/surface0, `Sel`: text/surface0/mauve).
- Font is JetBrainsMono Nerd Font size 10.

Tags/rules/layouts (`config.h:20-43`):
- Tags are plain `"1"` through `"9"`.
- Rules: Brave opens on tag 9; Zed on tag 1; `dropdown-terminal` and `dropdown-test` are floating native dropdowns.
- Layouts: tile `[]=`, floating `><>`, monocle `[M]`; default is tile, `mfact=0.55`, `nmaster=1`, `resizehints=1`, `refreshrate=120`.

Launcher/dropdowns (`config.h:55-69`):
```c
static const char *dmenucmd[] = { "rofi", "-show", "drun", "-theme", "~/.config/rofi/themes/rounded-nord-dark.rasi", NULL };
static const char *droptermcmd[] = {
	"st", "-c", "dropdown-terminal", "-n", "dropdown-terminal", "-T", "dropdown-terminal", "-A", "1.0",
	"-e", "tmux", "new-session", "-A", "-s", "dropdown",
	NULL
};
```
- Mod+Space runs the XDG rofi launcher script.
- Mod+grave toggles the dropdown st/tmux session; Mod+n toggles dropdown notes.

Patches/local changes already present:
- Native dropdown patch: `Client.dropdown`, `Rule.dropdown`, `Monitor.dropw/droph`, `MAXDROPDOWNS`, declarations in `dwm.c`, and `#include "features/dropdown.c"`.
- Dropdown behavior in `features/dropdown.c`: sticky while visible via `TAGMASK`, hidden by setting `tags=0`, remembers per-monitor size during current dwm runtime, moves to selected monitor, raises above normal windows, prevents duplicate spawn for 2s.
- Focus-follows-mouse `enternotify` is commented out in `dwm.c:781-798`; focus is click/keyboard driven.
- EWMH `_NET_CLIENT_LIST` support is present (`NetClientList` atom in `dwm.c:65-68` and setup/update hooks elsewhere).
- No systray, alpha, gaps, vanity layouts, pertag, swallow, xresources, autostart patch, or statuscmd patch found in searched code.

Status mechanism (`dwm.c:2048-2054`):
```c
if (!gettextprop(root, XA_WM_NAME, stext, sizeof(stext)))
    strcpy(stext, "dwm-"VERSION);
drawbar(selmon);
```
- Status text is root window name, i.e. set by `xsetroot -name ...` or a status daemon/script.
- No local `slstatus`, `dwmblocks`, `xsetroot` loop, or status script was found under `/home/razvan/.config`; `command -v` only found rofi among checked tools (`picom`, `feh`, `nitrogen`, `xwallpaper`, `xsetroot`, `slstatus`, `dwmblocks`, `rofi`). Note: this PATH inspection reported `/usr/sbin/rofi`; it did not report picom/wallpaper/status tools.

Startup/wallpaper/compositor:
- X session: `dwm.desktop:1-8` -> `Exec=/home/razvan/.config/X11/xinitrc`; `xinitrc:1-36` sets XDG/session/dbus/keyring environment then execs `$HOME/.local/bin/startdwm.sh`.
- `/home/razvan/.local/bin/startdwm.sh` is a symlink to `/home/razvan/.config/dwm/startdwm.sh`.
- `startdwm.sh:3-4` runs only `xrandr --output DP-2 --mode 1920x1080 --rate 144` and `setxkbmap us -option ctrl:nocaps,ctrl:swap_lalt_lctl` before the dwm restart loop.
- No wallpaper command/config found (`feh`, `nitrogen`, `xwallpaper`, `wall*` search empty under `.config`).
- No picom config found under `.config`, and `command -v picom` returned nothing in this environment.

Rofi setup:
- Dwm binding uses `/home/razvan/.config/rofi/themes/rounded-nord-dark.rasi:1-18`, a Nord rounded theme (`#2E3440F2`, `#88C0D0F2`, etc.).
- Rofi default config imports Catppuccin Macchiato (`/home/razvan/.config/rofi/config.rasi:1-5`) and there are Catppuccin theme files, but these are bypassed by the dwm keybinding.
- Visual mismatch risk: dwm is Catppuccin Mocha, rofi launcher is Nord.

## Architecture

Display manager/startx launches `dwm.desktop`, which executes `/home/razvan/.config/X11/xinitrc`. That script prepares session/dbus environment and execs `~/.local/bin/startdwm.sh` (symlink to this repo). `startdwm.sh` applies monitor and keyboard setup, then runs dwm in a checksum loop so rebuilding/reinstalling the `dwm` binary can relaunch it once.

Dwm visuals are compile-time config in `config.h`; changing colors/fonts/layouts requires rebuilding/reinstalling dwm. The bar is the built-in dwm bar: tags/layout/window title/status. Status is external via root window name; no external status producer was located locally. Dropdown terminal behavior is a local native patch, not stock scratchpad/tdrop, and has explicit terminology in `CONTEXT.md`.

## Start Here

Open `/home/razvan/.config/dwm/config.h` first. It contains the active user-facing visual choices and commands; then inspect `startdwm.sh` for low-risk startup additions and `features/dropdown.c` or `features/status.c` only if touching native feature behavior.

## Local implications and low-risk improvement candidates

- Align rofi theme with dwm: either change `dmenucmd` to a Catppuccin Mocha/Macchiato theme already present, or adjust `rounded-nord-dark.rasi`; current launcher is Nord while dwm is Catppuccin.
- Add startup commands cautiously in `startdwm.sh` rather than implementing dwm autostart patch: e.g. guarded background launches for compositor/wallpaper/status once tools are installed.
- Picom is configured under `~/.config/picom/picom.conf`; dropdown st windows are launched opaque (`st -A 1.0`) and picom owns their translucency/blur via dropdown WM_CLASS rules.
- If wallpaper is desired, choose one tool (`feh`/`xwallpaper`/etc.) and add a single command to `startdwm.sh`; no current wallpaper mechanism exists locally.
- If a status bar is desired, add a simple root-name loop or `slstatus`/`dwmblocks` in startup; dwm already consumes `XA_WM_NAME`, so no dwm patch is needed for basic status.
- Avoid overwriting dropdown semantics with a generic scratchpad patch; local code and docs intentionally implement custom native dropdown behavior.
