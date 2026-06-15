# Safety guard scripts

This document exists to give Graphify exact routing anchors for the dotfiles safety scripts.

## scripts/check-dotfiles-safety

`scripts/check-dotfiles-safety` is the mandatory dotfiles safety guard. It checks staged files, tracked files, and untracked non-ignored files. It blocks forbidden auth, session, state, cache, package, vendor, browser profile, generated Graphify output, and high-confidence secret content before commit.

Use this exact Graphify query anchor when investigating the safety guard:

```sh
graphify query "scripts/check-dotfiles-safety"
```

Related concepts:

- `dotfiles safety guard`
- `Hard-boundary guard`
- `Current-file guard scope`
- `git/hooks/pre-commit`

## scripts/claude-safety-hook

`scripts/claude-safety-hook` is the Claude Code early guard. It is configured by `claude/settings.json` as a `PreToolUse` and `PostToolUse` hook.

The hook blocks dangerous shell commands, `git commit --no-verify`, forced staging of ignored files, and writes to forbidden dotfiles paths. After mutating tools inside the dotfiles repo, the post-tool path runs `scripts/check-dotfiles-safety --quiet`.

Use this exact Graphify query anchor when investigating the Claude guard:

```sh
graphify query "scripts/claude-safety-hook"
```

Related concepts:

- `Claude safety hook`
- `PreToolUse`
- `PostToolUse`
- `bypassPermissions`
- `claude/settings.json`
