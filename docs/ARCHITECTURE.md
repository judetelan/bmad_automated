# Architecture Documentation

Comprehensive architecture documentation for `bmad-automate`.

## System Overview

`bmad-automate` is a Go CLI tool that orchestrates Claude AI to automate development workflows. It spawns Claude as a subprocess, parses streaming JSON output, and displays formatted results.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              bmad-automate                                  │
│                                                                             │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌────────────┐  │
│  │  CLI Layer  │───▶│   Workflow   │───▶│   Claude    │───▶│   Output   │  │
│  │   (Cobra)   │    │   (Runner)   │    │  (Executor) │    │  (Printer) │  │
│  └─────────────┘    └──────────────┘    └─────────────┘    └────────────┘  │
│         │                  │                   │                  │        │
│         ▼                  ▼                   ▼                  ▼        │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌────────────┐  │
│  │   Config    │    │    Status    │    │   Parser    │    │   Styles   │  │
│  │   (Viper)   │    │   (Reader)   │    │   (JSON)    │    │ (Lipgloss) │  │
│  └─────────────┘    └──────────────┘    └─────────────┘    └────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
                          ┌───────────────────────┐
                          │      Claude CLI       │
                          │  (External Process)   │
                          └───────────────────────┘
```

## Architecture Pattern

**Pattern:** Layered CLI Application with Dependency Injection

**Key Characteristics:**

- Single executable with subcommands
- Subprocess orchestration (wraps Claude CLI)
- Stateless execution model
- Event-driven streaming output
- Interface-based design for testability

## Layer Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     Entry Point Layer                           │
│                   cmd/bmad-automate/main.go                     │
│                       main() → cli.Execute()                    │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        CLI Layer                                │
│                      internal/cli/                              │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ App struct (Dependency Injection Container)               │  │
│  │   - Config    *config.Config                              │  │
│  │   - Executor  claude.Executor                             │  │
│  │   - Printer   output.Printer                              │  │
│  │   - Runner    *workflow.Runner                            │  │
│  │   - Queue     *workflow.QueueRunner                       │  │
│  │   - StatusReader *status.Reader                           │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  Commands: create-story, dev-story, code-review, git-commit,    │
│            run, queue, epic, raw                                │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Workflow Layer                              │
│                    internal/workflow/                           │
│                                                                 │
│  ┌─────────────────────┐     ┌─────────────────────────────┐    │
│  │  Runner             │     │  QueueRunner                │    │
│  │   - RunSingle()     │     │   - RunQueueWithStatus()    │    │
│  │   - RunRaw()        │     └─────────────────────────────┘    │
│  │   - RunFullCycle()  │                                        │
│  └─────────────────────┘                                        │
└─────────────────────────────────────────────────────────────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              ▼                   ▼                   ▼
┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐
│  Claude Layer     │  │  Output Layer     │  │  Config Layer     │
│  internal/claude/ │  │  internal/output/ │  │  internal/config/ │
│                   │  │                   │  │                   │
│  - Executor       │  │  - Printer        │  │  - Loader         │
│  - Parser         │  │  - Styles         │  │  - Config         │
│  - Event          │  │                   │  │  - GetPrompt()    │
└───────────────────┘  └───────────────────┘  └───────────────────┘
         │
         ▼
┌───────────────────────────────────────────────────────────────┐
│                    External: Claude CLI                       │
│                                                               │
│   claude --dangerously-skip-permissions                       │
│          -p "<prompt>"                                        │
│          --output-format stream-json                          │
└───────────────────────────────────────────────────────────────┘
```

## Package Dependencies

```
cmd/bmad-automate/main.go
         │
         ▼
    internal/cli (Cobra commands)
         │
         ├──► internal/workflow (orchestration)
         │         │
         │         ├──► internal/claude (Claude execution + JSON parsing)
         │         │
         │         ├──► internal/output (terminal formatting)
         │         │
         │         └──► internal/config (configuration)
         │
         ├──► internal/status (sprint status reading)
         │
         ├──► internal/router (workflow routing)
         │
         └──► internal/config (Viper configuration)
```

## Data Flow Diagram

### Single Workflow Execution

```
┌──────────────────────────────────────────────────────────────────────────┐
│  User: bmad-automate create-story PROJ-123                               │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  1. CLI Layer                                                            │
│     - Cobra parses command and arguments                                 │
│     - Routes to create-story command handler                             │
│     - Handler calls: runner.RunSingle(ctx, "create-story", "PROJ-123")   │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  2. Config Layer                                                         │
│     - config.GetPrompt("create-story", "PROJ-123")                       │
│     - Template: "/bmad:...:create-story - Create story: {{.StoryKey}}"   │
│     - Expanded: "/bmad:...:create-story - Create story: PROJ-123"        │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  3. Claude Layer                                                         │
│     - executor.ExecuteWithResult(ctx, prompt, handler)                   │
│     - Spawns: claude --dangerously-skip-permissions -p "..." ...         │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  4. Parser Layer                                                         │
│     - Reads JSON lines from stdout                                       │
│     - Converts StreamEvent → Event                                       │
│     - Emits events via channel                                           │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  5. Output Layer                                                         │
│     - handler(event) called for each event                               │
│     - printer.Text(msg) for text content                                 │
│     - printer.ToolUse(...) for tool invocations                          │
│     - printer.ToolResult(...) for tool results                           │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  6. Exit                                                                 │
│     - Claude subprocess completes                                        │
│     - Exit code propagated to CLI                                        │
│     - CLI returns ExitError or nil                                       │
└──────────────────────────────────────────────────────────────────────────┘
```

### Status-Based Routing (run command)

```
┌────────────────────────────────────────────────────────────────────────────┐
│  User: bmad-automate run PROJ-123                                          │
└────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│  1. Status Reader                                                          │
│     - Read: _bmad-output/implementation-artifacts/sprint-status.yaml       │
│     - Get status for PROJ-123: "ready-for-dev"                             │
└────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│  2. Router                                                                 │
│     - router.GetWorkflow("ready-for-dev") → "dev-story"                    │
│                                                                            │
│     Routing Table:                                                         │
│       backlog       → create-story                                         │
│       ready-for-dev → dev-story                                            │
│       in-progress   → dev-story                                            │
│       review        → code-review                                          │
│       done          → ErrStoryComplete                                     │
└────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│  3. Workflow Execution                                                     │
│     - runner.RunSingle(ctx, "dev-story", "PROJ-123")                       │
│     - (same flow as single workflow execution)                             │
└────────────────────────────────────────────────────────────────────────────┘
```

### Queue Processing

```
┌────────────────────────────────────────────────────────────────────────────┐
│  User: bmad-automate queue PROJ-123 PROJ-124 PROJ-125                      │
└────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│  QueueRunner.RunQueueWithStatus()                                          │
│                                                                            │
│  for each story:                                                           │
│    ┌────────────────────────────────────────────────────────────────────┐  │
│    │  1. Get status from sprint-status.yaml                             │  │
│    │  2. If "done" → skip                                               │  │
│    │  3. Route to workflow via router                                   │  │
│    │  4. Execute workflow                                               │  │
│    │  5. If exit code != 0 → stop queue                                 │  │
│    │  6. Record result (success/failure/skipped, duration)              │  │
│    └────────────────────────────────────────────────────────────────────┘  │
│                                                                            │
│  Print summary with all results                                            │
└────────────────────────────────────────────────────────────────────────────┘
```

## Component Diagrams

### Claude Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        DefaultExecutor                                  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ ExecuteWithResult(ctx, prompt, handler)
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  exec.CommandContext()                                                  │
│                                                                         │
│  cmd := exec.CommandContext(ctx, "claude",                              │
│      "--dangerously-skip-permissions",                                  │
│      "-p", prompt,                                                      │
│      "--output-format", "stream-json")                                  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
              ┌─────────────────────┴─────────────────────┐
              │                                           │
              ▼                                           ▼
┌─────────────────────────┐               ┌─────────────────────────────┐
│  stdout (JSON stream)   │               │  stderr (error output)      │
│                         │               │                             │
│  parser.Parse(stdout)   │               │  StderrHandler(line)        │
│          │              │               │  (logs to os.Stderr)        │
│          ▼              │               │                             │
│  chan Event             │               └─────────────────────────────┘
│          │              │
│          ▼              │
│  for event := range     │
│    handler(event)       │
└─────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  cmd.Wait()                                                             │
│  - Wait for Claude to complete                                          │
│  - Extract exit code from ExitError                                     │
│  - Return (exitCode, error)                                             │
└─────────────────────────────────────────────────────────────────────────┘
```

### Event Processing Pipeline

```
                      Claude CLI stdout
                            │
                            │ {"type":"system","subtype":"init",...}
                            │ {"type":"assistant","message":{"content":[...]}}
                            │ {"type":"user","tool_use_result":{...}}
                            │ {"type":"result",...}
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Parser.Parse()                                   │
│                                                                         │
│  bufio.Scanner reads JSON lines                                         │
│  json.Unmarshal → StreamEvent                                           │
│  NewEventFromStream → Event                                             │
└─────────────────────────────────────────────────────────────────────────┘
                            │
                            │ Event{Type, Subtype, Text, ToolName, ...}
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Event Handler                                      │
│                                                                         │
│  switch:                                                                │
│    event.IsText()       → printer.Text(event.Text)                      │
│    event.IsToolUse()    → printer.ToolUse(name, desc, cmd, path)        │
│    event.IsToolResult() → printer.ToolResult(stdout, stderr, limit)     │
└─────────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     Styled Terminal Output                              │
│                                                                         │
│  ┌─ Bash ──────────────────────────────────────────────────────────┐    │
│  │  List files in directory                                        │    │
│  │  $ ls -la                                                       │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  file1.txt                                                              │
│  file2.txt                                                              │
│  ...                                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Sequence Diagrams

### Complete Workflow Execution

```
User          CLI           Config        Runner        Executor       Parser        Printer
 │             │              │             │              │              │             │
 │─run PROJ-123│              │             │              │              │             │
 │             │              │             │              │              │             │
 │             │──GetPrompt()─▶             │              │              │             │
 │             │              │             │              │              │             │
 │             │◀─prompt──────│             │              │              │             │
 │             │              │             │              │              │             │
 │             │──RunSingle()──────────────▶│              │              │             │
 │             │              │             │              │              │             │
 │             │              │             │──SessionStart()─────────────────────────▶│
 │             │              │             │              │              │             │
 │             │              │             │──ExecuteWithResult()─▶│              │
 │             │              │             │              │              │             │
 │             │              │             │              │──Parse()────▶│             │
 │             │              │             │              │              │             │
 │             │              │             │              │◀─Event───────│             │
 │             │              │             │              │              │             │
 │             │              │             │◀─handler(event)─│              │             │
 │             │              │             │              │              │             │
 │             │              │             │──Text()───────────────────────────────────▶│
 │             │              │             │              │              │             │
 │             │              │             │◀─Event───────│              │             │
 │             │              │             │              │              │             │
 │             │              │             │──ToolUse()───────────────────────────────▶│
 │             │              │             │              │              │             │
 │             │              │             │◀─Event───────│              │              │
 │             │              │             │              │              │             │
 │             │              │             │──ToolResult()────────────────────────────▶│
 │             │              │             │              │              │             │
 │             │              │             │◀─exitCode────│              │             │
 │             │              │             │              │              │             │
 │             │              │             │──SessionEnd()────────────────────────────▶│
 │             │              │             │              │              │             │
 │             │◀─exitCode────│             │              │              │             │
 │             │              │             │              │              │             │
 │◀─Exit(code)─│              │             │              │              │             │
```

## Key Interfaces

### Executor Interface

```go
// Executor runs Claude CLI and returns streaming events.
type Executor interface {
    // Execute runs Claude with the given prompt and returns a channel of events.
    Execute(ctx context.Context, prompt string) (<-chan Event, error)

    // ExecuteWithResult runs Claude and waits for completion.
    ExecuteWithResult(ctx context.Context, prompt string, handler EventHandler) (int, error)
}
```

### Printer Interface

```go
// Printer handles terminal output formatting.
type Printer interface {
    // Session lifecycle
    SessionStart()
    SessionEnd(duration time.Duration, success bool)

    // Step progress
    StepStart(step, total int, name string)
    StepEnd(duration time.Duration, success bool)

    // Tool usage
    ToolUse(name, description, command, filePath string)
    ToolResult(stdout, stderr string, truncateLines int)

    // Content
    Text(message string)
    Divider()

    // Queue output
    QueueHeader(count int, stories []string)
    QueueStoryStart(index, total int, storyKey string)
    QueueSummary(results []StoryResult, allKeys []string, totalDuration time.Duration)
}
```

### Parser Interface

```go
// Parser reads Claude's streaming JSON output.
type Parser interface {
    Parse(reader io.Reader) <-chan Event
}
```

## Testing Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Test Setup                                     │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  MockExecutor                                                   │    │
│  │    - Events []Event     (predetermined events to return)        │    │
│  │    - ExitCode int       (exit code to return)                   │    │
│  │    - RecordedPrompts    (capture prompts for assertions)        │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Test Printer (NewPrinterWithWriter)                            │    │
│  │    - Writes to bytes.Buffer instead of os.Stdout                │    │
│  │    - Allows output capture and verification                     │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Test Config                                                    │    │
│  │    - DefaultConfig() provides sensible defaults                 │    │
│  │    - Custom configs for specific test scenarios                 │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Error Handling

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Error Flow                                      │
│                                                                        │
│  Command Handler                                                       │
│       │                                                                │
│       │  Error occurs                                                  │
│       ▼                                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Return ExitError{Code: N}                                      │   │
│  │  - Wraps exit code for Cobra compatibility                      │   │
│  │  - Implements error interface                                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                │
│       ▼                                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  RunWithConfig()                                                │   │
│  │  - Calls IsExitError() to extract code                          │   │
│  │  - Returns ExecuteResult{ExitCode, Err}                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                │
│       ▼                                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Execute()                                                      │   │
│  │  - Calls os.Exit(code) for non-zero codes                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────┘
```

## Configuration Loading

```
┌────────────────────────────────────────────────────────────────────────┐
│                     Configuration Priority                             │
│                                                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  1. Environment Variables (BMAD_*)                              │   │
│  │     - BMAD_CONFIG_PATH → custom config file                     │   │
│  │     - BMAD_CLAUDE_PATH → custom Claude binary                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                          │                                             │
│                          ▼                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  2. Config File                                                 │   │
│  │     - $BMAD_CONFIG_PATH if set                                  │   │
│  │     - OR ./config/workflows.yaml                                │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                          │                                             │
│                          ▼                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  3. Default Configuration                                       │   │
│  │     - Built-in defaults via DefaultConfig()                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                        │
│  Result: Merged Config struct                                          │
└────────────────────────────────────────────────────────────────────────┘
```

## Design Principles

1. **Dependency Injection** - All dependencies injected via App struct
2. **Interface Segregation** - Small, focused interfaces (Executor, Printer, Parser)
3. **Single Responsibility** - Each package has one clear purpose
4. **Stateless Design** - No state between command invocations
5. **Event-Driven Processing** - Stream-based handling of Claude output
6. **Testability First** - Interfaces and mocks for isolated testing
7. **Graceful Degradation** - Queue continues processing, skips completed stories
