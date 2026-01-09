package router

import (
	"bmad-automate/internal/status"
)

// LifecycleStep represents a single step in the story lifecycle.
// Each step contains the workflow to execute and the status to transition to after completion.
type LifecycleStep struct {
	Workflow   string
	NextStatus status.Status
}

// GetLifecycle returns the sequence of lifecycle steps from the given status to done.
// Returns ErrStoryComplete for done stories, ErrUnknownStatus for invalid status values.
func GetLifecycle(s status.Status) ([]LifecycleStep, error) {
	switch s {
	case status.StatusBacklog:
		return []LifecycleStep{
			{Workflow: "create-story", NextStatus: status.StatusReadyForDev},
			{Workflow: "dev-story", NextStatus: status.StatusReview},
			{Workflow: "code-review", NextStatus: status.StatusDone},
			{Workflow: "git-commit", NextStatus: status.StatusDone},
		}, nil
	case status.StatusReadyForDev, status.StatusInProgress:
		return []LifecycleStep{
			{Workflow: "dev-story", NextStatus: status.StatusReview},
			{Workflow: "code-review", NextStatus: status.StatusDone},
			{Workflow: "git-commit", NextStatus: status.StatusDone},
		}, nil
	case status.StatusReview:
		return []LifecycleStep{
			{Workflow: "code-review", NextStatus: status.StatusDone},
			{Workflow: "git-commit", NextStatus: status.StatusDone},
		}, nil
	case status.StatusDone:
		return nil, ErrStoryComplete
	default:
		return nil, ErrUnknownStatus
	}
}
