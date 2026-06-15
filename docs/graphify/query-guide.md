# Graphify query guide for dotfiles

Graphify is a routing layer, not a replacement for reading source files. Use it to find likely files, concepts, and cross-community paths before opening files.

## Query with exact domain terms

Prefer domain terms, file names, and community labels over broad natural-language questions.

Good queries:

```sh
graphify query "Hard-boundary guard dotfiles safety guard"
graphify query "scripts/check-dotfiles-safety"
graphify query "scripts/claude-safety-hook"
graphify query "Agent Config Boundaries"
graphify query "Claude Code Adapter Canonical Cross-Agent Instructions"
graphify path "scripts/claude-safety-hook" "scripts/check-dotfiles-safety"
graphify explain "scripts/check-dotfiles-safety"
```

Avoid broad queries that contain overloaded words such as `setup`, `config`, `run`, or `main` unless combined with exact domain terms. For example, `how does the agent safety setup work?` may match unrelated symbols such as `setup()` in dwm.

## Refinement loop

1. Start with a community name, concept name, or file name.
2. If results include unrelated code, add another exact term.
3. Use `path` when checking how two known things connect.
4. Use `explain` when starting from one file or symbol.
5. Read the returned files before editing.

## Useful community labels

- Agent Config Boundaries
- Graphify Integration
- Playwright MCP Config
- Dotfiles Secret Loader ADR
- Dotfiles Secret Env Contract
- DWM Window Manager Core
- ST Terminal Core
- Keyboard Contract

## Agent safety query anchors

Use these terms for safety work:

- `Hard-boundary guard`
- `Autonomous within hard boundaries`
- `dotfiles safety guard`
- `check-dotfiles-safety`
- `claude-safety-hook`
- `PreToolUse`
- `PostToolUse`
- `bypassPermissions`
- `Canonical Agent Instructions`
- `Claude Code Adapter`
- `Agent configuration boundaries`
