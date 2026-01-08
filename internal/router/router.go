// Package router provides workflow routing based on story status.
package router

import (
	"errors"

	"bmad-automate/internal/status"
)

// Sentinel errors for workflow routing.
var (
	// ErrStoryComplete indicates the story is done and no workflow is needed.
	ErrStoryComplete = errors.New("story is complete, no workflow needed")

	// ErrUnknownStatus indicates the status value is not recognized.
	ErrUnknownStatus = errors.New("unknown status value")
)

// GetWorkflow returns the workflow name for the given story status.
// Returns ErrStoryComplete for done stories, ErrUnknownStatus for invalid status values.
func GetWorkflow(s status.Status) (string, error) {
	// TODO: implement
	return "", nil
}
