package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bmad-automate/internal/claude"
	"bmad-automate/internal/config"
	"bmad-automate/internal/output"
	"bmad-automate/internal/status"
	"bmad-automate/internal/workflow"
)

func TestEpicCommand_FindsAndRunsEpicStories(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  6-1-first: backlog
  6-2-second: ready-for-dev
  6-3-third: in-progress`)

	app, mockExecutor, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	// Should have run all 3 stories
	assert.Len(t, mockExecutor.RecordedPrompts, 3, "should run workflow for all three epic stories")
}

func TestEpicCommand_NumericSorting(t *testing.T) {
	tmpDir := t.TempDir()
	// Stories 1, 2, 10 should run in order 1, 2, 10 (not 1, 10, 2 alphabetically)
	createSprintStatusFile(t, tmpDir, `development_status:
  6-10-last: backlog
  6-2-middle: ready-for-dev
  6-1-first: in-progress`)

	app, mockExecutor, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	// Should have run all 3 stories in numeric order
	assert.Len(t, mockExecutor.RecordedPrompts, 3)
}

func TestEpicCommand_SkipsDoneStories(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  6-1-first: backlog
  6-2-done-story: done
  6-3-third: ready-for-dev`)

	app, mockExecutor, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	// Should have run 2 workflows (skipping the done story)
	assert.Len(t, mockExecutor.RecordedPrompts, 2, "should run workflow for 2 stories, skip done")
}

func TestEpicCommand_StopsOnFailure(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  6-1-first: backlog
  6-2-second: ready-for-dev`)

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
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)

	// Should have only run one workflow (stopped on first failure)
	assert.Len(t, mockExecutor.RecordedPrompts, 1, "should stop after first failure")
}

func TestEpicCommand_NoStoriesFound(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  7-1-other-epic: backlog`)

	app, _, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}

func TestEpicCommand_MissingSprintStatusFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Don't create sprint-status.yaml

	app, _, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}

func TestEpicCommand_FiltersOutOtherEpics(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  6-1-story: backlog
  6-2-story: ready-for-dev
  7-1-other: in-progress
  8-1-another: done`)

	app, mockExecutor, _ := setupQueueTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"epic", "6"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	// Should have run only the 2 stories from epic 6
	assert.Len(t, mockExecutor.RecordedPrompts, 2, "should only run workflows for epic 6 stories")
}
