# EEA Website Backend (Plone)

## Environment
- Plone: 6.1.x (develop Makefile sets `PLONE_VERSION=6.1.3`)
- Python: README recommends 3.10/3.11; `develop/Makefile` defaults `PYTHON=python3.12`

## Repo layout
- Backend root: `backend/`
- Dev env: `backend/develop/`
- Sources list: `backend/develop/sources.ini`
- Backend requirements: `backend/develop/requirements.txt`, `backend/develop/constraints.txt`

## Common commands (run from `backend/develop/`)
- Bootstrap: `make` (runs bootstrap, install, develop)
- Start: `make start`
- RelStorage: `make relstorage`
- Develop add-ons: `make develop`
- Release info: `make release`

## Development notes
- Hot reload is not available; restart Plone after Python changes.
- Optional reload helper: use `plone.reload` / `dm.plonepatches.reload`, visit `http://localhost:8080/@@reload`.

## Testing
- CLI: `bin/zope-testrunner --test-path sources/<package>`
- VS Code / editor debugging: configure test paths.
