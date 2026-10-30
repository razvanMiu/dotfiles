---
name: redmine
description: Manage Redmine issues across configured instances via the key-safe local API wrapper. Use when querying, triaging, creating, updating, or commenting on Redmine tickets.
---

# Redmine

Use the local key-safe wrapper (`~/.config/scripts/redmine.ts`) instead of embedding API calls or raw secrets in prompts.

## Safety rules

- Never ask the user to paste API keys into chat.
- Never print API keys or write them into tracked files.
- Run Redmine commands through `with-secrets --` so API keys come from the local env/Bitwarden flow.
- Instance names such as `EEA` or `EDW` are local keys from `~/.config/redmine/config.json`. Do not guess names; run `with-secrets -- ~/.config/scripts/redmine.ts instances` or inspect local config before picking an instance.
- Ask which Redmine instance to use when ambiguous.

### Mandatory Write Permission Gate
- Read and query actions (`instances`, `trackers`, `statuses`, `priorities`, `projects`, `issues list`, `issue get`, `request GET ...`) may be executed autonomously when relevant to the user's task.
- **ALL write/mutation operations (`issue create`, `issue update`, `issue comment`, and mutating `request POST/PUT/DELETE`) require explicit user permission before execution.**
- For every write operation:
  1. Run the command **without** `--execute` (mutating commands are dry-run by default).
  2. Present the dry-run output to the user: state the target instance, issue/project ID, subject/notes, and full payload.
  3. Pause and explicitly ask for the user's approval to apply the change.
  4. Only append `--execute` after the user gives explicit consent in the current turn.

## Discovery and Querying

```sh
# List configured instances (* denotes default)
with-secrets -- ~/.config/scripts/redmine.ts instances

# Discover numeric metadata IDs for an instance
with-secrets -- ~/.config/scripts/redmine.ts trackers EEA
with-secrets -- ~/.config/scripts/redmine.ts statuses EEA
with-secrets -- ~/.config/scripts/redmine.ts priorities EEA
with-secrets -- ~/.config/scripts/redmine.ts projects EEA --limit 50

# List issues
with-secrets -- ~/.config/scripts/redmine.ts issues list EEA --directly-assigned-to-me --status open --limit 20
with-secrets -- ~/.config/scripts/redmine.ts issues list EEA --project my-project --query tracker_id=1 --query sort=updated_on:desc

# Inspect single issue with history and attachments
with-secrets -- ~/.config/scripts/redmine.ts issue get EEA 123 --include journals,attachments,relations,children
```

## Creating, Updating, and Commenting (Requires Permission)

```sh
# 1. Dry run (default) to show proposed creation payload:
with-secrets -- ~/.config/scripts/redmine.ts issue create EEA \
  --project my-project \
  --subject "Fix map legend rendering" \
  --tracker-id 1 \
  --priority-id 2 \
  --parent-id 100 \
  --custom-field 5="Sprint 42"

# 2. Only after explicit user approval:
with-secrets -- ~/.config/scripts/redmine.ts issue create EEA \
  --project my-project \
  --subject "Fix map legend rendering" \
  --tracker-id 1 \
  --priority-id 2 \
  --parent-id 100 \
  --custom-field 5="Sprint 42" \
  --execute

# Update issue status or assign:
with-secrets -- ~/.config/scripts/redmine.ts issue update EEA 123 --status-id 3 --notes "PR ready for review" --execute

# Add a comment:
with-secrets -- ~/.config/scripts/redmine.ts issue comment EEA 123 --notes "Implemented in commit abc123" --execute
```

## Arbitrary API Requests (`request`)

For endpoints or operations not covered by dedicated verbs (e.g., logging time entries, managing relations, project settings):

```sh
# Read arbitrary endpoint (autonomous)
with-secrets -- ~/.config/scripts/redmine.ts request EEA GET /users/current.json
with-secrets -- ~/.config/scripts/redmine.ts request EEA GET /time_entries.json

# Mutating endpoint (dry-run first, require approval before --execute)
with-secrets -- ~/.config/scripts/redmine.ts request EEA POST /time_entries.json \
  --data-json '{"time_entry":{"issue_id":123,"hours":1.5,"comments":"Investigated bug"}}' \
  --execute
```

## Workflow integration

- Use `triage` for incoming issue classification and acceptance criteria.
- Use `to-tickets` for breaking plans into actionable tickets.
- Use `implement` or `tdd` for implementation, then add a Redmine comment with validation evidence.
- Prefer `--directly-assigned-to-me` over `--assigned-to me` when the user asks for their own queue, because Redmine's `me` filter can include group-assigned issues.
- Prefer linking commits, test output, or concise reproduction notes in Redmine comments, but still get approval before posting the comment.
