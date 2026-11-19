# XDG Dotfiles

This context describes the language used for Razvan's XDG-based dotfiles repository rooted at `~/.config`.

## Language

**Dotfiles-wide env file**:
A local, untracked environment file at `~/.config/.env` used by user-level dotfiles tools that need machine-specific or secret-bearing values. Its tracked template is `~/.config/.env.example`.
_Avoid_: global env, root env, project env

**User binary**:
An executable installed outside the dotfiles repo, typically under `~/.local/bin`, and discovered through `PATH`. User binaries are not tracked as configuration.
_Avoid_: vendored binary, config binary

**Cross-agent layer**:
The tracked, tool-neutral agent resource area at `~/.config/agents` that holds shared instructions and skills consumed by agent-specific adapters. Agent-specific directories such as Claude Code and Pi adapt from this layer rather than becoming the source of truth.
_Avoid_: Claude layer, Pi layer, shared Claude config

**Mixed tool directory**:
A tool directory under `~/.config` that contains both authored configuration and tool-generated data because the tool does not cleanly separate XDG config, state, and cache. Only authored configuration is tracked; generated data, auth, state, cache, package installs, and transcripts are ignored.
_Avoid_: stateful config, dirty config directory, XDG violation

**Shared skill directory**:
The cross-agent skill source at `~/.config/agents/skills`, also exposed as `~/.agents/skills`. Agent-specific skill directories may symlink to it when their skill format is compatible.
_Avoid_: Claude skills, Pi skills, copied skill cache

**Adapter instruction file**:
A tool-specific instruction file that loads or mirrors the cross-agent instructions while adding only the details needed by that tool. Adapter instruction files should use portable paths where possible so the dotfiles can be reused across machines.
_Avoid_: duplicated global instructions, tool fork, separate policy

**Portable dotfiles path**:
A tracked path expression that can work across Linux, macOS, and WSL without embedding a machine-specific home directory. Prefer relative paths, `~`, and XDG environment variables; reserve machine-specific absolute paths for untracked local files.
_Avoid_: hardcoded home path, machine path, Linux-only path

**Agent package directory**:
An XDG data directory for installed agent packages, vendor code, and runtime assets, separate from tracked agent configuration where the agent tool supports relocating those assets. If a tool hardcodes package assets under a mixed tool directory, those assets remain ignored until supported relocation exists.
_Avoid_: config package cache, tracked package install, Pi npm config

**Strict-supported XDG split**:
A migration policy that moves tool state, sessions, packages, and cache out of `~/.config` only where the tool exposes supported paths or environment variables. Tool-generated files that cannot be relocated safely remain ignored in their mixed tool directory.
_Avoid_: symlink hack, forced XDG split, config-state merge

**Autonomous within hard boundaries**:
A unified agent safety posture where coding agents may act autonomously after planning, while deterministic guards block non-negotiable risks such as committing secrets, auth files, state, cache, vendor assets, browser profiles, or running destructive commands outside intended roots without approval.
_Avoid_: dangerous mode, unrestricted agent, approval-first workflow

**Hard-boundary guard**:
A deterministic safety check that enforces non-negotiable dotfiles rules independently of agent instructions. Commit-time guards are mandatory final gates; tool-time guards are additive early warnings.
_Avoid_: safety prompt, best-effort instruction, vibe guard

**Current-file guard scope**:
The hard-boundary guard scope that checks staged files, tracked files, and untracked non-ignored files in the dotfiles repo. Ignored files are excluded from blocking checks by default and may be inspected by an explicit audit mode.
_Avoid_: full home scan, ignored-file blocker, staged-only check

**Bitwarden secret key reference**:
A required secret placeholder in the dotfiles-wide env file whose value has the form `bws:<secret-key>` and resolves by matching the Bitwarden Secrets Manager secret `key` within the configured project. If the referenced key is missing or ambiguous, the secret loader fails before running the child command.
_Avoid_: secret id reference, UUID reference

**Optional Bitwarden secret key reference**:
An optional secret placeholder in the dotfiles-wide env file whose value has the form `bws?:<secret-key>`. If the referenced key is missing, the secret loader unsets the corresponding child environment variable instead of failing.
_Avoid_: nullable secret, maybe secret

**BWS control environment**:
The Bitwarden Secrets Manager variables used by the secret loader itself, including `BWS_ACCESS_TOKEN`, `BWS_PROJECT_ID`, and `BWS_SERVER_URL`. These are scrubbed from child commands by default.
_Avoid_: Bitwarden app env, secret env
