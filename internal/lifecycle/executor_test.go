package lifecycle

import (
	"context"
	"errors"
	"testing"

	"bmad-automate/internal/router"
	"bmad-automate/internal/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockWorkflowRunner implements WorkflowRunner for testing.
type MockWorkflowRunner struct {
	// RunSingleFunc allows tests to control workflow execution behavior.
	RunSingleFunc func(ctx context.Context, workflowName, storyKey string) int
	// Calls records all RunSingle calls for verification.
	Calls []struct {
		WorkflowName string
		StoryKey     string
	}
}

func (m *MockWorkflowRunner) RunSingle(ctx context.Context, workflowName, storyKey string) int {
	m.Calls = append(m.Calls, struct {
		WorkflowName string
		StoryKey     string
	}{workflowName, storyKey})

	if m.RunSingleFunc != nil {
		return m.RunSingleFunc(ctx, workflowName, storyKey)
	}
	return 0 // success by default
}

// MockStatusReader implements StatusReader for testing.
type MockStatusReader struct {
	// GetStoryStatusFunc allows tests to control status reading behavior.
	GetStoryStatusFunc func(storyKey string) (status.Status, error)
}

func (m *MockStatusReader) GetStoryStatus(storyKey string) (status.Status, error) {
	if m.GetStoryStatusFunc != nil {
		return m.GetStoryStatusFunc(storyKey)
	}
	return status.StatusBacklog, nil
}

// MockStatusWriter implements StatusWriter for testing.
type MockStatusWriter struct {
	// UpdateStatusFunc allows tests to control status writing behavior.
	UpdateStatusFunc func(storyKey string, newStatus status.Status) error
	// Calls records all UpdateStatus calls for verification.
	Calls []struct {
		StoryKey  string
		NewStatus status.Status
	}
}

func (m *MockStatusWriter) UpdateStatus(storyKey string, newStatus status.Status) error {
	m.Calls = append(m.Calls, struct {
		StoryKey  string
		NewStatus status.Status
	}{storyKey, newStatus})

	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(storyKey, newStatus)
	}
	return nil // success by default
}

func TestNewExecutor(t *testing.T) {
	runner := &MockWorkflowRunner{}
	reader := &MockStatusReader{}
	writer := &MockStatusWriter{}

	executor := NewExecutor(runner, reader, writer)

	assert.NotNil(t, executor)
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		storyKey       string
		currentStatus  status.Status
		getStatusErr   error
		workflowResult int
		updateErr      error
		wantErr        error
		wantWorkflows  []string
		wantUpdates    []status.Status
	}{
		{
			name:          "story in backlog runs full lifecycle",
			storyKey:      "EPIC-1-story",
			currentStatus: status.StatusBacklog,
			wantWorkflows: []string{"create-story", "dev-story", "code-review", "git-commit"},
			wantUpdates:   []status.Status{status.StatusReadyForDev, status.StatusReview, status.StatusDone, status.StatusDone},
		},
		{
			name:          "story in ready-for-dev skips create-story",
			storyKey:      "EPIC-1-story",
			currentStatus: status.StatusReadyForDev,
			wantWorkflows: []string{"dev-story", "code-review", "git-commit"},
			wantUpdates:   []status.Status{status.StatusReview, status.StatusDone, status.StatusDone},
		},
		{
			name:          "story in review runs code-review and git-commit",
			storyKey:      "EPIC-1-story",
			currentStatus: status.StatusReview,
			wantWorkflows: []string{"code-review", "git-commit"},
			wantUpdates:   []status.Status{status.StatusDone, status.StatusDone},
		},
		{
			name:          "story already done returns ErrStoryComplete",
			storyKey:      "EPIC-1-story",
			currentStatus: status.StatusDone,
			wantErr:       router.ErrStoryComplete,
			wantWorkflows: nil,
			wantUpdates:   nil,
		},
		{
			name:         "get status error propagates",
			storyKey:     "EPIC-1-story",
			getStatusErr: errors.New("file not found"),
			wantErr:      errors.New("file not found"),
		},
		{
			name:           "workflow failure stops execution",
			storyKey:       "EPIC-1-story",
			currentStatus:  status.StatusBacklog,
			workflowResult: 1, // first workflow fails
			wantWorkflows:  []string{"create-story"},
			wantUpdates:    nil, // no status updates if workflow fails
		},
		{
			name:          "status update failure stops execution",
			storyKey:      "EPIC-1-story",
			currentStatus: status.StatusBacklog,
			updateErr:     errors.New("write failed"),
			wantWorkflows: []string{"create-story"},                  // workflow runs
			wantUpdates:   []status.Status{status.StatusReadyForDev}, // update attempted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockWorkflowRunner{
				RunSingleFunc: func(ctx context.Context, workflowName, storyKey string) int {
					return tt.workflowResult
				},
			}
			reader := &MockStatusReader{
				GetStoryStatusFunc: func(storyKey string) (status.Status, error) {
					if tt.getStatusErr != nil {
						return "", tt.getStatusErr
					}
					return tt.currentStatus, nil
				},
			}
			writer := &MockStatusWriter{
				UpdateStatusFunc: func(storyKey string, newStatus status.Status) error {
					return tt.updateErr
				},
			}

			executor := NewExecutor(runner, reader, writer)
			err := executor.Execute(context.Background(), tt.storyKey)

			// Check error
			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, router.ErrStoryComplete) {
					assert.ErrorIs(t, err, router.ErrStoryComplete)
				} else {
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
			} else if tt.workflowResult != 0 {
				// Workflow failure should return an error
				require.Error(t, err)
				assert.Contains(t, err.Error(), "workflow failed")
			} else if tt.updateErr != nil {
				// Update failure should return an error
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.updateErr.Error())
			} else {
				require.NoError(t, err)
			}

			// Check workflow calls
			if tt.wantWorkflows != nil {
				require.Len(t, runner.Calls, len(tt.wantWorkflows))
				for i, wantWorkflow := range tt.wantWorkflows {
					assert.Equal(t, wantWorkflow, runner.Calls[i].WorkflowName)
					assert.Equal(t, tt.storyKey, runner.Calls[i].StoryKey)
				}
			}

			// Check status updates
			if tt.wantUpdates != nil {
				require.Len(t, writer.Calls, len(tt.wantUpdates))
				for i, wantStatus := range tt.wantUpdates {
					assert.Equal(t, tt.storyKey, writer.Calls[i].StoryKey)
					assert.Equal(t, wantStatus, writer.Calls[i].NewStatus)
				}
			}
		})
	}
}
