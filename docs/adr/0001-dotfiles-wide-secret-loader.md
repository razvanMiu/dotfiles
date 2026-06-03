---
status: accepted
---

# Dotfiles-wide secret loader

The dotfiles repository uses `~/.config/.env` as the local untracked env contract and `~/.config/.env.example` as its tracked template, with `~/.config/scripts/with-secrets.ts` resolving `bws:<secret-key>` and `bws?:<secret-key>` references through the Bitwarden Secrets Manager CLI discovered on `PATH`. The `bws` executable is treated as a user binary under `~/.local/bin` rather than vendored in `~/.config`; secret references use Bitwarden secret keys for readability instead of immutable IDs, and BWS control variables are scrubbed from child commands by default to keep the Bitwarden access boundary separate from application secrets.

## Considered Options

- Vendor platform-specific `bws` binaries in the dotfiles repo; rejected because `~/.config` should track stateless configuration, not downloaded executables.
- Resolve secrets by Bitwarden UUID; rejected because key-based references make the dotfiles-wide env file readable and easier to maintain.
- Pass `BWS_*` variables through to child commands; rejected because the Bitwarden access token can fetch more than the individual resolved secrets.
