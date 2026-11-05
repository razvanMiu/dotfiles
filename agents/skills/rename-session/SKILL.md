---
name: rename-session
description: "Print copy-pasteable session name candidates from session work using YYYY-MM-DD_<project>_<verb-task>[_<#ticket>]."
disable-model-invocation: true
---

# Rename Session Name Suggestions

Derive candidate names from the session work. Print only; never rename or edit files.

## Gather evidence

Use the least evidence needed, in this order:
1. session start date and `cwd` (to determine `<project>`)
2. first meaningful user request
3. referenced issue or ticket numbers (e.g. `#310194`)
4. compaction summaries, if any
5. later requests only if the task changed

Stop once the dominant task and target are clear. Do not scan the whole transcript.

## Name form

`YYYY-MM-DD_<project>_<verb-task>[_<#ticket>]`

- **Date**: session start date (`YYYY-MM-DD`)
- **Project**: repository or addon slug (e.g. `dotfiles`, `eea-backend`, `eea-frontend`, `volto-eea-kitkat`)
- **Verb-Task**: 2-4 lowercase hyphenated words led by the main action (e.g. `fix-catalog-indexing`, `add-approval-mode`)
- **Ticket**: optional issue ID prefixed with `#` (e.g. `_#310194`) when an issue or PR is referenced

Normalize components: lowercase ASCII, hyphens between words in task slug, no spaces, no trailing punctuation.

## Output

Print exactly three candidates: primary recommendation, shorter alternative, and more specific alternative. Label each outside its code block; put only the name inside. Follow with one short sentence identifying the evidence and any material assumption.

## Examples

- Dotfiles work: `2026-11-04_dotfiles_add-approval-mode-and-tmux-sync`
- Redmine issue fix: `2026-11-02_eea-backend_fix-catalog-indexing_#310194`
- Frontend feature: `2026-10-25_eea-frontend_merge-systems-services_#39124`
- Addon development: `2026-10-18_volto-eea-kitkat_upgrade-cypress-tests`

Omit ticket suffix when no ticket or issue is referenced.

## Completion criteria

- Exactly three copy-pasteable candidates shown
- Primary recommendation traces to session evidence
- Follows the canonical `YYYY-MM-DD_<project>_<verb-task>[_<#ticket>]` format
- No mutation of any session
