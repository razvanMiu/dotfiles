---
name: ai-workflow
description: "Routes non-trivial coding, debugging, architecture, research, and UI tasks through Razvan's AI stack: Serena, Context7, Playwright, NotebookLM, subagents, and specialized skills. Use when the user asks for the full workflow, asks how to approach work, starts a non-trivial code task, or wants coordinated tool/skill usage."
---

# AI Workflow

Use this as the orchestrator skill. It does not replace specialized skills; it decides when to load them and which MCP/tool evidence is needed.

## Default flow

1. Clarify objective, constraints, and measurable acceptance criteria.
2. Classify the task using the routing table below.
3. Load the specialized skill(s) only when the task matches.
4. Use the cheapest useful orientation step first, then narrow to files/symbols/docs/tests.
5. Implement in a small vertical slice.
6. Verify with concrete evidence.
7. Finalize with project-specific checks and a concise summary.

## Routing table

| Task type | First move | Then use | Optional skill |
| --- | --- | --- | --- |
| Architecture, file relationships, where to look | Inspect project docs/README/CONTEXT and use focused file search | Read targeted files, use Serena where language support exists | improve-codebase-architecture |
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

- Use Serena for symbol-aware navigation when the project language server supports it.
- Use targeted `rg`/`find`/project docs for broad orientation when symbols are unknown.
- Use Context7 before changing fast-moving external APIs.
- Use Playwright when the success criterion is visual, interactive, browser, console, or network behavior.
- Use NotebookLM for video, paper, or large-source research synthesis.
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
