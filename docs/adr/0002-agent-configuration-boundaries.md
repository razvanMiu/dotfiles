---
status: accepted
---

# Agent configuration boundaries

The dotfiles repo uses `~/.config/agents` as the canonical cross-agent layer for shared instructions and skills, while Pi and Claude Code keep adapter-specific configuration under `~/.config/pi/agent` and `~/.config/claude`. Claude Code is treated as a mixed tool directory because it does not cleanly separate authored config from generated state; generated auth, cache, transcripts, package/vendor assets, and browser-like profiles are ignored or moved to XDG state/data/cache only where the tool exposes supported paths. For Pi, `PI_CODING_AGENT_DIR` is set to `$XDG_DATA_HOME/pi/agent` and `PI_CODING_AGENT_SESSION_DIR` to `$XDG_STATE_HOME/pi/sessions`, with state and cache linked into `$XDG_STATE_HOME/pi` and `$XDG_CACHE_HOME/pi`. Authored declarative configs (`settings.json`, `mcp.json`, `extensions/`, `AGENTS.md`) remain cleanly tracked under `~/.config/pi/agent` and symlinked into `$XDG_DATA_HOME/pi/agent`, keeping dotfiles free of npm packages, binaries, caches, and live credentials. Coding agents use an "autonomous within hard boundaries" posture: Claude may use `permissions.defaultMode = "bypassPermissions"` and Pi may keep broad tool access, but repo-local safety guards must block non-negotiable risks across staged, tracked, and untracked non-ignored dotfiles files. Claude Code also uses PreToolUse/PostToolUse hooks as additive early guards; the repo-local pre-commit guard remains the final mandatory boundary.

## Considered Options

- Make Claude Code the source of truth for agent resources; rejected because shared instructions and skills should not be coupled to one tool.
- Force every tool-generated file out of `~/.config` with symlinks; rejected because unsupported relocation can break tool updates and corrupt state.
- Use approval-first agent permissions everywhere; rejected because the intended workflow is planned autonomous execution with deterministic guardrails rather than repeated prompt-level approval.
