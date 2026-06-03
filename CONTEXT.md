# XDG Dotfiles

This context describes the language used for Razvan's XDG-based dotfiles repository rooted at `~/.config`.

## Language

**Dotfiles-wide env file**:
A local, untracked environment file at `~/.config/.env` used by user-level dotfiles tools that need machine-specific or secret-bearing values. Its tracked template is `~/.config/.env.example`.
_Avoid_: global env, root env, project env

**User binary**:
An executable installed outside the dotfiles repo, typically under `~/.local/bin`, and discovered through `PATH`. User binaries are not tracked as configuration.
_Avoid_: vendored binary, config binary

**Bitwarden secret key reference**:
A required secret placeholder in the dotfiles-wide env file whose value has the form `bws:<secret-key>` and resolves by matching the Bitwarden Secrets Manager secret `key` within the configured project. If the referenced key is missing or ambiguous, the secret loader fails before running the child command.
_Avoid_: secret id reference, UUID reference

**Optional Bitwarden secret key reference**:
An optional secret placeholder in the dotfiles-wide env file whose value has the form `bws?:<secret-key>`. If the referenced key is missing, the secret loader unsets the corresponding child environment variable instead of failing.
_Avoid_: nullable secret, maybe secret

**BWS control environment**:
The Bitwarden Secrets Manager variables used by the secret loader itself, including `BWS_ACCESS_TOKEN`, `BWS_PROJECT_ID`, and `BWS_SERVER_URL`. These are scrubbed from child commands by default.
_Avoid_: Bitwarden app env, secret env
