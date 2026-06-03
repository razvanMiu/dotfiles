# Shared Agent Resources

This directory contains shared agent resources used by Pi and any other agent harness that understands the `~/.agents` convention.

## Layout

```text
~/.config/agents/
├── skills/                         # copied Agent Skills plus local custom skills
├── skill-sources.json              # declarative Git skill sources
└── .managed-skills.manifest        # skills managed by update-agent-skills.sh
```

`~/.agents` should be a symlink to this directory:

```sh
ln -sfn ~/.config/agents ~/.agents
```

Pi loads skills through `~/.agents/skills`, not through a Pi-specific path. This keeps the skills reusable across agent tools.

## Local custom skills

Put your own skills directly under:

```text
~/.config/agents/skills/<skill-name>/SKILL.md
```

The updater only removes skills listed in `.managed-skills.manifest`, so local custom skills are preserved.

## External skill repositories

External skill repos are configured in:

```text
~/.config/agents/skill-sources.json
```

Format:

```json
{
  "sources": [
    {
      "name": "source-name",
      "repo": "git@github.com:example/team-skills.git",
      "ref": "main",
      "skills": [
        { "path": "skills/debugging" },
        { "path": "skills/code-review", "name": "team-code-review" }
      ]
    }
  ]
}
```

`ref` may be empty for the repository default branch. Use `name` on a skill to avoid destination name collisions.

## Updating copied skills

Run:

```sh
~/.config/scripts/update-agent-skills.sh
```

The script clones/updates each source under:

```text
~/.local/share/agent-skill-sources/<source-name>
```

Then it copies the selected skills into:

```text
~/.config/agents/skills
```

The copied skills are committed to dotfiles, which pins the exact skill versions used by Pi. Review `git diff` after updating before committing.
