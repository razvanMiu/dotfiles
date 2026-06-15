# Agent Workflows in This Dotfiles Repo

The reusable workflow router lives as a shared skill:

```text
~/.config/agents/skills/ai-workflow/SKILL.md
```

Use that skill for non-trivial coding, debugging, architecture, UI, or research tasks. It routes work through the rest of the local AI stack:

```text
Graphify  -> project architecture and file relationship routing
Serena    -> symbol-aware code navigation and semantic edits
Context7  -> fresh library/API documentation
Playwright -> browser/UI evidence
NotebookLM -> videos, papers, large-source research
subagents/forks -> noisy scouting, review, diagnosis, risk checks
specialized skills -> diagnose, tdd, grill-with-docs, prototype, triage, etc.
```

## How to invoke it

Ask the agent explicitly:

```text
Use the ai-workflow skill. Follow the full workflow for this task.
```

Or choose a branch:

```text
Use ai-workflow and the UI/browser workflow. Validate with Playwright.
```

```text
Use ai-workflow for a library/framework change. Use Context7 before editing.
```

```text
Use ai-workflow for architecture work. Use Graphify first and grill-with-docs if terminology or ADRs are involved.
```

## Dotfiles-specific finalization

For this repo only, finalization includes the normal project checks plus the dotfiles safety guard:

```sh
cd ~/.config
scripts/check-dotfiles-safety --quiet
git diff --check
```

The safety guard is repo-specific. Other projects should run their own tests, type checks, linters, builds, and validation commands instead.

## Graphify refresh policy for this repo

Do not run full Graphify after every edit.

- Use `graphify query ...` freely.
- Use `graphify update` after meaningful code/doc changes.
- Use full `/graphify . --wiki --no-viz` only after large architectural changes, major refactors, or before major review.
- Refresh Obsidian/community labels only when the visual/wiki output matters.
