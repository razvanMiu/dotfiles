---
name: orchestrate
description: "Orchestrate work by delegating to subagents instead of doing everything in the parent: frame the goal, dispatch recon/research/implement/review/audit children, integrate their output. Use when the user explicitly asks to orchestrate, to work as an orchestrator, to use or fan out subagents, or to delegate parts of a task to other agents."
---

# Orchestrate

You are the orchestrator. You hold the goal, the acceptance criteria, and the final answer. Children do bounded work in their own context windows and hand back results, not transcripts.

Delegation has a coordination cost. Spend it deliberately.

## Decide first: delegate or do it yourself

Delegate when separation buys something concrete:

- **Parallelism** — independent areas can be investigated, implemented, or compared at once.
- **Async progress** — slow or broad work continues while you keep working.
- **Role specialism** — an independent read-only critique that you cannot give your own work.

Do it yourself when the work is small, sequential, blocked on user intent, or mostly synthesis. A task that fits in one command or one file is cheaper inline; dispatching it only moves the waiting.

State in one line which mode you chose and why, then proceed. An explicit request to orchestrate wins over the size gate: say the overhead once, then orchestrate as asked.

## Roles, by capability

| Role | Job | Hands back |
|---|---|---|
| **recon** | locate code, map call paths, list conventions | compressed summary (a file when the harness can hand one back) |
| **research** | external docs, APIs, changelogs, comparisons | brief (file or message) |
| **implement** | one bounded change with its own write scope | changed files + report |
| **review** | independent critique of a plan or diff | findings with evidence |
| **audit** | check whether a claim rests on its source | per-claim verdicts |

Map these onto whatever this harness actually offers. If it exposes named agents, match by capability and read their documented contracts; if it offers only a generic delegation primitive, write the role into the task text. For recon, research, implement, and audit, do the work inline when no matching role exists. Review is the exception: self-review is not independent review, so when no independent reviewer exists, say so plainly and let the user decide.

## Where this fits the workflow

This skill changes *who* does the work, not *when* stages happen. The global workflow keeps ownership of stage transitions, review policy, and commit authorization; here you only decide how work is delegated inside those stages.

1. **Frame** — the goal, the files or subsystems in play, and one done criterion you can check. Dispatch recon, plus research in parallel when external facts are needed. If the slice list runs past a handful, re-check scope with the user instead of slicing further.
2. **Plan** — you write the slice list: ordered slices, each with a done criterion and the write scope it may touch. Independent slices with disjoint write scopes may run in parallel; serialize slices where one builds on another.
3. **Per slice** — one child per slice, given the slice text, the pointers, and the criterion. You run that slice's verification command; a slice that fails its check never becomes the base of the next one. Keep it in the working tree, report one progress line, take the next.
4. **Review** — once every slice is in place, dispatch one independent review of the assembled packet.
5. **Commit** — through the workflow's authorization flow, and only with no critical or high finding still open.

A **packet** is the complete working-tree change for one user request plus the criteria it must satisfy — the unit one review covers and one commit sequence records.

## Delegation contract

- **One task, one write scope, one done criterion.** "Implement the slice" invites premature completion; "these three tests pass and the diff is clean" does not.
- **Pointers, not payloads.** Hand over paths, commit ranges, spec files. Pasted content defeats the isolation you are paying for.
- **Fresh context by default.** Give a child the parent conversation only when the harness supports that and the child genuinely needs it — it is slower and more expensive.
- **A failed child is information.** Narrow the scope or take the work back; change something before re-dispatching.
- **Fan out only as wide as the harness can actually run at once.** Extra dispatches queue, fail, or cost a turn; either way the plan is wrong.
- **Two failures on one slice ends the loop.** Stop, report what is verified, what is unverified, and what you would change — do not silently re-plan.
- **A plan-invalidating discovery goes back to the user.** When a slice reveals the plan is wrong, surface the finding and the revised plan rather than improvising.

## Review economy

- One implementation review per packet — not per slice. Reviewing a plan before coding is a separate decision the global workflow owns.
- After low or medium findings, fix and validate directly; do not open another pass.
- Repeat review only for critical or high remediation, security-sensitive changes, or an explicit request. **Critical/high** means correctness, security, or data loss; **low/medium** means style, naming, or structure.
- Fan several reviewers over one change only when the user asked for independent lenses, never merely to raise confidence.

## Parent accountability

Only you speak to the user; children never do. You own the final verification and the commit, and you label what you know: confirmed by a command or artifact, versus inferred. A child's report is evidence to check, not a conclusion to forward.

## Cost

Delegation multiplies tokens: every child re-reads its own context. Keep the parent lean, prefer the lightest agent that can honestly do the role, and stop delegating when the next dispatch would only re-derive what you already have.
