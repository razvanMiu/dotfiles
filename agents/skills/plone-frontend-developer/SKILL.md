---
name: plone-frontend-developer
description: Develop and debug Plone Volto frontend code and add-ons in frontend/src/addons. Use when modifying Volto components, blocks, themes, routing, or debugging frontend tests and UI bugs.
---

# Plone Frontend Developer

## Overview

Diagnose and fix Volto/Plone frontend issues with a production-grade workflow. Use repo-specific guidance from `references/eea-website-frontend.md`.

## Core workflow

1. Clarify the symptom, the affected URL, and the expected behavior. Note the language/locale if the site is multilingual.
2. Locate the code path in `frontend/src/` or a local add-on in `frontend/src/addons/`.
3. Reproduce and isolate: use minimal changes, add short-lived logs if needed, and confirm the exact component/block responsible.
4. Implement the fix in the most appropriate layer (add-on). Favor tokenized values and shared helpers when updating styles.
5. Verify and establish done criteria: run the smallest relevant check (unit test, lint, or targeted runtime verification). The step is complete when the targeted suite passes with zero errors and git diff confirms minimal surgical changes.

## Common tasks

- Debugging a Volto add-on: find the add-on entry in `frontend/package.json` and trace the component in `frontend/src/addons/`.
- Token updates: replace hard-coded values with design tokens and update shared theme/config entries to keep consistency.
- Multilingual links: check locale-aware routing and URL generation before hardcoding language paths.

## Good practices

- For new features, add or update Cypress and/or Jest coverage.
- When writing Cypress tests, drive all state setup, user flows, and assertions strictly through the browser UI as a real user would interact with the application; avoid backend API shortcuts.
- When modifying add-ons, run `make lint` inside the add-on folder.
- Add-ons live in `frontend/src/addons/`; check each add-on's `make help` for available targets.
- When developing an add-on, start the Volto dev environment with `make; make start` in the add-on folder for faster feedback loops.
- Volto core code is available in `volto17/`, `volto18/`, and `volto19/` (read-only reference clones) for consulting core Volto source.
- To run cypress tests, use `make cypress-run` in the add-on folder.
- After defining i18n messages, run `make i18n`.
- Use Playwright MCP for quick UI verification when needed; Volto hot reload covers most changes.
- If adding new shadow customizations, ask to restart the frontend server.

## Repo context

Use `references/eea-website-frontend.md` for environment details, repo layout, and commands.

## Notes

- Prefer readable, minimal diffs and avoid refactors unless they are required for the fix.
- Keep temporary logs out of the final patch.
- When unsure about conventions, inspect nearby code or similar components in this repo.
