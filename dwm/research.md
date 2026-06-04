# Research: tasteful lightweight visual improvements for suckless dwm

## Summary
The best lightweight path is to keep dwm’s core bar and add a small set of official suckless patches: `status2d` or `statuscolors` for colored status text, `alpha` for controlled bar transparency, `barpadding` for breathing room, and optionally `vanitygaps` if gaps are desired. Use Catppuccin only as a static color palette in `config.h`/status scripts, prefer `slstatus` or a minimal `dwmblocks` over heavyweight bars, and run picom conservatively for shadows/tear-free compositing rather than blur-heavy effects.

## Findings
1. **Most tasteful patch set: status2d + alpha + barpadding; optional vanitygaps.** `status2d` adds colors and rectangle drawing to dwm’s native status bar, avoiding an external panel; `alpha` makes the bar translucent while keeping text opaque and can also keep borders opaque under a compositor; `barpadding` adds small vertical/side spacing; `vanitygaps` explicitly describes itself as pure visual “eyecandy,” so use small gaps only if the look is worth the screen loss. [status2d](https://dwm.suckless.org/patches/status2d/) [alpha](https://dwm.suckless.org/patches/alpha/) [barpadding](https://dwm.suckless.org/patches/barpadding/) [vanitygaps](https://dwm.suckless.org/patches/vanitygaps/)
2. **Status bar: stay native unless you need modules/clicks.** The suckless status monitor page lists `slstatus` as the canonical suckless status monitor and `dwmstatus` as a barebones C approach; `slstatus` supports common modules including battery, CPU, disk, entropy, shell commands, and date/time. For clickable/modular blocks, `dwmblocks-async` is a reputable lightweight option, but adds more moving parts than `slstatus`. [suckless status monitor](https://dwm.suckless.org/status_monitor/) [slstatus](https://codeberg.org/beastie/slstatus) [dwmblocks-async](https://github.com/DusanLesan/dwmblocks)
3. **Catppuccin integration should be palette-level, not a full rice.** There does not appear to be an official maintained `catppuccin/dwm` port in the Catppuccin ecosystem; official Catppuccin docs emphasize ports/palette usage broadly, while dwm examples are mostly personal forks/rices. Recommendation: copy Catppuccin Mocha colors into dwm `colors[][]`, dmenu/st config, status script color codes, GTK/terminal/rofi themes, and avoid importing a whole third-party dwm tree. [Catppuccin ports](https://catppuccin.com/ports/) [Catppuccin main repo](https://github.com/catppuccin/catppuccin) [example third-party dwm rice](https://github.com/mozkomor05/dwm)
4. **Fonts: Nerd Font + emoji fallback gives polish without runtime cost.** Reputable dwm examples commonly use JetBrainsMono Nerd Font plus Noto Color Emoji; this gives icons/glyphs for status blocks while staying just an Xft font list in config. Tradeoff: Nerd Font icons can make scripts less portable if the font is missing. [example dwm build](https://github.com/mozkomor05/dwm)
5. **Compositor: picom is appropriate, but keep effects restrained.** picom is a standalone Xorg compositor for WMs without compositing and supports opacity, shadows, fading, and blur. Its own man page warns background blur is “bad in performance” and driver-dependent; ArchWiki also notes dwm-specific opacity/exclusion caveats. Recommendation: enable vsync/backend as needed, subtle shadows/fading, maybe inactive dim/opacity; avoid blur/animations if “lightweight” is the priority. [picom man page](https://picom.app/) [ArchWiki picom](https://wiki.archlinux.org/title/Compton)
6. **Other visual patches worth considering.** `statuscolors` is simpler than `status2d` if only colored text is needed; `scheme_switch` allows multiple color schemes such as dark/light switching; `systray` is practical for tray icons but adds complexity and has patch interactions with barpadding/status2d. [statuscolors](https://dwm.suckless.org/patches/statuscolors/) [scheme_switch](https://dwm.suckless.org/patches/scheme_switch/) [systray](https://dwm.suckless.org/patches/systray/)

## Recommended options
- **Minimal polished:** Catppuccin Mocha colors in `config.h`, JetBrainsMono Nerd Font + Noto Color Emoji, `slstatus`, no compositor or picom with only vsync/shadows.
- **Best balanced:** `status2d`, `alpha`, `barpadding`, small `vanitygaps`, Catppuccin palette, `slstatus` or minimal `dwmblocks-async`, picom with subtle shadows and no blur.
- **Avoid for lightweight dwm:** polybar/eww-heavy setups, full third-party dwm rices, blur-heavy picom forks/configs, large patch stacks that fight during upgrades.

## Sources
- Kept: suckless `status2d` (https://dwm.suckless.org/patches/status2d/) — official patch for colored/native bar rendering.
- Kept: suckless `alpha` (https://dwm.suckless.org/patches/alpha/) — official patch for translucent bar and opaque text/borders.
- Kept: suckless `barpadding` (https://dwm.suckless.org/patches/barpadding/) — official patch for tasteful spacing.
- Kept: suckless `vanitygaps` (https://dwm.suckless.org/patches/vanitygaps/) — official patch for gaps, with clear eyecandy tradeoff.
- Kept: suckless status monitor page (https://dwm.suckless.org/status_monitor/) — authoritative status-bar ecosystem overview.
- Kept: slstatus Codeberg (https://codeberg.org/beastie/slstatus) — canonical suckless-style status monitor.
- Kept: picom man page (https://picom.app/) — primary evidence for compositor features/performance warnings.
- Kept: Catppuccin ports/main repo (https://catppuccin.com/ports/, https://github.com/catppuccin/catppuccin) — authoritative palette/port ecosystem.
- Dropped: random Catppuccin dwm rices — useful inspiration but not authoritative and often include heavyweight dependencies.
- Dropped: old compton/xcompmgr commentary — superseded by picom docs and ArchWiki.

## Gaps
No official Catppuccin dwm port was confirmed; next step would be manually mapping the Catppuccin palette into the local `config.h`, dmenu/st, and status scripts. Patch compatibility depends on the current dwm version and existing local patches, so inspect current patch stack before applying anything.
