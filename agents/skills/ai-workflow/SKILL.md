---
name: ai-workflow
description: Routes non-trivial coding, debugging, architecture, research, and UI tasks through Razvan's full AI stack: Graphify, Serena, Context7, Playwright, NotebookLM, subagents, and specialized skills. Use when the user asks for the full workflow, asks how to approach work, starts a non-trivial code task, or wants coordinated tool/skill usage.
---

# AI Workflow

Use this as the orchestrator skill. It does not replace specialized skills; it decides when to load them and which MCP/tool evidence is needed.

## Default flow

1. Clarify objective, constraints, and measurable acceptance criteria.
2. Classify the task using the routing table below.
3. Load the specialized skill(s) only when the task matches.
4. Use the cheapest routing tool first, then narrow to files/symbols/docs/tests.
5. Implement in a small vertical slice.
6. Verify with concrete evidence.
7. Finalize with project-specific checks and a concise summary.

## Routing table

| Task type | First move | Then use | Optional skill |
| --- | --- | --- | --- |
| Architecture, file relationships, where to look | `graphify query` when `graphify-out/graph.json` exists | `graphify path`, `graphify explain`, then read targeted files | graphify |
| Symbol/code change | Serena symbol overview | `find_symbol`, `find_referencing_symbols`, focused edits/tests | tdd |
| Library/framework/API change | Context7 docs | code_search/upstream source, tests/type checks | librarian |
| UI/browser behavior | Playwright reproduce/snapshot | console/network/screenshot, fix, Playwright proof | diagnose |
| Hard bug/regression | Reproduce/minimize | hypothesis/instrument/fix/regression test | diagnose |
| Test-first feature/fix | Write failing test | red-green-refactor | tdd |
| Architecture/design decision | Interview and inspect docs | update CONTEXT.md/ADR if warranted | grill-with-docs |
| Video/paper/large-source research | NotebookLM notebook/sources | synthesize into docs/issues | notebooklm |
| Prototype/state-machine/UI options | Throwaway prototype | compare options, discard or promote | prototype |
| Planning into PRD/issues | Clarify scope | publish PRD/issues | to-prd, to-issues, triage |

## Tool order rules

- If `graphify-out/graph.json` exists and the question is about project structure, query Graphify before broad `rg`/`find`.
- Use Serena after Graphify identifies likely code files, or directly for known symbols.
- Use Context7 before changing fast-moving external APIs.
- Use Playwright when the success criterion is visual, interactive, browser, console, or network behavior.
- Use subagents/forks for noisy scouting, independent review, test failure diagnosis, and architecture risk review.

## Finalization

For any repo:

```sh
git diff --check
```

Also run relevant project validation such as tests, type checks, lint, build, Playwright checks, or documented commands.

For the dotfiles repo rooted at `~/.config`, additionally run when available:

```sh
scripts/check-dotfiles-safety --quiet
```

Do not commit secrets, auth files, sessions, caches, generated state, browser profiles, or package/vendor installs.

## Graphify refresh policy

- `graphify query ...`: cheap; use anytime.
- `graphify update`: use after meaningful code/doc changes when the graph should stay current.
- Full `/graphify . --wiki --no-viz`: use after architecture-scale changes, large refactors, or before a major review.
- Obsidian export/community relabel: only when the visual/wiki output matters.

Do not run full Graphify after every small edit.
