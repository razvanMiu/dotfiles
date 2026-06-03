# Global Agent Instructions

I am an experienced developer. Do not treat work as vibe coding.

Engineering discipline:

- Understand before editing.
- Read relevant files before changing them.
- Make small, validated changes.
- Preserve project conventions.
- Prefer tests, type checks, and concrete evidence.
- Do not invent architecture unless asked.

Use subagents for noisy or separable work:

- codebase scouting
- web/library research
- second opinions
- code review
- test failure diagnosis
- architecture risk review

Keep parent context clean. Subagents should return concise findings with evidence.

Safety and dotfiles:

- Do not commit secrets, auth files, sessions, caches, browser profiles, generated state, or package installs.
- Stateless config belongs in `~/.config`.
- State belongs in `~/.local/state`.
- Cache belongs in `~/.cache`.
- Vendor/runtime assets belong in `~/.local/share` unless explicitly tracked.
- For commands requiring credentials, prefer `with-secrets <command>` rather than embedding secrets in files, shell history, or prompts.
