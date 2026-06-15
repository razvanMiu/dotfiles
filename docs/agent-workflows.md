# AI Workflow Playbook

This playbook describes how to use the local AI CLI stack without turning every task into setup work.

Core rule: start with the cheapest source of orientation, then narrow to exact files, symbols, docs, tests, or UI evidence.

## Tool roles

- **Graphify**: project map and routing layer. Use first for architecture, file relationships, communities, and "where should I look?" questions when `graphify-out/graph.json` exists.
- **Serena**: symbol-aware code navigation and semantic edits. Use inside active code projects for symbol overview, definitions, references, declarations, implementations, rename, and symbol-body edits.
- **Context7**: fresh external library/API docs. Use before changing framework, SDK, or package usage.
- **Playwright**: browser/UI evidence. Use for screenshots, accessibility snapshots, console/network checks, and interaction validation.
- **NotebookLM**: research synthesis over videos, papers, docs, and large external sources.
- **Subagents/forks**: fresh-context scouting, review, diagnosis, or parallel investigation.
- **Safety guards**: deterministic boundaries. Do not bypass `scripts/check-dotfiles-safety`, Claude safety hooks, or pre-commit hooks.

## Default coding workflow

1. Clarify objective and acceptance criteria.
2. If the project has `graphify-out/graph.json`, query Graphify with exact domain terms, file names, or community labels.
3. Use Serena to inspect symbols in the likely files.
4. Read only the files needed to understand the change.
5. Make the smallest vertical slice.
6. Validate with tests, type checks, linters, Playwright, or concrete command output.
7. Run safety checks before commit.

Example:

```sh
graphify query "Resource filtering MapPageShell MarkerResourceMap"
# then use Serena: get_symbols_overview, find_symbol, find_referencing_symbols
npm test -- MapPageShell
```

## Codebase investigation

Use when the question is architectural or cross-file.

1. `graphify query "exact domain terms filenames community labels"`
2. If results are broad or fuzzy, refine with more exact terms.
3. Use `graphify path "A" "B"` when tracing a relationship.
4. Use Serena on the files Graphify identifies.
5. Read files only after routing narrows the target.

Do not begin with broad `rg`/`find` unless Graphify is absent, stale, or irrelevant to the inspected path.

## Symbol-level change

Use when changing a function, component, composable, class, or exported value.

1. Serena `get_symbols_overview` on the file.
2. Serena `find_symbol` for the target.
3. Serena `find_referencing_symbols` before changing behavior.
4. Edit the smallest symbol/body or nearby file region.
5. Run focused tests.

Prefer Serena over raw text search for references when the language server supports the project.

## Library or framework change

Use when touching Nuxt, Vue, Leaflet, Playwright, API clients, MCP SDKs, or other external packages.

1. Context7 `resolve-library-id` with a concrete task query.
2. Context7 `query-docs` for the exact API or pattern.
3. If docs are insufficient, use `code_search` or upstream source.
4. Implement and validate with tests/type checks.

Do not rely on model memory for fast-moving framework APIs.

## UI bug or UI feature

Use when behavior is visible in a browser.

1. Reproduce manually or with Playwright navigation.
2. Capture accessibility snapshot and/or screenshot.
3. Inspect console and network requests if relevant.
4. Fix in a small slice.
5. Re-run the Playwright interaction and focused tests.

Evidence should include what was seen before and after, not just "looks good".

## Hard bug or regression

Use the diagnosis loop.

1. Reproduce.
2. Minimize.
3. Form hypotheses.
4. Instrument only where needed.
5. Fix.
6. Add or update regression tests.

Use a fork/subagent for noisy scouting or independent failure diagnosis.

## Architecture or workflow change

Use when the decision affects conventions, safety, toolchain, or many files.

1. Interview first; make acceptance criteria explicit.
2. Use Graphify to inspect existing terminology and relationships.
3. Use `grill-with-docs` for domain-language alignment.
4. Update `CONTEXT.md` when a term becomes part of the domain language.
5. Add an ADR only for hard-to-reverse, surprising, or trade-off-heavy decisions.
6. Implement one vertical slice and validate.

## Research or video-heavy work

Use NotebookLM when the corpus is videos, papers, long docs, or multiple external sources.

1. Create or reuse a notebook.
2. Add sources.
3. Ask focused research questions.
4. Distill conclusions into tracked docs or issues only after review.

## Graphify refresh policy

Do not run a full Graphify rebuild after every edit.

- `graphify query ...`: cheap; use anytime.
- `graphify update`: use after meaningful code/doc changes when the graph should stay current.
- Full `/graphify . --wiki --no-viz`: use after architecture-scale changes, large refactors, or before a major review.
- Obsidian export/community relabel: use only when the visual/wiki output matters.

Small config changes and one-line docs usually do not need semantic refresh.

## Commit workflow

Before committing dotfiles work:

```sh
cd ~/.config
scripts/check-dotfiles-safety --quiet
git diff --check
```

Stage only intended files unless explicitly asked to commit everything. Let the pre-commit hook run.

## Stop condition

Stop adding plugins when the stack already covers the task. Prefer daily-use feedback over speculative tools.
