# Roadmap: BMAD Automate

## Milestones

- ✅ **v1.0 Status-Based Workflow Routing** — Phases 1-5 (shipped 2026-01-08)
- ✅ **v1.1 Full Story Lifecycle** — Phases 6-13 (complete)

## Completed Milestones

- ✅ [v1.0 Status-Based Workflow Routing](milestones/v1.0-ROADMAP.md) (Phases 1-5) — SHIPPED 2026-01-08

<details>
<summary>✅ v1.0 Status-Based Workflow Routing (Phases 1-5) — SHIPPED 2026-01-08</summary>

**Delivered:** Automatic workflow routing based on sprint-status.yaml, eliminating manual workflow selection.

- [x] Phase 1: Sprint Status Reader (1/1 plans) — completed 2026-01-08
- [x] Phase 2: Workflow Router (1/1 plans) — completed 2026-01-08
- [x] Phase 3: Update Run Command (1/1 plans) — completed 2026-01-08
- [x] Phase 4: Update Queue Command (1/1 plans) — completed 2026-01-08
- [x] Phase 5: Epic Command (1/1 plans) — completed 2026-01-08

</details>

### ✅ v1.1 Full Story Lifecycle (Complete)

**Milestone Goal:** Run the complete story lifecycle (create→dev→review→commit) for each story before moving to the next, with error recovery, dry-run mode, and enhanced progress UI.

#### Phase 6: Lifecycle Definition ✅

**Goal**: Define the full workflow sequence per status and status transitions
**Depends on**: v1.0 complete
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 06-01: Story Lifecycle Sequence (TDD) — completed 2026-01-08

**Details:**

- Add `git-commit` to the workflow chain after `code-review`
- Define status transitions: create→ready-for-dev, dev→review, review→done
- Update router package with lifecycle sequence logic

#### Phase 7: Story Lifecycle Executor ✅

**Goal**: New package that runs the complete workflow sequence for one story
**Depends on**: Phase 6
**Research**: Unlikely (internal patterns)
**Plans**: 2/2 complete

Plans:

- [x] 07-01: Status Writer (TDD) - Add UpdateStatus to status package — completed 2026-01-09
- [x] 07-02: Lifecycle Executor (TDD) - Orchestrate full story lifecycle — completed 2026-01-09

**Details:**

- New `lifecycle` package in `internal/lifecycle`
- Executor runs: status→workflow→update status→next workflow→etc.
- Ends with git-commit+push→mark done
- Updates sprint-status.yaml after each step

#### Phase 8: Update Run Command ✅

**Goal**: `run <story>` executes full lifecycle, not just one workflow
**Depends on**: Phase 7
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 08-01: Run command with lifecycle execution (TDD) — completed 2026-01-09

**Details:**

- `run` command uses lifecycle executor
- Completes story entirely: create→dev→review→commit→done

#### Phase 9: Update Epic Command ✅

**Goal**: Epic uses lifecycle executor, full cycle per story before moving to next
**Depends on**: Phase 8
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 09-01: Epic command with lifecycle execution (TDD) — completed 2026-01-09

**Details:**

- Epic command uses lifecycle executor
- Each story runs to completion before next story starts
- Maintains existing fail-fast behavior

#### Phase 10: Update Queue Command ✅

**Goal**: Queue also uses lifecycle executor for consistency
**Depends on**: Phase 9
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 10-01: Queue command with lifecycle execution (TDD) — completed 2026-01-09

**Details:**

- Queue command uses lifecycle executor
- Consistent behavior with epic and run commands

#### Phase 11: Error Recovery & Resume ✅

**Goal**: Save progress state when workflow fails, resume from failure point
**Depends on**: Phase 10
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 11-01: State Persistence (TDD) — completed 2026-01-09

**Details:**

- Save lifecycle state to file when workflow fails
- `--resume` flag to continue from failure point
- Track which step failed and resume from there

#### Phase 12: Dry Run Mode ✅

**Goal**: Preview what would happen without executing
**Depends on**: Phase 11
**Research**: Unlikely (internal patterns)
**Plans**: 2/2 complete

Plans:

- [x] 12-01: GetSteps Method (TDD) - Add GetSteps to lifecycle executor — completed 2026-01-09
- [x] 12-02: Dry Run Flags - Add --dry-run to run, queue, epic commands — completed 2026-01-09

**Details:**

- `--dry-run` flag for run, queue, epic commands
- Shows workflow sequence without executing
- Lists stories and their lifecycle steps

#### Phase 13: Enhanced Progress UI ✅

**Goal**: Better visibility into lifecycle progress
**Depends on**: Phase 12
**Research**: Unlikely (internal patterns)
**Plans**: 1/1 complete

Plans:

- [x] 13-01: Step Progress UI — completed 2026-01-09

**Details:**

- Show current step and remaining steps
- Overall epic/queue progress indicator
- Estimated time based on previous story durations

## Progress

| Phase                       | Milestone | Plans Complete | Status   | Completed  |
| --------------------------- | --------- | -------------- | -------- | ---------- |
| 1. Sprint Status Reader     | v1.0      | 1/1            | Complete | 2026-01-08 |
| 2. Workflow Router          | v1.0      | 1/1            | Complete | 2026-01-08 |
| 3. Update Run Command       | v1.0      | 1/1            | Complete | 2026-01-08 |
| 4. Update Queue Command     | v1.0      | 1/1            | Complete | 2026-01-08 |
| 5. Epic Command             | v1.0      | 1/1            | Complete | 2026-01-08 |
| 6. Lifecycle Definition     | v1.1      | 1/1            | Complete | 2026-01-08 |
| 7. Story Lifecycle Executor | v1.1      | 2/2            | Complete | 2026-01-09 |
| 8. Update Run Command       | v1.1      | 1/1            | Complete | 2026-01-09 |
| 9. Update Epic Command      | v1.1      | 1/1            | Complete | 2026-01-09 |
| 10. Update Queue Command    | v1.1      | 1/1            | Complete | 2026-01-09 |
| 11. Error Recovery & Resume | v1.1      | 1/1            | Complete | 2026-01-09 |
| 12. Dry Run Mode            | v1.1      | 2/2            | Complete | 2026-01-09 |
| 13. Enhanced Progress UI    | v1.1      | 1/1            | Complete | 2026-01-09 |
