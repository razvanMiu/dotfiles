# Agent configuration boundaries

This repo separates tool-neutral agent resources from tool-specific adapter configuration.

## Layers

- `Cross-agent layer`: tracked tool-neutral agent resources under `agents/`.
- `Shared skill directory`: `agents/skills/`, also exposed as `~/.agents/skills` and `claude/skills`.
- `Pi adapter`: `pi/agent/`, with `pi/agent/AGENTS.md` symlinked to `../../agents/AGENTS.md`.
- `Claude Code adapter`: `claude/`, with `claude/CLAUDE.md` importing `@../agents/AGENTS.md`.
- `Mixed tool directory`: tool config directories that contain generated state because the tool does not fully separate config, state, and cache.

## Canonical instruction loading

The canonical instruction file is:

```text
agents/AGENTS.md
```

Pi loads it through:

```text
pi/agent/AGENTS.md -> ../../agents/AGENTS.md
```

Claude Code loads it through:

```text
claude/CLAUDE.md -> @../agents/AGENTS.md
```

## Generated data policy

Tracked config belongs in `~/.config`. State belongs in `~/.local/state`. Cache belongs in `~/.cache`. Package/vendor assets belong in `~/.local/share` only when the tool supports relocation.

Pi currently installs user packages under `PI_CODING_AGENT_DIR/npm` and `PI_CODING_AGENT_DIR/git`, so those directories stay ignored rather than using unsupported relocation.

Claude Code stores some auth, transcript, cache, and generated state inside `CLAUDE_CONFIG_DIR`, so `claude/` is treated as a mixed tool directory with strict ignore rules.

## Related files

- `agents/AGENTS.md`
- `agents/skills/`
- `pi/agent/settings.json`
- `pi/agent/AGENTS.md`
- `claude/CLAUDE.md`
- `claude/settings.json`
- `docs/adr/0002-agent-configuration-boundaries.md`
- `CONTEXT.md`
