---
status: accepted
---

# Antigravity CLI XDG layout

The Antigravity CLI (`agy`) hardcodes `$HOME/.gemini` as its single runtime root and honours no override environment variable, so the dotfiles keep the binary at `$XDG_BIN_HOME/agy` and relocate its state by symlinking `~/.gemini` to `$XDG_DATA_HOME/antigravity`. The root holds both declarative config (`config/config.json`, `config/mcp_config.json`, `config/projects/`) and generated runtime state (`antigravity-cli/conversations/`, `brain/`, `log/`, `crashes/`, `cache/`, `updater/`, `builtin/`, `conversation_summaries.db*`, `installation_id`, `jetski_state.pbtxt` with mode `0600`), which is why the whole tree is treated as state rather than splitting config into `$XDG_CONFIG_HOME`. Authentication lives in the macOS Keychain, not in the tree, so relocating it does not touch credentials. The symlink is created after moving, never in place of a real directory, and `mv` on the same APFS volume preserves modes and timestamps so the `0600` state file survives intact. `~/.local/bin/agy` reaches `PATH` through `$XDG_BIN_HOME` in `zsh/paths.zsh`, so the installer's profile edits are redundant and removed.

## Considered Options

- Set an override environment variable such as `AGY_HOME`, `AGY_APP_DATA`, `GEMINI_DIR`, or `ANTIGRAVITY_CONFIG_DIR`; rejected because every candidate was set simultaneously against v1.5.9 and none created a directory or redirected a write, and upstream still tracks the gap in google-antigravity/antigravity-cli#155.
- Point the tool at a different directory using `agy install --dir`; rejected because that flag only configures `PATH`, it does not relocate the runtime root.
- Leave `~/.gemini` in `$HOME`; rejected because it is generated state, and the dotfiles keep state in `$XDG_STATE_HOME` and `$XDG_DATA_HOME` rather than scattering dot-directories in the home directory.
- Symlink `~/.gemini` to `$XDG_CONFIG_HOME/antigravity`; rejected because the tree is dominated by conversations, logs, SQLite databases, and caches, so treating it as configuration would put writable runtime state inside `~/.config`.
- Split config from state with two symlinks; rejected because upstream resolves a single root and no supported mechanism separates the two subtrees.

## Consequences

`~/.gemini` must stay a symlink for the tool to work. A fresh reinstall hardcodes `$HOME/.cache/antigravity/staging` again and can recreate or replace the root, so the symlink is re-applied after any reinstall and the stray staging directory is deleted. `agy update` self-updates the binary in place and any installer re-run is invoked as `agy install --skip-path --skip-aliases` to avoid re-appending `PATH` exports to `~/.zshrc` and `~/.zprofile`. The official Gemini CLI shares `~/.gemini` by design, so installing it alongside this CLI would make both tools write into the same tree; that collision is inherent to upstream and not introduced by the symlink. When upstream ships a working root override, the symlink is retired in favour of the variable.
