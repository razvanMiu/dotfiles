---
name: commit
description: "Git commit standard: group changes into revertible atomic commits, draft Conventional Commits messages, and present each for approval before committing or pushing. Use when committing work, writing commit messages, or splitting a bundled diff into commits."
---

# Git Commit Standard

Record finished, verified work as history that reads like a story: every commit is one logical change, independently revertible. Verification, review, and safety checks belong to the workflow (AGENTS.md)—this skill records work; it does not re-verify.

## Steps

1. **Group & split.** Read `git status` and `git diff` (plus `git diff --cached` if anything is already staged). Assign every changed file to exactly one commit, or leave it explicitly uncommitted and say so.
   - One logical change per commit: a refactor, a bug fix, a doc update, and a dependency bump are four commits.
   - Order: build/config changes first, refactors next, behavior changes after, tests alongside the feature they cover.
   - An explicit user instruction overrides the default: "commit everything in one commit" means one commit.
   - *Done when*: every changed file is assigned to a commit or explicitly excluded, and each planned commit is independently revertible.
2. **Stage.** Build each group with `git add <paths>`, or `git add -p` for files that mix concerns.
   - *Done when*: `git diff --cached` shows exactly one group.
3. **Draft the message** for each commit against the format below.
   - *Done when*: each subject is <=50 characters, imperative, period-free; each body states why the change was made.
4. **Present for approval.** For each commit show: the file list, `git diff --cached --stat`, and the full proposed message. Wait for explicit approval; fold in corrections and re-present.
   - *Done when*: the user approves each message.
5. **Commit & push.** Commit each approved group in order, passing the approved message verbatim to `git commit -m`. Push only when the user asked.
   - *Done when*: `git status` is clean apart from the intentionally uncommitted remainder.

## Message format

```
<type>(<scope>): <description>      <- subject: <=50 chars, imperative, no period
                                     <- blank line
<body>                              <- wrap at 72: why the change, notable consequences
                                     <- blank line
<footer>                            <- issue refs, BREAKING CHANGE
```

- **Subject**: imperative present (`Add`, `Fix`, `Drop`), first letter capitalized, no trailing period.
- **Scope** (optional): the module or directory the change touches.
- **Body** (optional): the reasoning and context a reviewer or `git blame` reader needs—the why, not a restatement of the diff.
- **Footer** (optional): `Closes #123` issue references; `BREAKING CHANGE: <what and why>` for breaking changes (equivalently `type(scope)!: ...` in the subject).

## Types

| type | change kind |
|---|---|
| `feat` | new capability |
| `fix` | bug fix |
| `refactor` | code change that neither fixes a bug nor adds a feature |
| `perf` | performance improvement |
| `docs` | documentation only |
| `style` | formatting, no behavior change |
| `test` | adding or correcting tests |
| `build` | build system or dependencies |
| `ci` | CI configuration |
| `chore` | maintenance and tooling |

## Examples

Bad: `update code` / Good: `fix(auth): reissue refresh token on expiry`

Bad: one commit bundling a refactor, a bug fix, and a README change / Good:

```
refactor(billing): extract invoice total calculation
fix(billing): round proration to minor units
docs: document proration rounding
```

Full message example:

```
feat(api): paginate session list endpoint

Unbounded lists made the /sessions response grow with workspace
size, stalling the dashboard on large tenants.

Closes #142
```
