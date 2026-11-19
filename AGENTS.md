# Dotfiles Project Instructions

Scope: this repository only (`~/.config`). The global instructions in `~/.config/agents/AGENTS.md` still apply; this file adds what is specific to these dotfiles.

## Layout contract

- Keep stateless config in `~/.config`, state in `~/.local/state`, cache in `~/.cache`, and vendor/runtime packages in `~/.local/share`.
- Authored configuration lives here and is tracked; tool-generated state, credentials, packages, and caches are ignored or relocated to XDG paths.
- Document non-obvious configuration boundaries in `docs/adr/` with the decision, the rejected alternatives, and the consequences.

## Gates in this repository

- Run `scripts/check-dotfiles-safety --quiet` and `scripts/check-dotfiles-lint` before finalizing; `git diff --check` catches whitespace.
- The tracked pre-commit hook in `git/hooks/` runs all three gates, so a commit that fails any of them is rejected.
- Never commit credentials: `.gitignore` is a default-deny allowlist, and secrets belong in `with-secrets`.

## Agent skills

### Issue tracker

Issues and specs live as markdown under `.scratch/<feature-slug>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Five canonical roles, label string equal to role name. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `GLOSSARY.md` plus `docs/adr/` at the repo root. See `docs/agents/domain.md`.
