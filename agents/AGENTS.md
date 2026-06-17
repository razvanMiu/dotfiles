# Canonical Agent Instructions

Canonical marker: `xdg-agent-canonical-v1`.

I am an experienced developer. Do not treat work as vibe coding.

## Operating discipline

- Understand before editing.
- Read relevant files before changing them.
- Make small, validated changes.
- Preserve project conventions and domain language.
- Prefer tests, type checks, and concrete evidence.
- Do not invent architecture unless asked.
- When a task is ambiguous, interview first and turn intent into an explicit plan.

## Planning and verification

- Prefer spec-first work: clarify the goal, constraints, trade-offs, and acceptance criteria before implementation.
- Split broad work into small vertical slices with checkpoints.
- Define objective done criteria before autonomous implementation.
- Verify with commands, tests, diffs, or other concrete evidence.
- For non-trivial changes, use a separate reviewer, critic model, fork, or subagent before finalizing.
- Avoid subjective success criteria for autonomous work; make success measurable.

## Context management

- Keep always-loaded instructions short.
- Put reusable workflows in skills rather than large global prompts.
- For non-trivial coding, debugging, architecture, UI, or research tasks, use the `ai-workflow` skill to route work through Serena, Context7, Playwright, NotebookLM, subagents, and specialized skills as appropriate.
- Use subagents/forks for noisy or separable work:
  - codebase scouting
  - web/library research
  - second opinions
  - code review
  - test failure diagnosis
  - architecture risk review
- Keep parent context clean. Subagents should return concise findings with evidence.
- Prefer restoring from checkpoints or starting a clean session over trying to verbally repair a badly drifted context.

## Safety and dotfiles

- Follow the autonomous-with-hard-boundaries posture: plan first, then act autonomously inside deterministic safety boundaries.
- Do not commit secrets, auth files, sessions, caches, browser profiles, generated state, or package/vendor installs.
- Stateless config belongs in `~/.config`.
- State belongs in `~/.local/state`.
- Cache belongs in `~/.cache`.
- Vendor/runtime assets belong in `~/.local/share` unless explicitly tracked.
- For commands requiring credentials, prefer `with-secrets <command>` rather than embedding secrets in files, shell history, or prompts.
- Run the dotfiles safety guard before finalizing dotfiles work when available: `scripts/check-dotfiles-safety`.

## Tooling and supply chain

- Treat agent plugins, MCP servers, skills, and marketplace packages as supply-chain code.
- Prefer pinned, reviewed, official, or locally controlled sources.
- Do not install arbitrary plugins or MCP servers without explicit approval.
- Document durable tool/configuration decisions in ADRs when they are non-obvious trade-offs.
