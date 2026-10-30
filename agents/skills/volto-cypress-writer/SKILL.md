---
name: volto-cypress-writer
description: Cypress E2E writing and maintenance for Plone Volto add-ons under frontend/src/addons, including UI-first setup, selectors, helpers, flake reduction, and add-on Cypress commands.
---

# Volto Cypress Writer

Use this skill when writing, updating, reviewing, or debugging Cypress tests for Volto add-ons in `frontend/src/addons/*`.

## Workflow

1. Inspect the target add-on before editing:
   - `cypress.config.js`
   - `cypress/support/e2e.js`
   - `cypress/support/commands.js`
   - existing `cypress/support/*.js` helpers
   - nearby `cypress/e2e/*.cy.js` specs
2. Follow the add-on's local file layout, naming, and helper style.
3. Put specs in `frontend/src/addons/<addon>/cypress/e2e/*.cy.js`.
4. Put reusable feature helpers in `cypress/support/<feature>.js`.
5. Keep global Cypress commands minimal; prefer imported helpers for feature-specific flows.
6. Improve fragile selectors and waits when touching nearby test code.

## UI-First Rule

New tests must set up state, navigate, edit, save, and assert through the browser UI.

- Drive all test setup exclusively through the UI; restrict `cy.request()` and REST helpers (`autologin`, `createContent`) to legacy suites where explicitly requested.
- Read existing REST helpers for context when helpful, but author all new flows using user-visible UI interactions.
- Use the UI for the behavior under test even when an older suite uses API shortcuts.
- If UI-only setup is impractical for a narrow reason, stop and explain the constraint before changing the approach.

## Selectors

- Target user-facing roles, labels, placeholders, visible text, or stable `data-cy` / `data-testid` attributes.
- If stable selectors are missing and component edits are in scope, add explicit test selectors to the component instead of relying on DOM shape.
- Avoid positional selectors (`:nth-child()`) and fragile style-driven classes (`.ui.basic.icon.button...`).
- When Volto or Semantic UI markup forces a complex selector, encapsulate it in a named helper such as `openBlockChooser()` or `saveCurrentPage()`.
- Keep assertions user-visible where possible: URL, heading, block text, saved view output, modal state, toolbar state, validation message.

## Reliability

- Each spec must pass independently; isolate test cases using `beforeEach` for browser state and content setup.
- Rely on Cypress retry-ability: `cy.get(...).should(...)`, `cy.contains(...).should(...)`, and URL assertions.
- Wait on observable UI or network state (`cy.intercept().as()`, `cy.wait('@alias')`) rather than fixed sleeps (`cy.wait(ms)`). If a Volto editor limitation requires a sleep, keep it local and document the reason.
- Rely on standard Cypress timeouts; reserve per-command `{ timeout: ... }` overrides for documented heavy operations.
- Split chains after actions before assertions:

```js
cy.get('[data-cy="toolbar-save"]').click();
cy.get('[data-cy="view-page"]').should('be.visible');
```

- Prefer natural user interactions; reserve `{ force: true }` strictly for cases where Volto overlays make it unavoidable, accompanied by an explanatory comment.
- Let application errors surface cleanly; reserve exception suppression strictly for verified third-party vendor scripts where the user accepts the tradeoff.

## Volto Add-on Patterns

- Add-on Cypress config usually sets `baseUrl: 'http://localhost:3000'` and passes Plone API settings through `CYPRESS_API_PATH`, `RAZZLE_DEV_PROXY_API_PATH`, or `RAZZLE_INTERNAL_API_PATH`.
- Add-on tests normally live beside the add-on, not in root `frontend/cypress/e2e`.
- Cypress tests for add-ons must pass against both Volto 18 and Volto 17 unless the user explicitly scopes the task to one version.
- For block tests, cover the edit flow and the saved view:
  - create or navigate to a page through the UI
  - add the block through the block chooser
  - configure the fields through the sidebar or inline editor
  - save
  - assert the view mode output
- For editor interactions, prefer existing Slate helpers if the add-on already has them. Add focused helpers only when they reduce repeated low-level selection/type logic.

## Commands

Run add-on Cypress commands from `frontend/src/addons/<addon>` when the add-on has its own Makefile.

Start by checking the add-on's available targets:

```sh
make help
```

Default Volto 18 flow:

```sh
make; make start
make cypress-run
make cypress-open
```

Volto 17 compatibility flow:

```sh
VOLTO_VERSION=17 make; make start
make cypress-run
```

`VOLTO_VERSION` selects the add-on environment during `make`; it does not change Cypress itself. Run Cypress against the Volto version that is currently started.

Run one spec from the add-on directory:

```sh
../../../node_modules/cypress/bin/cypress run --spec "cypress/e2e/<spec>.cy.js"
```

If the local binary is not convenient, `npx` can run the installed Cypress package:

```sh
npx cypress run --spec "cypress/e2e/<spec>.cy.js"
```

Useful debugging and CI-style targets:

```sh
make shell
make check-ci
make cypress-ci
```

`make shell` opens a shell in the add-on frontend container. `make check-ci` and `make cypress-ci` are common hidden targets in these add-ons; confirm they exist in the target add-on's Makefile before using them.

Useful non-Cypress checks:

```sh
make lint
make prettier
make stylelint
make test
```

Use the smallest relevant command for verification. For a Cypress-only add-on change, prefer the target add-on's `make cypress-run` in both Volto 18 and Volto 17 environments.

## Review Checklist

- The test exercises the feature like a user, not by mutating backend state.
- Selectors are stable or wrapped in named helpers.
- The test asserts both edit-time behavior and saved output when relevant.
- The suite can run test cases independently.
- The Cypress behavior works in both Volto 18 and Volto 17, or the version-specific limitation is clearly reported.
- Network or async transitions wait on observable conditions, not fixed sleeps.
- Per-command `{ timeout: ... }` overrides are avoided or explicitly justified.
- Helpers are local to the add-on unless multiple specs in the add-on need them.
- The command used to verify is reported back to the user.
