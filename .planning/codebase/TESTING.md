# Testing Patterns

**Analysis Date:** 2026-01-08

## Test Framework

**Runner:**

- Go standard `testing` package
- Run via `go test` or `just test`

**Assertion Library:**

- `github.com/stretchr/testify` v1.11.1
- `assert` for non-fatal assertions
- `require` for fatal assertions

**Run Commands:**

```bash
just test                              # Run all tests
just test-verbose                      # Run tests with verbose output
just test-pkg ./internal/claude        # Single package
just test-coverage                     # Generate coverage.html
```

## Test File Organization

**Location:**

- Co-located with source files (`foo.go` → `foo_test.go`)
- Same package (white-box testing)

**Naming:**

- `*_test.go` suffix (Go convention)
- Test functions: `Test<Type>_<Method>` or `Test<Function>`

**Structure:**

```
internal/
  cli/
    root.go
    cli_test.go
  claude/
    client.go
    client_test.go
    parser.go
    parser_test.go
    types.go
    types_test.go
  config/
    config.go
    config_test.go
  output/
    printer.go
    printer_test.go
  workflow/
    workflow.go
    workflow_test.go
```

## Test Structure

**Suite Organization:**

```go
func TestMockExecutor_Execute(t *testing.T) {
    events := []Event{
        {Type: EventTypeSystem, SessionStarted: true},
        {Type: EventTypeAssistant, Text: "Hello!"},
        {Type: EventTypeResult, SessionComplete: true},
    }

    mock := &MockExecutor{Events: events}

    ctx := context.Background()
    ch, err := mock.Execute(ctx, "test prompt")

    require.NoError(t, err)

    // Collect all events
    var collected []Event
    for event := range ch {
        collected = append(collected, event)
    }

    assert.Equal(t, events, collected)
    assert.Equal(t, []string{"test prompt"}, mock.RecordedPrompts)
}
```

**Patterns:**

- Table-driven tests for multiple cases
- Setup helpers: `setupTestRunner()`, `setupTestApp()`
- `assert` for most assertions, `require` for fatal prerequisites
- One logical assertion per test (multiple `assert` calls OK)

## Mocking

**Framework:**

- Custom mock implementations (not a mocking library)
- Mocks defined in production code, not test files

**Patterns:**

```go
// MockExecutor in internal/claude/client.go
type MockExecutor struct {
    Events          []Event
    Error           error
    ExitCode        int
    RecordedPrompts []string
}

func (m *MockExecutor) Execute(ctx context.Context, prompt string) (<-chan Event, error) {
    m.RecordedPrompts = append(m.RecordedPrompts, prompt)
    if m.Error != nil {
        return nil, m.Error
    }
    // Return pre-configured events via channel
    events := make(chan Event)
    go func() {
        defer close(events)
        for _, event := range m.Events {
            events <- event
        }
    }()
    return events, nil
}
```

**What to Mock:**

- Claude CLI execution (`MockExecutor`)
- Terminal output (`NewPrinterWithWriter(buf)`)
- Environment variables (via `os.Setenv` with `defer os.Unsetenv`)

**What NOT to Mock:**

- Pure functions (Event methods, string utilities)
- Configuration structs
- Simple type conversions

## Fixtures and Factories

**Test Data:**

```go
// Setup helpers in test files
func setupTestRunner() (*Runner, *claude.MockExecutor, *bytes.Buffer) {
    buf := &bytes.Buffer{}
    printer := output.NewPrinterWithWriter(buf)
    cfg := config.DefaultConfig()
    mockExecutor := &claude.MockExecutor{
        Events: []claude.Event{
            {Type: claude.EventTypeResult, SessionComplete: true},
        },
        ExitCode: 0,
    }
    runner := NewRunner(mockExecutor, printer, cfg)
    return runner, mockExecutor, buf
}
```

**Location:**

- Factory functions defined at top of test file
- Shared fixtures not currently used (tests are self-contained)

## Coverage

**Requirements:**

- No enforced threshold
- Coverage for awareness, not gate

**Configuration:**

- Built into Go tooling
- Output: `coverage.html`, `coverage.out`

**View Coverage:**

```bash
just test-coverage
open coverage.html
```

## Test Types

**Unit Tests:**

- All tests in codebase are unit tests
- Mock external dependencies (Claude CLI, terminal)
- Fast execution (milliseconds per test)

**Integration Tests:**

- Not currently implemented
- Would test actual Claude CLI invocation

**E2E Tests:**

- Not currently implemented
- Manual testing via actual CLI usage

## Common Patterns

**Async Testing (Channels):**

```go
func TestParser_Parse(t *testing.T) {
    parser := NewParser()
    reader := strings.NewReader(`{"type":"system"}`)

    events := parser.Parse(reader)

    // Collect from channel
    var collected []Event
    for event := range events {
        collected = append(collected, event)
    }

    assert.Len(t, collected, 1)
}
```

**Error Testing:**

```go
func TestMockExecutor_Execute_WithError(t *testing.T) {
    mock := &MockExecutor{
        Error: assert.AnError,
    }

    ctx := context.Background()
    _, err := mock.Execute(ctx, "test prompt")

    assert.Error(t, err)
}
```

**Output Capture:**

```go
func TestPrinter_SessionStart(t *testing.T) {
    var buf bytes.Buffer
    p := NewPrinterWithWriter(&buf)

    p.SessionStart()

    output := buf.String()
    assert.Contains(t, output, "Session started")
}
```

**Context Cancellation:**

```go
func TestMockExecutor_Execute_ContextCancellation(t *testing.T) {
    mock := &MockExecutor{Events: events}

    ctx, cancel := context.WithCancel(context.Background())
    ch, err := mock.Execute(ctx, "test prompt")
    require.NoError(t, err)

    // Read first event
    <-ch

    // Cancel context
    cancel()

    // Drain remaining (should close gracefully)
    for range ch {
    }
}
```

**Snapshot Testing:**

- Not used in this codebase
- Explicit assertions preferred

---

_Testing analysis: 2026-01-08_
_Update when test patterns change_
