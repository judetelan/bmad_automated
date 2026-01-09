package lifecycle

import (
	"context"

	"bmad-automate/internal/status"
)

// WorkflowRunner is the interface for running workflows.
type WorkflowRunner interface {
	RunSingle(ctx context.Context, workflowName, storyKey string) int
}

// StatusReader is the interface for reading story status.
type StatusReader interface {
	GetStoryStatus(storyKey string) (status.Status, error)
}

// StatusWriter is the interface for updating story status.
type StatusWriter interface {
	UpdateStatus(storyKey string, newStatus status.Status) error
}

// Executor runs the complete lifecycle for a story.
type Executor struct {
	runner       WorkflowRunner
	statusReader StatusReader
	statusWriter StatusWriter
}

// NewExecutor creates a new Executor with the given dependencies.
func NewExecutor(runner WorkflowRunner, reader StatusReader, writer StatusWriter) *Executor {
	return &Executor{
		runner:       runner,
		statusReader: reader,
		statusWriter: writer,
	}
}

// Execute runs the complete lifecycle for a story from its current status to done.
func (e *Executor) Execute(ctx context.Context, storyKey string) error {
	// TODO: Implement
	return nil
}
