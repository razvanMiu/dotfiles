---
name: redmine
description: Work with Redmine issue trackers through the local key-safe Redmine API wrapper. Use when the user asks to inspect, create, triage, update, comment on, or synchronize Redmine issues across configured Redmine instances.
---

# Redmine

Use the local wrapper instead of embedding API calls or secrets in prompts.

## Safety rules

- Never ask the user to paste API keys into chat.
- Never print API keys or write them into tracked files.
- Run Redmine commands through `with-secrets --` so API keys come from the local env/Bitwarden flow.
- Instance names such as `EEA` or `EDW` are local keys from `~/.config/redmine/config.json`; do not assume `primary` or `secondary` are real. Run `with-secrets -- ~/.config/scripts/redmine.ts instances` or read the local config before choosing an instance.
- Ask which Redmine instance to use when ambiguous.
- Read/list/get actions may be run directly when they match the user's request.
- Create, update, delete, and comment actions must be guarded by explicit user approval.
- For every mutation: first run or show the dry-run command without `--execute`, summarize the exact target instance, issue/project, and payload, then ask for approval.
- Do not add `--execute` until the user explicitly approves that exact mutation in the current conversation.
- Mutating commands are dry-run by default; `--execute` is the final step after approval, not a permission substitute.

## Commands

Config lives outside tracked secrets:

```text
~/.config/redmine/config.json          # local, ignored
~/.config/redmine/config.example.json  # tracked template
```

API keys are referenced by env var name in the config and should be populated by `~/.config/.env` via `bws:<secret-key>` references.

Typical commands:

```sh
with-secrets -- ~/.config/scripts/redmine.ts instances
with-secrets -- ~/.config/scripts/redmine.ts issues list EEA --directly-assigned-to-me --status open --limit 20
with-secrets -- ~/.config/scripts/redmine.ts issue get EEA 123 --include journals,attachments
with-secrets -- ~/.config/scripts/redmine.ts issue create EEA --project my-project --subject "Fix map legend" --description "..."
# Only after explicit user approval for the exact dry-run payload:
with-secrets -- ~/.config/scripts/redmine.ts issue create EEA --project my-project --subject "Fix map legend" --description "..." --execute
with-secrets -- ~/.config/scripts/redmine.ts issue comment EEA 123 --notes "Implemented in commit abc123" --execute
```

## Workflow integration

- Use `triage` for incoming issue classification and acceptance criteria.
- Use `to-issues` for breaking plans into Redmine tickets.
- Use `ai-workflow` for implementation, then add a Redmine comment with validation evidence.
- Prefer `--directly-assigned-to-me` over `--assigned-to me` when the user asks for their own queue, because Redmine's `me` filter can include group-assigned issues.
- Prefer linking commits, test output, screenshots, or concise reproduction notes in Redmine comments, but still get approval before posting the comment.
