---
name: rename-session
description: Print copy-pasteable session name candidates from session work using YYYY-MM-DD_AgentName_Verb-Task_TargetID.
disable-model-invocation: true
---

# Rename Session Name Suggestions

Derive candidate names from the session work. Print only; never rename or edit files.

## Gather evidence

Use the least evidence needed, in this order:
1. session start date and `cwd`
2. first meaningful user request
3. dominant role implied by the request/cwd
4. compaction summaries, if any
5. referenced issue IDs, repo/service/domain names
6. later requests only if the task changed

Stop once the dominant task and target are clear. Do not scan the whole transcript.

## Name form

`YYYY-MM-DD_AgentName_Verb-Task_TargetID`

- **Date**: session start date
- **AgentName**: dominant role (`pi`, `researcher`, `reviewer`, `developer`, etc.)
- **Verb-Task**: 3-7 lowercase hyphenated words led by the main action
- **TargetID**: optional lowercase hyphenated target slug

Normalize components: lowercase ASCII, hyphens between words, no spaces, no trailing punctuation.

## Output

Print exactly three candidates: primary recommendation, shorter alternative, and more specific alternative. Label each outside its code block; put only the name inside. Follow with one short sentence identifying the evidence and any material assumption.

## Examples

- Current session: `2026-10-21_system-manager_help-with-maintaining-mac-system-integrity`
- Ranking pi models: `2026-10-21_researcher_rank-available-models`
- Redmine tracker work: `2026-10-15_pi_investigate-issue-export_redmine-task-tracker`

Omit TargetID when no concrete target exists.

## Completion criteria

- Exactly three copy-pasteable candidates shown
- Primary recommendation traces to session evidence
- No mutation of any session