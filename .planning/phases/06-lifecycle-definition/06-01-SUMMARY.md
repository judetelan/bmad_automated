---
phase: 06-lifecycle-definition
plan: 01
subsystem: router
tags: [lifecycle, status, workflow-sequence, tdd]

# Dependency graph
requires:
  - phase: 02-workflow-router
    provides: GetWorkflow function and sentinel errors
provides:
  - LifecycleStep struct for workflow + status transition
  - GetLifecycle function returning full lifecycle sequence
affects:
  [07-story-lifecycle-executor, 08-update-run-command, 09-update-epic-command]

# Tech tracking
tech-stack:
  added: []
  patterns: [switch-statement-routing, table-driven-tests]

key-files:
  created: [internal/router/lifecycle.go, internal/router/lifecycle_test.go]
  modified: []

key-decisions:
  - "git-commit has NextStatus=done (status already done after code-review)"
  - "LifecycleStep uses status.Status type for type safety"

patterns-established:
  - "Lifecycle sequence pattern: list of {Workflow, NextStatus} tuples"

issues-created: []

# Metrics
duration: 3min
completed: 2026-01-08
---

# Phase 6 Plan 01: Story Lifecycle Sequence Summary

**GetLifecycle function returning full workflow sequence from any status to done, with git-commit as final step**

## Performance

- **Duration:** 3 min
- **Started:** 2026-01-08T21:00:00Z
- **Completed:** 2026-01-08T21:03:00Z
- **Tasks:** RED → GREEN (no refactor needed)
- **Files modified:** 2

## Accomplishments

- LifecycleStep struct defined with Workflow and NextStatus fields
- GetLifecycle function implemented returning sequence from current status to done
- Full test coverage for all status values including error cases
- git-commit workflow added as final step after code-review

## Task Commits

TDD cycle commits:

1. **RED: Failing test** - `2a9abb0` (test)
2. **GREEN: Implementation** - `90bacdf` (feat)

**Plan metadata:** (this commit)

## Files Created/Modified

- `internal/router/lifecycle.go` - LifecycleStep type, GetLifecycle function
- `internal/router/lifecycle_test.go` - TestGetLifecycle with table-driven tests

## Decisions Made

- git-commit has NextStatus=done (story is already marked done after code-review, git-commit just finalizes)
- Reused existing sentinel errors from router.go (ErrStoryComplete, ErrUnknownStatus)
- LifecycleStep uses status.Status type rather than string for type safety

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None

## Next Phase Readiness

- Lifecycle definition complete with GetLifecycle function
- Ready for Phase 7 (Story Lifecycle Executor) to use GetLifecycle for running complete story workflows
- LifecycleStep provides both workflow name and next status for the executor to update sprint-status.yaml

---

_Phase: 06-lifecycle-definition_
_Completed: 2026-01-08_
