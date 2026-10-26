# Zed Editor Configuration

Structured, reproducible Zed configuration optimized for full-stack TypeScript (React & Vue), Python, Go, and Tailwind CSS, with seamless support for both ESLint/Prettier and Biome.

## Architecture

- **`settings.json`**: Primary editor configuration. Declares required extensions via `auto_install_extensions`, language-specific formatting/linting rules, and UI preferences.
- **`keymap.json`**: User keymap customizations.
- **`templates/biome.json`**: Portable `.zed/settings.json` template for Biome-based projects.
- **`prompts/`**: Ignored local prompt library state.

## Formatting & Linting Coexistence

### Default (ESLint + Prettier)

By default, all JavaScript, TypeScript, TSX, and Vue files run **Prettier** on save and trigger `source.fixAll.eslint`. Biome is explicitly excluded (`!biome`) in global settings so it never conflicts with Prettier or throws errors in non-Biome repos.

### Biome Projects

To enable Biome for a project:

```bash
zed-enable-biome
# or manually:
mkdir -p .zed && cp "$XDG_CONFIG_HOME/zed/templates/biome.json" .zed/settings.json
```

This local override switches formatting to Biome LSP, enables Biome fix-all and import sorting on save, and disables ESLint (`!eslint`).

`.zed/` is globally ignored in `~/.config/git/ignore`, so these workspace settings will not dirty git status unless explicitly tracked.

## Stack Summary

| Language / Tool       | Language Server(s)            | Formatter / Linter on Save                           | Indentation        |
| :-------------------- | :---------------------------- | :--------------------------------------------------- | :----------------- |
| **JS / TS / TSX**     | `vtsls`, `eslint`             | Prettier + ESLint Fix                                | 2 spaces           |
| **Vue**               | `vue`, `vtsls`, `eslint`      | Prettier + ESLint Fix                                | 2 spaces           |
| **Tailwind CSS**      | `tailwindcss-language-server` | Built-in (with `clsx`, `cva`, `twMerge`, `cn`, `cx`) | -                  |
| **Python**            | `basedpyright`, `ruff`        | Ruff (format, fixAll, organizeImports)               | 4 spaces           |
| **Go**                | `gopls`, `golangci-lint`      | `gopls` (format, organizeImports)                    | 4 tabs (hard tabs) |
| **HTML / CSS / JSON** | Built-in                      | Prettier / Built-in                                  | 2 spaces           |

## Reviewing Agent Changes (Native Zed Shortcuts)

- **Toggle Git Panel**: `Ctrl + Shift + G` (or click git branch in status bar)
- **Next Diff Hunk**: `Cmd + F8`
- **Previous Diff Hunk**: `Cmd + Shift + F8`
- **Toggle Inline Diff Hunk**: `Cmd + '`
- **Restore / Revert Hunk**: `Cmd + Alt + Z` (while focused on a diff)
