@../agents/AGENTS.md

# Claude Code Adapter

- Use the canonical cross-agent instructions imported above as the source of truth.
- Claude Code may run in bypass-permissions mode only under the autonomous-with-hard-boundaries policy.
- Prefer planning mode or explicit user confirmation before broad, destructive, or cross-repository changes.
- Use Claude hooks/permissions as additive early guards; the repo safety guard remains the final gate for dotfiles commits.
- Keep Claude-generated state, auth, transcripts, caches, plugins, and shell snapshots out of git.
