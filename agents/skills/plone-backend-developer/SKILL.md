---
name: plone-backend-developer
description: Develop and debug Plone 6.1 backend code and add-ons in backend/develop. Use when modifying Python packages, sources.ini, server configuration, or debugging Zope tests.
---

# Plone Backend Developer

## Overview

Implement and debug Plone 6 backend changes with a reliable workflow and repo-specific conventions. Consult `references/eea-website-backend.md` for environment details and commands.

## Core workflow

1. Confirm the failing endpoint, request path, and expected behavior.
2. Locate the relevant package in `backend/develop/sources.ini` or installed add-ons.
3. Reproduce locally; add minimal logging only when needed.
4. Apply the fix in the add-on or config, then restart Plone (or use reload helpers when appropriate).
5. Validate and establish done criteria: run a targeted `zope-testrunner` invocation (e.g. `bin/zope-testrunner --test-path sources/<package>`). The step is complete when all targeted tests pass, Plone starts cleanly without traceback, and git diff confirms only intended changes.

## Common tasks

- Add or update a backend add-on in `backend/develop/` and run `make develop` to refresh sources.
- Adjust Plone configuration and verify the app starts with `make start` or `make relstorage`.
- Investigate errors in `bin/runwsgi` output and trace back to the relevant add-on module.

## Good practices

- Restart Plone after Python changes; hot reload is not available.
- Use the reload endpoint (`@@reload`) only for quick iteration, then restart for confirmation.
- Keep changes small and add targeted tests for regressions.

## Repo context

Use `references/eea-website-backend.md` for environment details, commands, and dev notes.
