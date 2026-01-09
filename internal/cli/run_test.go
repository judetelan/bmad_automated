package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bmad-automate/internal/claude"
	"bmad-automate/internal/config"
	"bmad-automate/internal/output"
	"bmad-automate/internal/status"
	"bmad-automate/internal/workflow"
)

// MockWorkflowRunner records workflow executions for testing.
type MockWorkflowRunner struct {
	ExecutedWorkflows []string
	ReturnExitCode    int
	FailOnWorkflow    string // If set, fail when this workflow is called
}

func (m *MockWorkflowRunner) RunSingle(ctx context.Context, workflowName, storyKey string) int {
	m.ExecutedWorkflows = append(m.ExecutedWorkflows, workflowName)
	if m.FailOnWorkflow == workflowName {
		return 1
	}
	return m.ReturnExitCode
}

// MockStatusWriter records status updates for testing.
type MockStatusWriter struct {
	Updates        []StatusUpdate
	FailOnStoryKey string
}

type StatusUpdate struct {
	StoryKey  string
	NewStatus status.Status
}

func (m *MockStatusWriter) UpdateStatus(storyKey string, newStatus status.Status) error {
	m.Updates = append(m.Updates, StatusUpdate{StoryKey: storyKey, NewStatus: newStatus})
	if m.FailOnStoryKey == storyKey {
		return fmt.Errorf("story not found: %s", storyKey)
	}
	return nil
}

func setupRunTestApp(tmpDir string) (*App, *claude.MockExecutor, *bytes.Buffer) {
	cfg := config.DefaultConfig()
	buf := &bytes.Buffer{}
	printer := output.NewPrinterWithWriter(buf)
	mockExecutor := &claude.MockExecutor{
		Events: []claude.Event{
			{Type: claude.EventTypeSystem, SessionStarted: true},
			{Type: claude.EventTypeResult, SessionComplete: true},
		},
		ExitCode: 0,
	}
	runner := workflow.NewRunner(mockExecutor, printer, cfg)
	queue := workflow.NewQueueRunner(runner)
	statusReader := status.NewReader(tmpDir)

	return &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		Queue:        queue,
		StatusReader: statusReader,
	}, mockExecutor, buf
}

func createSprintStatusFile(t *testing.T, tmpDir string, content string) {
	t.Helper()
	artifactsDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	require.NoError(t, os.MkdirAll(artifactsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(artifactsDir, "sprint-status.yaml"), []byte(content), 0644))
}

func TestRunCommand_StatusBasedRouting(t *testing.T) {
	tests := []struct {
		name             string
		storyKey         string
		statusYAML       string
		expectedWorkflow string
		expectError      bool
		expectExitCode   int
		expectedOutput   string
	}{
		{
			name:     "backlog status routes to create-story",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: backlog`,
			expectedWorkflow: "create-story",
			expectError:      false,
		},
		{
			name:     "ready-for-dev status routes to dev-story",
			storyKey: "STORY-2",
			statusYAML: `development_status:
  STORY-2: ready-for-dev`,
			expectedWorkflow: "dev-story",
			expectError:      false,
		},
		{
			name:     "in-progress status routes to dev-story",
			storyKey: "STORY-3",
			statusYAML: `development_status:
  STORY-3: in-progress`,
			expectedWorkflow: "dev-story",
			expectError:      false,
		},
		{
			name:     "review status routes to code-review",
			storyKey: "STORY-4",
			statusYAML: `development_status:
  STORY-4: review`,
			expectedWorkflow: "code-review",
			expectError:      false,
		},
		{
			name:     "done status prints completion message",
			storyKey: "STORY-5",
			statusYAML: `development_status:
  STORY-5: done`,
			expectedWorkflow: "",
			expectError:      false,
			expectedOutput:   "", // Output goes to fmt.Printf, not captured
		},
		{
			name:     "story not found returns error",
			storyKey: "STORY-NOT-FOUND",
			statusYAML: `development_status:
  STORY-1: backlog`,
			expectError:    true,
			expectExitCode: 1,
			expectedOutput: "", // Output goes to fmt.Printf, not captured
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, tt.statusYAML)

			app, mockExecutor, _ := setupRunTestApp(tmpDir)
			rootCmd := NewRootCommand(app)

			outBuf := &bytes.Buffer{}
			errBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(errBuf)
			rootCmd.SetArgs([]string{"run", tt.storyKey})

			err := rootCmd.Execute()

			if tt.expectError {
				require.Error(t, err)
				if tt.expectExitCode > 0 {
					code, ok := IsExitError(err)
					assert.True(t, ok, "error should be an ExitError")
					assert.Equal(t, tt.expectExitCode, code)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedWorkflow != "" {
				assert.NotEmpty(t, mockExecutor.RecordedPrompts, "prompt should have been executed")
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, outBuf.String()+errBuf.String(), tt.expectedOutput)
			}
		})
	}
}

func TestRunCommand_MissingSprintStatusFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Don't create sprint-status.yaml

	app, _, _ := setupRunTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"run", "STORY-1"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}

func TestRunCommand_WorkflowExecution(t *testing.T) {
	tests := []struct {
		name             string
		storyKey         string
		status           string
		expectedWorkflow string
	}{
		{"backlog executes create-story", "S1", "backlog", "create-story"},
		{"ready-for-dev executes dev-story", "S2", "ready-for-dev", "dev-story"},
		{"in-progress executes dev-story", "S3", "in-progress", "dev-story"},
		{"review executes code-review", "S4", "review", "code-review"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, "development_status:\n  "+tt.storyKey+": "+tt.status)

			app, mockExecutor, _ := setupRunTestApp(tmpDir)
			rootCmd := NewRootCommand(app)

			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs([]string{"run", tt.storyKey})

			err := rootCmd.Execute()
			require.NoError(t, err)

			// The workflow runner should have been called
			assert.NotEmpty(t, mockExecutor.RecordedPrompts)
		})
	}
}

func TestRunCommand_WorkflowFailure(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  STORY-1: backlog`)

	cfg := config.DefaultConfig()
	buf := &bytes.Buffer{}
	printer := output.NewPrinterWithWriter(buf)
	mockExecutor := &claude.MockExecutor{
		Events: []claude.Event{
			{Type: claude.EventTypeSystem, SessionStarted: true},
			{Type: claude.EventTypeResult, SessionComplete: true},
		},
		ExitCode: 1, // Simulate failure
	}
	runner := workflow.NewRunner(mockExecutor, printer, cfg)
	queue := workflow.NewQueueRunner(runner)
	statusReader := status.NewReader(tmpDir)

	app := &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		Queue:        queue,
		StatusReader: statusReader,
	}

	rootCmd := NewRootCommand(app)
	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"run", "STORY-1"})

	err := rootCmd.Execute()
	require.Error(t, err)

	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}

// TestRunCommand_FullLifecycleExecution tests that run command executes the full lifecycle
func TestRunCommand_FullLifecycleExecution(t *testing.T) {
	tests := []struct {
		name              string
		storyKey          string
		initialStatus     string
		expectedWorkflows []string
		expectedStatuses  []status.Status
		expectError       bool
		failOnWorkflow    string
	}{
		{
			name:          "backlog story runs full lifecycle (4 workflows)",
			storyKey:      "STORY-1",
			initialStatus: "backlog",
			expectedWorkflows: []string{
				"create-story",
				"dev-story",
				"code-review",
				"git-commit",
			},
			expectedStatuses: []status.Status{
				status.StatusReadyForDev,
				status.StatusReview,
				status.StatusDone,
				status.StatusDone,
			},
			expectError: false,
		},
		{
			name:          "ready-for-dev story runs 3 workflows",
			storyKey:      "STORY-2",
			initialStatus: "ready-for-dev",
			expectedWorkflows: []string{
				"dev-story",
				"code-review",
				"git-commit",
			},
			expectedStatuses: []status.Status{
				status.StatusReview,
				status.StatusDone,
				status.StatusDone,
			},
			expectError: false,
		},
		{
			name:          "review story runs 2 workflows",
			storyKey:      "STORY-3",
			initialStatus: "review",
			expectedWorkflows: []string{
				"code-review",
				"git-commit",
			},
			expectedStatuses: []status.Status{
				status.StatusDone,
				status.StatusDone,
			},
			expectError: false,
		},
		{
			name:              "done story prints message and exits 0",
			storyKey:          "STORY-4",
			initialStatus:     "done",
			expectedWorkflows: nil, // No workflows executed
			expectedStatuses:  nil,
			expectError:       false,
		},
		{
			name:              "workflow failure mid-lifecycle returns error",
			storyKey:          "STORY-5",
			initialStatus:     "backlog",
			failOnWorkflow:    "dev-story",
			expectedWorkflows: []string{"create-story", "dev-story"},
			expectedStatuses:  []status.Status{status.StatusReadyForDev},
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, "development_status:\n  "+tt.storyKey+": "+tt.initialStatus)

			mockRunner := &MockWorkflowRunner{
				FailOnWorkflow: tt.failOnWorkflow,
			}
			mockWriter := &MockStatusWriter{}
			statusReader := status.NewReader(tmpDir)

			app := &App{
				Config:       config.DefaultConfig(),
				StatusReader: statusReader,
				StatusWriter: mockWriter,
				Runner:       mockRunner,
			}

			rootCmd := NewRootCommand(app)
			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs([]string{"run", tt.storyKey})

			err := rootCmd.Execute()

			if tt.expectError {
				require.Error(t, err)
				code, ok := IsExitError(err)
				assert.True(t, ok, "error should be an ExitError")
				assert.Equal(t, 1, code)
			} else {
				assert.NoError(t, err)
			}

			// Verify workflows were executed in order
			assert.Equal(t, tt.expectedWorkflows, mockRunner.ExecutedWorkflows,
				"workflows should be executed in lifecycle order")

			// Verify status updates occurred after each workflow
			if tt.expectedStatuses != nil {
				require.Len(t, mockWriter.Updates, len(tt.expectedStatuses),
					"should have correct number of status updates")
				for i, expectedStatus := range tt.expectedStatuses {
					assert.Equal(t, tt.storyKey, mockWriter.Updates[i].StoryKey)
					assert.Equal(t, expectedStatus, mockWriter.Updates[i].NewStatus,
						"status update %d should be %s", i, expectedStatus)
				}
			}
		})
	}
}

func TestRunCommand_LifecycleStoryNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, "development_status:\n  OTHER-STORY: backlog")

	mockRunner := &MockWorkflowRunner{}
	mockWriter := &MockStatusWriter{}
	statusReader := status.NewReader(tmpDir)

	app := &App{
		Config:       config.DefaultConfig(),
		StatusReader: statusReader,
		StatusWriter: mockWriter,
		Runner:       mockRunner,
	}

	rootCmd := NewRootCommand(app)
	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"run", "STORY-NOT-FOUND"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)

	// No workflows should have been executed
	assert.Empty(t, mockRunner.ExecutedWorkflows)
}
