# Agent safety architecture

The dotfiles repo uses an autonomous-with-hard-boundaries posture for coding agents.

## Core concepts

- `Autonomous within hard boundaries`: agents may act autonomously after planning, but deterministic guards block non-negotiable risks.
- `Hard-boundary guard`: deterministic safety checks independent of prompts or model behavior.
- `dotfiles safety guard`: repo-local safety command at `scripts/check-dotfiles-safety`.
- `Claude safety hook`: Claude Code early guard at `scripts/claude-safety-hook`.

## Safety layers

1. Canonical instructions in `agents/AGENTS.md` define operating discipline.
2. Claude adapter instructions in `claude/CLAUDE.md` import canonical instructions and add Claude-specific safety expectations.
3. `claude/settings.json` sets `permissions.defaultMode` to `bypassPermissions` and relies on deterministic hooks plus repo guards.
4. `scripts/claude-safety-hook` runs as Claude Code `PreToolUse` and `PostToolUse` hook.
5. `scripts/check-dotfiles-safety` runs as the mandatory repo guard and pre-commit hook via `git/hooks/pre-commit`.

## Guard responsibilities

`scripts/claude-safety-hook` blocks early:

- dangerous broad shell commands
- `git commit --no-verify`
- forced staging of ignored files
- writes to forbidden auth, session, state, cache, package, vendor, browser, or generated paths

`scripts/check-dotfiles-safety` checks staged files, tracked files, and untracked non-ignored files by default. It blocks forbidden paths and high-confidence secret content. Ignored files are excluded from blocking checks unless explicitly audited.

## Important files

- `agents/AGENTS.md`
- `claude/CLAUDE.md`
- `claude/settings.json`
- `scripts/claude-safety-hook`
- `scripts/check-dotfiles-safety`
- `git/hooks/pre-commit`
- `docs/adr/0002-agent-configuration-boundaries.md`
- `CONTEXT.md`
