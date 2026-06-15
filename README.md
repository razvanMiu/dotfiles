# Razvan's XDG Dotfiles

This repository is rooted at `~/.config` and tracks mostly stateless configuration. Runtime state, caches, secrets, installed packages, and generated files should stay out of git.

Remote:

```sh
git clone git@github.com:razvanMiu/dotfiles.git ~/.config
```

If `~/.config` already exists on a new machine, clone elsewhere first and copy/adopt files carefully.

## Core conventions

- Config: `~/.config`
- Data/vendor assets: `~/.local/share`
- State/logs/sessions: `~/.local/state`
- Cache: `~/.cache`
- Local binaries: `~/.local/bin`
- Shared agent resources: `~/.agents -> ~/.config/agents`
- Keyboard contract: `~/.config/KEYBINDINGS.md`

## Dotfiles-wide secrets

Tracked template:

```text
~/.config/.env.example
```

Local secret-bearing env file, never committed:

```text
~/.config/.env
```

`with-secrets` loads `~/.config/.env`, resolves Bitwarden Secrets Manager references, scrubs `BWS_*` control variables from the child environment by default, and runs a command:

```sh
with-secrets -- pi
with-secrets -- npm run dev
with-secrets --env-file ~/Work/project/.env -- npm run dev
with-secrets --pass-bws-env -- bws secret list "$BWS_PROJECT_ID"
```

Secret reference syntax:

```dotenv
REQUIRED_SECRET=bws:SECRET_KEY
OPTIONAL_SECRET=bws?:OPTIONAL_SECRET_KEY
```

Required references fail if the Bitwarden secret key is missing or duplicated in the configured project. Optional references are unset when missing.

Install the Bitwarden Secrets Manager CLI as a user binary on `PATH`, for example:

```text
~/.local/bin/bws
```

Expose the tracked script as a local command:

```sh
ln -sfn ~/.config/scripts/with-secrets.ts ~/.local/bin/with-secrets
```

See also:

```text
~/.config/docs/adr/0001-dotfiles-wide-secret-loader.md
```

## Shell: zsh

Tracked zsh config lives in:

```text
~/.config/zsh/
├── zshrc
├── env.zsh
├── paths.zsh
├── oh-my-zsh.zsh
├── aliases.zsh
└── tools.zsh
```

Bootstrap the home entrypoint:

```sh
cat > ~/.zshrc <<'EOF'
source "$HOME/.config/zsh/zshrc"
EOF
```

Set zsh as the login shell if needed:

```sh
chsh -s "$(command -v zsh)"
```

The zsh config sets XDG paths and Pi paths, including:

```sh
PI_CODING_AGENT_DIR="$XDG_CONFIG_HOME/pi/agent"
PI_CODING_AGENT_SESSION_DIR="$XDG_STATE_HOME/pi/sessions"
```

## Oh My Zsh

Oh My Zsh is not vendored in this repo. Install it under XDG data:

```sh
git clone https://github.com/ohmyzsh/ohmyzsh.git ~/.local/share/oh-my-zsh
```

Plugins expected by `~/.config/zsh/oh-my-zsh.zsh` live as vendor code under XDG data. The tracked `~/.config/zsh/oh-my-zsh/plugins/*` entries are symlinks to these directories:

```sh
mkdir -p ~/.local/share/zsh/plugins

git clone https://github.com/zsh-users/zsh-autosuggestions \
  ~/.local/share/zsh/plugins/zsh-autosuggestions

git clone https://github.com/zsh-users/zsh-syntax-highlighting.git \
  ~/.local/share/zsh/plugins/zsh-syntax-highlighting
```

Theme customization is tracked at:

```text
~/.config/zsh/oh-my-zsh/themes/eastwood.zsh-theme
```

## Git

Tracked git config:

```text
~/.config/git/config
```

Bootstrap `~/.gitconfig`:

```sh
cat > ~/.gitconfig <<'EOF'
[include]
	path = ~/.config/git/config
EOF
```

This repo uses a repo-local hooks path:

```sh
cd ~/.config
git config core.hooksPath git/hooks
```

The pre-commit hook runs the guard quietly:

```sh
scripts/check-dotfiles-safety --quiet
```

The guard checks staged files, tracked files, and untracked non-ignored files for forbidden dotfiles content such as secrets, auth files, sessions, caches, browser profiles, package/vendor installs, and generated state. Ignored files are not blocking by default; inspect them explicitly with:

```sh
scripts/check-dotfiles-safety --audit-ignored
```

## Agent configuration

Shared, tool-neutral agent resources live in:

```text
~/.config/agents
~/.agents -> ~/.config/agents
```

Canonical global agent instructions live in:

```text
~/.config/agents/AGENTS.md
```

Pi and Claude Code consume those canonical instructions differently:

```text
~/.config/pi/agent/AGENTS.md -> ../../agents/AGENTS.md
~/.config/claude/CLAUDE.md   # imports @../agents/AGENTS.md
```

The canonical marker for loading tests is:

```text
xdg-agent-canonical-v1
```

Verify pi loading:

```sh
cd ~/.config
PI_OFFLINE=1 pi --no-session -p "Reply only with the canonical marker from your loaded instructions."
```

Verify Claude Code loading:

```sh
cd ~/.config
CLAUDE_CONFIG_DIR="$XDG_CONFIG_HOME/claude" \
  claude -p "Reply only with the canonical marker from your loaded instructions."
```

Expected output for both:

```text
xdg-agent-canonical-v1
```

See also:

```text
~/.config/docs/adr/0002-agent-configuration-boundaries.md
```

## Claude Code

Claude Code uses the XDG config root through:

```sh
CLAUDE_CONFIG_DIR="$XDG_CONFIG_HOME/claude"
```

Tracked Claude config is limited to portable authored files such as:

```text
~/.config/claude/settings.json
~/.config/claude/CLAUDE.md
~/.config/claude/skills -> ../agents/skills
~/.config/claude/.gitignore
```

Claude Code is a mixed tool directory: it stores generated state next to authored config. Keep auth, history, transcripts, caches, plugin state, shell snapshots, backups, and downloads ignored.

The default Claude permission posture is autonomous within hard boundaries:

```json
{
  "permissions": {
    "defaultMode": "bypassPermissions"
  }
}
```

Do not rely on prompts alone for safety. The dotfiles safety guard is the final commit-time gate.

Claude also runs an early safety hook from `settings.json`:

```text
~/.config/scripts/claude-safety-hook
```

The hook blocks obvious destructive Bash commands, `git commit --no-verify`, forced `git add`, and writes to forbidden dotfiles paths such as auth, cache, state, browser profile, and package/vendor directories. It also runs `scripts/check-dotfiles-safety --quiet` after mutating tool calls in the dotfiles repo.

### Claude plugins and MCP policy

Treat plugins and MCP servers as supply-chain code. Prefer official, pinned, reviewed integrations. Avoid installing marketplace plugins just because they are available.

Current baseline:

```text
playwright MCP -> @playwright/mcp@0.0.76 --browser chromium
```

Candidate Claude plugins to review before installing:

```text
claude-code-setup
claude-md-management
code-review
pr-review-toolkit
security-guidance
context7
serena
```

Install only after checking plugin inventory, token cost, MCP servers, hooks, commands, and required credentials.

## Pi coding agent

Pi runtime config lives in:

```text
~/.config/pi/agent/settings.json
```

Global Pi instructions are the canonical cross-agent instructions exposed through a relative symlink:

```text
~/.config/pi/agent/AGENTS.md -> ../../agents/AGENTS.md
```

These instructions bias Pi toward disciplined engineering work: read before edit, small validated changes, subagents for noisy scouting/review, measurable verification, and no secret/runtime commits.

Do not commit Pi runtime/secrets. The local ignore file excludes:

```text
~/.config/pi/agent/auth.json
~/.config/pi/agent/sessions/
~/.config/pi/agent/bin/
~/.config/pi/agent/npm/
~/.config/pi/agent/git/
~/.config/pi/agent/mcp-cache.json
~/.config/pi/agent/mcp-npx-cache.json
~/.config/pi/agent/run-history.jsonl
~/.config/pi/agent/trust.json
*.log
```

Pi currently hardcodes user-installed extension packages under `PI_CODING_AGENT_DIR` as `npm/` and `git/`. Do not set `PI_PACKAGE_DIR` for this purpose; pi uses that variable for its own executable/package asset directory. Keep generated package directories ignored until pi exposes supported relocation paths.

### Install Pi packages

Current recommended package stack:

```sh
pi install npm:pi-web-access
pi install npm:pi-subagents
pi install npm:pi-observational-memory
pi install npm:pi-sensitive-guard
pi install git:github.com/elpapi42/pi-fork
pi install npm:pi-mcp-adapter
pi install npm:pi-rtk-optimizer
pi install npm:graphify-pi
```

Restart Pi after package changes.

### Pi MCP baseline

Pi uses `pi-mcp-adapter`. Keep global Pi MCP config in:

```text
~/.config/pi/agent/mcp.json
```

Current baseline matches Claude's MCP server set:

```text
playwright MCP -> @playwright/mcp@0.0.76 --browser chromium
context7 MCP  -> @upstash/context7-mcp@3.2.1
serena MCP    -> git+https://github.com/oraios/serena@v1.5.3
```

`context7` provides fresh library/API docs. `serena` provides semantic code navigation; Claude uses Serena's `claude-code` context while Pi uses Serena's generic `ide` context. Serena runs with `SERENA_HOME=${XDG_STATE_HOME:-$HOME/.local/state}/serena`; project-local `.serena/` data is generated, ignored, and blocked by the dotfiles safety guard if force-staged. Serena reminder hooks are not enabled for now because the existing Graphify and safety hooks already modify Claude tool behavior.

Validation note: Context7 starts cleanly on stdio. Serena starts and exposes MCP tools. For the local `~/.config` project, generated `.serena/project.yml` is configured for TypeScript/JavaScript, Python, Go, C/C++, Bash, and Lua; the matching language-server commands are expected on `PATH`.

Prefer standard shared MCP files such as `~/.config/mcp/mcp.json` only after deciding which servers should be cross-tool rather than Pi-specific.

### Graphify

Graphify maps project files into a queryable knowledge graph and wiki.

Global CLI install:

```sh
uv tool install graphifyy
```

Shared Claude skill location:

```text
~/.config/agents/skills/graphify
~/.config/claude/skills/graphify -> ../agents/skills/graphify
```

Pi integration is provided by `npm:graphify-pi`, which adds Pi graph-first guidance and a native `/graphify` command when `graphify-out/` exists.

Current Obsidian AI knowledge vault:

```text
~/.local/share/obsidian/vaults/ai-knowledge
```

For this dotfiles repo, treat `graphify-out/` and `.graphify/` as generated local state. Keep them ignored and export only reviewed notes to the Obsidian vault. Other projects may choose to commit `graphify-out/` when team graph sharing is desired.

The dotfiles corpus is constrained by:

```text
~/.config/.graphifyignore
```

That file excludes secrets, auth/session files, caches, browser profiles, package/vendor installs, Claude plugin checkouts, and generated Graphify output.

Dotfiles Graphify workflow:

```sh
cd ~/.config
# In Pi, prefer the Graphify skill / command so semantic extraction can use Pi-hosted agents.
/graphify . --wiki --no-viz

# Refresh Obsidian export after graph rebuilds or label changes.
graphify export obsidian --dir ~/.local/share/obsidian/vaults/ai-knowledge/dotfiles-graph
```

After first generation, rename generic communities in `graphify-out/.graphify_labels.json`, then regenerate report/wiki/Obsidian export. Keep `graphify-out/` ignored.

Useful query commands:

```sh
graphify query "Hard-boundary guard dotfiles safety guard"
graphify query "scripts/check-dotfiles-safety"
graphify query "scripts/claude-safety-hook"
graphify query "Agent Config Boundaries"
graphify query "Claude Code Adapter Canonical Cross-Agent Instructions"
graphify explain "scripts/check-dotfiles-safety"
graphify path "scripts/claude-safety-hook" "scripts/check-dotfiles-safety"
```

Prefer exact domain terms, file names, and community labels over broad questions. See `docs/graphify/query-guide.md` for query patterns and known anchors.

### Current Pi settings

`~/.config/pi/agent/settings.json` points Pi at shared skills through a portable relative path:

```json
{
  "enableSkillCommands": true,
  "skills": ["../../agents/skills"]
}
```

`~/.agents -> ~/.config/agents` is still created for tools that understand the shared agent convention.

## Shared agent skills

Shared skills are stored in:

```text
~/.config/agents/skills
```

Create the shared symlink:

```sh
ln -sfn ~/.config/agents ~/.agents
```

This lets Pi and other agent tools use the same location:

```text
~/.agents/skills
```

### External and local skills

Skills are copied, not symlinked, from configured Git sources. Sources are declared in:

```text
~/.config/agents/skill-sources.json
```

Current source includes:

```text
git@github.com:mattpocock/skills.git
```

Add more repositories with:

```json
{
  "name": "my-team",
  "repo": "git@github.com:example/team-skills.git",
  "ref": "main",
  "skills": [
    { "path": "skills/debugging" },
    { "path": "skills/code-review", "name": "team-code-review" }
  ]
}
```

Use an empty `ref` for the repo default branch. Use a skill `name` to avoid destination name collisions.

Your own skills can live directly under:

```text
~/.config/agents/skills/<skill-name>/SKILL.md
```

The updater preserves local skills and only removes skills listed in:

```text
~/.config/agents/.managed-skills.manifest
```

Update external skills with:

```sh
~/.config/scripts/update-agent-skills.sh
```

The script:

1. clones/updates each source into `~/.local/share/agent-skill-sources/<source-name>`
2. copies selected skills into `~/.config/agents/skills`
3. records managed skills in `~/.config/agents/.managed-skills.manifest`
4. ensures `~/.agents -> ~/.config/agents`

After running it:

```sh
cd ~/.config
git diff -- agents scripts pi/agent/settings.json
```

Review changes before committing. This pins the exact skill content in dotfiles.

## Suggested setup order on a new machine

```sh
# 1. clone dotfiles
git clone git@github.com:razvanMiu/dotfiles.git ~/.config

# 2. shell entrypoints
cat > ~/.zshrc <<'EOF'
source "$HOME/.config/zsh/zshrc"
EOF

cat > ~/.gitconfig <<'EOF'
[include]
	path = ~/.config/git/config
EOF

# 3. oh-my-zsh and plugins
git clone https://github.com/ohmyzsh/ohmyzsh.git ~/.local/share/oh-my-zsh
mkdir -p ~/.local/share/zsh/plugins
git clone https://github.com/zsh-users/zsh-autosuggestions ~/.local/share/zsh/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.local/share/zsh/plugins/zsh-syntax-highlighting

# 4. expose dotfile-managed scripts
ln -sfn ~/.config/scripts/with-secrets.ts ~/.local/bin/with-secrets

# 5. create local env from the tracked template
cp ~/.config/.env.example ~/.config/.env
$EDITOR ~/.config/.env

# 6. install bws as a user binary on PATH, e.g. ~/.local/bin/bws

# 7. shared agent symlink and skills
ln -sfn ~/.config/agents ~/.agents
~/.config/scripts/update-agent-skills.sh

# 8. install Pi packages
pi install npm:pi-web-access
pi install npm:pi-subagents
pi install npm:pi-observational-memory
pi install npm:pi-sensitive-guard
pi install git:github.com/elpapi42/pi-fork
pi install npm:pi-mcp-adapter
pi install npm:pi-rtk-optimizer
pi install npm:graphify-pi
uv tool install graphifyy

# 9. enable repo-local safety hook
cd ~/.config
git config core.hooksPath git/hooks
scripts/check-dotfiles-safety
```

## Commit workflow

```sh
cd ~/.config
git status
git diff
git add README.md CONTEXT.md docs/adr agents claude pi/agent scripts git zsh
git commit -m "update dotfiles"
git push
```

Never commit secrets, auth files, sessions, caches, browser profiles, package installs, or generated state.
