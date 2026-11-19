# Compaction Policy: keep decisions alive across context boundaries

Status: ready-for-agent

## Problem Statement

Long sessions silently lose the rules and decisions that govern them. Measured on
this machine's own session history: 115 compactions across 32 sessions, **none**
of which waited for a task boundary — every one fired with tool calls in flight
immediately before and after, 13 with failing-test output in that window. The
agent then re-read 1,215 files it had already read across compaction boundaries,
against a 23.6% baseline in sessions that never compacted, and the user had to
re-send the same instruction five times in a single session.

The information that dies is specific: standing constraints and *how to work*
instructions, not descriptions of what to build — the literature measures them
being dropped ~6x more often, with constraint violations rising from 0% to
30-59% once they leave context.

Two mechanisms cause it. First, memory eviction: 5,166 observations recorded
against 2,689 dropped (52%), the worst session retaining almost nothing across 29
compactions at a pool budget of 8K tokens — the exact budget at which the
literature finds every eviction decision irreversible. Second, a trigger that has
no notion of task state: the compaction moment is a pure function of the active
model's context window, so switching to a smaller-window model retroactively
lowers the threshold below the current context and compacts instantly (35 of 115
compactions followed a model switch; two confirmed cases compacted at 133,920 and
247,082 tokens against a 131,072 threshold).

## Solution

Compaction still happens automatically — every production harness ships it — but
it becomes boundary-aware, non-destructive, and measured.

Memory stops being a lossy buffer and becomes an append-only ledger: nothing is
destroyed to make room. Pressure is relieved by compressing the oldest material
into reflections and, once compression is exhausted, by archiving it behind a
discoverable index line so the agent knows a fact exists and can recall it. A
pinned tier holds constraints, decisions and rejected alternatives and is exempt
from every budget.

The trigger gains a phase gate: it fires only when nothing is in flight — no tool
running, no failing test, no question awaiting the user — with a hard ceiling that
forces a checkpoint rather than a silent compaction. Window-dependent limits scale
with the window so small models degrade gracefully instead of thrashing.

Finally, the policy becomes measurable: a cut-and-resume harness scores retention
at 25%, 50% and 75% of a trajectory, so "did compaction lose something" stops
being a feeling and becomes a number with a baseline.

## User Stories

1. As a developer working a long session, I want compaction to wait until I am
   between tasks, so that it never interrupts a plan I am in the middle of.
2. As a developer, I want my stated constraints preserved verbatim, so that the
   agent does not silently start violating a rule I set hours ago.
3. As a developer, I want decisions *and the alternatives I rejected* preserved,
   so that the agent does not re-propose something I already ruled out.
4. As a developer, I want the agent to know that archived material exists, so that
   it can ask for it instead of guessing or re-deriving.
5. As a developer, I want memory to be append-only, so that a budget never deletes
   the reason a design exists.
6. As a developer, I want compaction to be forbidden while a tool call is in
   flight, so that side effects are never torn down mid-operation.
7. As a developer, I want compaction to be forbidden while a test run is failing,
   so that I do not lose the failure context I was about to diagnose.
8. As a developer, I want a hard ceiling that forces a checkpoint, so that a
   session can never blow past the window silently.
9. As a developer switching to a cheaper or smaller model mid-session, I want the
   threshold to scale with that model's window instead of compacting instantly, so
   that a model switch is not a data-loss event.
10. As a developer running a session on a small-window model, I want the verbatim
    tail to shrink with the window, so that the session does not compact every
    turn.
11. As a developer, I want observation density high enough that memory reflects
    recent work, so that what I just decided is not missing from the record.
12. As a developer, I want the observation gap to stay far below the verbatim
    tail, so that unobserved work is always still recoverable from the transcript.
13. As a developer, I want the memory worker's scope ceiling documented, so that I
    know a 262K worker localizes memory even when the session model has 1M.
14. As a developer, I want the work product itself to live on disk, so that
    compaction cannot affect it.
15. As a developer, I want the durable record to cover constraints as well as the
    plan, because keeping only the plan measurably does not recover the loss.
16. As a developer, I want a reproducible acceptance test, so that a policy change
    can be judged against a baseline instead of an impression.
17. As a developer, I want the re-read rate reported, so that lost context shows up
    as duplicated work rather than as a vague suspicion.
18. As a developer, I want the number of times I am asked the same question
    counted, so that re-asking is treated as a defect.
19. As a coding agent, I want the compaction summary to state which constraints are
    in force, so that I can comply without re-reading the whole session.
20. As a coding agent, I want an index of archived observations, so that `recall`
    is discoverable rather than a lucky guess.
21. As a maintainer, I want the policy expressed in configuration where possible,
    so that the first slice is reversible without touching third-party code.
22. As a maintainer, I want the remaining behaviour changes isolated, so that they
    can be proposed upstream rather than maintained as a fork.

## Implementation Decisions

**Scope: two slices, deliberately sequential.**

*Slice A — configuration only, no code.* Expresses everything the existing knobs
can carry:

- Raise the observation pool budget far above the current 8K target, since the
  eviction literature finds all eviction decisions irreversible at 8K and no
  policy beats another once accuracy is matched.
- Restore observation and reflection density to the defaults (10K / 20K) from the
  current 3x-larger gaps; the memory worker is a free in-house model, so density
  costs latency, not money.
- Keep the trigger in ratio mode, which already scales down safely for small
  windows, and add per-model compaction overrides so the verbatim tail and reserve
  scale with each model's window. This removes the small-window thrash regime,
  where a fixed 60K tail plus a 16K reserve exceeds the window and the session
  compacts every turn.

Slice A deliberately does **not** implement the absolute-budget ceiling: ratio
mode is linear through the origin and cannot express `min(budget, window − reserve
− tail)`. That is recorded as a known limitation, not an oversight.

*Slice B — behaviour change in the memory extension, to be proposed upstream
rather than forked:*

- A pinned tier for constraints, decisions and rejected alternatives, exempt from
  pool pressure.
- Compression before archival: pool pressure collapses the oldest observations
  into reflections rather than tombstoning them.
- An index line naming what was archived and how to recall it, replacing silent
  tombstoning as the pressure valve.
- A phase gate on the compaction trigger: refuse to fire while a tool call is
  active, a test run is failing, or a question is awaiting the user; defer to the
  next clean boundary, and force a checkpoint at the hard ceiling.

**Domain vocabulary this spec fixes** (to be recorded in the glossary):

- *compaction* — the event that replaces the transcript prefix with a summary.
- *pruning* — the transcript prefix leaving context; the record stays on disk.
- *eviction* — a stored observation being excluded from future summaries.
- *checkpoint* — a compaction deliberately taken at a task boundary.
- *pinned tier* — records exempt from every budget.
- *scope ceiling* — the largest span a memory worker can relate at once, set by its
  own context window, independent of the session model.

Invariants:

- The observation gap stays well below the verbatim tail, or the gap becomes a
  real hole.
- The ledger is append-only; pressure is handled by compression, never destruction.
- The compaction summary is a cache, never the only copy of a decision.

## Testing Decisions

**What makes a good test here:** the observable behaviour is what a resumed agent
*does*, not what the policy code computes. Tests assert on retained content and
subsequent action, never on internal counters or message shapes.

**The single seam:** the `compaction` entry in a session transcript. Everything
this spec changes — trigger moment, retained content, archived-index presence — is
visible in that entry, and the transcript is already a first-class artifact of the
system. No new seam is introduced for slice A.

**Protocol (cut-and-resume):** take a session with an explicit plan and stated
constraints; cut at 25%, 50% and 75% of the trajectory; resume with a matched
action budget; score three things:

1. the resumed agent restates goal, constraints and open decisions from the
   durable artifact without being prompted;
2. re-reads of already-read files stay at or below the 23.6% no-compaction
   baseline, against a measured 1,215 cross-boundary re-reads as the "before";
3. zero re-asks of the same question.

**Prior art in this repo:** the session-integrity checker walks session JSONL and
validates links; this harness walks the same files and validates retention. The
forensic pass that produced this spec's numbers is the same shape and should be
promoted into the harness rather than rewritten.

**Slice B, if it happens:** the projection `entries → summary text` is pure and
model-free, so it becomes a unit seam — pressure applied to a synthetic ledger,
asserting pinned records survive, compression precedes archival, and the index line
names what left.

## Out of Scope

- Replacing the memory extension, or writing a competing compaction engine.
- Changing Pi's native compaction algorithm, which this policy treats as a fixed
  overflow guard.
- Any change to subagent context handling.
- Tuning the memory worker's model choice, or eliminating the scope ceiling caused
  by a smaller worker window.
- Cross-session or cross-agent memory.
- Publishing measured results for the harness beyond the single-session protocol.

## Further Notes

**Evidence base.** Harness practice: Claude Code, Codex CLI, Cursor, Aider,
opencode, Cline and Roo all summarize at a window threshold with a manual override
and keep a recent tail; none publishes a compaction evaluation. Measurements used
here: chained summarization drops ~1/3 of task rules on the first summary
(arXiv:2610.31778); append-only progressive summaries beat re-summarization by
8.66 points at 25% lower cost (CompactionBench, Microsoft); dropped constraints
convert to 30-59% violations (Governance Decay); disk-resident state survives
compaction while conversation-only state degrades sharply (CompAct-Bench);
irreversible eviction is 1.00 for every policy at an 8K budget (What Eviction
Destroys); keeping only the plan does not recover the loss (Plans Don't Persist);
a learned compaction trigger gains 9.2% on SWE-bench Verified (AutoCompact).

**Honest gaps.** No controlled study compares mid-task to boundary compaction, so
the phase gate rests on this machine's 115-compaction dataset plus anecdotal bug
reports of in-flight tool teardown. The absolute-budget ceiling is unimplementable
in ratio mode. The observation-density change increases ledger volume, which slice
A's larger pool accommodates only up to a point — if the pool refills and eviction
pressure returns, that is the signal that slice B is required rather than optional.
