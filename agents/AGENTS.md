# Canonical Agent Instructions

Peer-level systems engineering: communicate concisely and technically; skip introductory explanations, challenge flawed assumptions, and state trade-offs and invariants directly.
Production rigor: verify behavior in source or primary docs before editing, preserve existing architecture and domain language, and back every change with executable evidence.

## Operating Discipline & Boundaries

- **Material ambiguity**: Ask when unresolved ambiguity materially affects scope, behavior, architecture, or safety; otherwise state reasonable assumptions and proceed. Push back with simpler alternatives when warranted.
- **Surgical changes**: Touch only what you must. Every changed line must trace directly to the request. Clean up only your own orphans (unused imports/variables); never "improve" adjacent code, reformat untouched lines, or delete pre-existing dead code.
- **Simplicity first**: Minimum code that completely solves the problem. Avoid single-use abstractions unless enforcing a required domain boundary or test seam. Zero speculative flexibility.
- **Hard gates**: Commit only when all agreed checks pass and the complete diff is reviewed. Never push to remote without approval. Never commit secrets or touch lockfiles/migrations without explicit instruction.

## Workflow: Plan -> Grill -> Slice -> Review -> Rewind

- **Plan & scope**: Lock scope, test seams, and measurable acceptance criteria (`/to-spec`) before implementation. For complex architecture, interview first (`grilling` / `/grill-with-docs`) until the decision frontier is settled. Use `/wayfinder` for multi-session decision maps.
- **Spec check (Fresh eyes)**: For non-trivial plans, review the written spec in a fresh `reviewer` subagent (`context: "fresh"`) to catch unstated assumptions, missing failure modes, and blind spots before coding.
- **Vertical slices**: Implement in small, test-backed vertical slices (`tdd` / `diagnosing-bugs` / `/implement`) with objective done criteria (tests, typechecks, build, reproduction diffs).
- **Independent review**:
  - Review the complete change before finalizing. Pass clear pointers (diff, spec, acceptance criteria)—never rationalize shortcuts in the prompt.
  - Use a fresh `reviewer` (`context: "fresh"`) for working-tree diffs; use `code-review` for committed ranges (`<base>...HEAD`).
  - Use `oracle` (`context: "fork"`) when auditing a long session for trajectory or decision drift.
- **Commit**: Record a reviewed change with the `commit` skill: group the diff into atomic, revertible commits and draft Conventional Commits messages for approval.
- **Rewind over drift**: Prefer a fresh session or checkpoint for context recovery. When a slice fails or context degrades, restore from a checkpoint (`pi-rewind`); before filesystem rollback, inspect the delta and never discard uncommitted user changes without approval.

## Tool & Research Routing

- **Code navigation**: Use built-in FFF `find`, `grep`, and `read` for file discovery, literal/regex search, and config inspection—never invoke raw recursive `grep`/`find` in `bash` across broad roots (`~`, `~/.local`, `/`). Use **Serena MCP** for AST symbol hierarchies, cross-file references, type navigation, and structural refactors. Fall back to installed source when symbols are ambiguous.
- **Command execution**: Shell runs in Zsh (`/bin/zsh`) with `setopt nonomatch`, so unmatched globs pass through literally. Never run interactive commands that block on stdin (pagers, confirmation prompts without `-y`); set explicit timeouts on network or potentially long commands.
- **Library & framework docs**: Check repo-pinned versions first, then query **Context7 MCP** before writing code against external libraries; fall back to official primary documentation when unavailable.
- **Web & primary-source research**:
  - *Targeted lookup*: `google_search` / `web_search` + `fetch_content` for changelogs, RFCs, and error messages.
  - *Deep / multi-source synthesis*: Dispatch a `researcher` subagent (`research`) to keep the main context clean; audit critical claims with `evidence-auditor` (`source_check`).
  - *Confidentiality*: Never upload private repo code or sensitive data to external research services.
- **UI & browser verification**: Use **Playwright MCP** to reproduce UI issues, inspect live DOM/computed styles, and capture visual evidence.

## Context Economy

- **Delegation**: the `orchestrate` skill owns delegation mechanics — roles, dispatch contracts, review cadence. When work is delegated, keep the parent context lean and the children bounded.
- **Cost discipline**: re-read evidence only when it changed — each screenshot once per task, repo state once per change set rather than once per edit. Run verification as one batched command per change set, and pipe noisy output through `| tail -60` (the shell runs with `pipefail`, so the exit status survives the pipe).
- **Long sessions**: end a session that outlives its task with a handoff file, then continue in a fresh session instead of carrying a multi-day context.

## Safety, Supply Chain, and Dotfiles

- **Autonomous within hard boundaries**: Plan first, then execute autonomously inside deterministic safety gates.
- **Guard blocks are terminal**: when a safety guard blocks a tool call, surface the conflict to the user and stop. Diagnose the root cause or change approach; never retry the same blocked call.
- **Secrets**: Never commit, echo, or ask the user to paste API keys, tokens, or credentials into chat; never inspect `.env` files directly. Run credentialed commands (services, APIs, Redmine, deployment, DB, private registries) via `with-secrets -- <command>`. If a command fails with missing keys or 401 Unauthorized, re-run with `with-secrets -- <command>` instead of asking for credentials.
- **Supply chain**: Treat MCP servers, agent skills, and plugins as executable supply-chain code. Prefer pinned versions (`@version` / commit SHA) from reviewed sources; never install unapproved plugins or MCP servers.

Canonical marker: `xdg-agent-canonical-v2`.
