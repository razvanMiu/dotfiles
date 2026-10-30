# EEA Website Frontend (Volto)

## Environment
- Plone: 6.1
- Volto: 18.x (see `frontend/package.json` for exact version)
- Node: ^20 || ^22 (preferred: 22)
- Package manager: Yarn 3.2.3

## Repo layout
- Root frontend: `frontend/`
- Source: `frontend/src/`
- Local add-ons: `frontend/src/addons/`
- Overrides/customizations: `frontend/src/customizations/`
- Theme entry: `frontend/src/theme.js`
- App config: `frontend/src/config.js`
- Routes: `frontend/src/routes.js`

## Common commands (run from `frontend/`)
- Start dev server: `yarn start`
- Build: `yarn build`
- Lint JS/TS: `yarn lint`
- Format check: `yarn prettier`
- Stylelint: `yarn stylelint`
- Tests: `yarn test`
- Cypress: `yarn cypress:open` or `yarn cypress:run`

## Reference Volto source
- `volto17/`: Volto 17 core code (branch `17.x.x`)
- `volto18/`: Volto 18 core code (branch `18.x.x`)
- `volto19/`: Volto 19 core code (branch `main`)
- These are read-only reference clones for AI agents to consult core Volto source.

## Add-ons
- Workspace add-ons are in `frontend/src/addons/*` and registered via `frontend/package.json` `addons`.
- Check `frontend/package.json` `resolutions` for pinned add-on versions.

## Add-on Make targets (example)
From `frontend/src/addons/volto-accordion-block` (`make help`):
- `make lint` / `make lint-fix`
- `make test` / `make test-update`
- `make cypress-open` / `make cypress-run`
- `make stylelint` / `make stylelint-fix`
- `make prettier` / `make prettier-fix`
- `make start`, `make install`, `make shell`, `make clean`
