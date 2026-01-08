# Architecture

**Analysis Date:** 2026-01-08

## Pattern Overview

**Overall:** Layered CLI Application with Dependency Injection

**Key Characteristics:**

- Single executable with subcommands
- Subprocess orchestration (wraps Claude CLI)
- Stateless execution model
- Event-driven streaming output
- Interface-based design for testability

## Layers

**Entry Point Layer:**

- Purpose: Process startup and delegation
- Contains: `main()` function only
- Location: `cmd/bmad-automate/main.go`
- Depends on: CLI layer
- Used by: OS process

**CLI Layer:**

- Purpose: Command parsing, dependency wiring, error handling
- Contains: Cobra command definitions, App struct for DI
- Location: `internal/cli/`
- Depends on: Workflow, Config layers
- Used by: Entry point

**Workflow Layer:**

- Purpose: Orchestrate step execution, manage timing
- Contains: `Runner` for single/full cycles, `QueueRunner` for batch
- Location: `internal/workflow/`
- Depends on: Claude, Output, Config layers
- Used by: CLI commands

**Claude Integration Layer:**

- Purpose: Execute Claude CLI, parse streaming JSON output
- Contains: `Executor` interface, `Parser`, Event types
- Location: `internal/claude/`
- Depends on: None (leaf dependency)
- Used by: Workflow layer

**Output Layer:**

- Purpose: Terminal formatting and display
- Contains: `Printer` interface, Lipgloss styling
- Location: `internal/output/`
- Depends on: None (leaf dependency)
- Used by: Workflow layer

**Config Layer:**

- Purpose: Load configuration, expand prompt templates
- Contains: `Loader`, `Config` struct, defaults
- Location: `internal/config/`
- Depends on: None (leaf dependency)
- Used by: CLI, Workflow layers

## Data Flow

**Single Workflow Execution:**

1. User runs: `bmad-automate create-story STORY-123`
2. Cobra parses args, routes to command handler (`internal/cli/create_story.go`)
3. Handler calls `Runner.RunSingle(ctx, "create-story", "STORY-123")`
4. Runner gets prompt via `Config.GetPrompt()` - expands Go template
5. Runner calls `Executor.ExecuteWithResult()` - spawns Claude subprocess
6. Parser reads streaming JSON, emits `Event` structs via channel
7. Runner calls `Printer.*` methods for each event
8. On completion, returns exit code to CLI
9. CLI returns `ExitError` or nil to Cobra

**State Management:**

- Stateless - no persistent state between commands
- All state passed via subprocess calls to Claude
- Configuration loaded fresh on each invocation

## Key Abstractions

**Executor:**

- Purpose: Run Claude CLI and stream events
- Examples: `DefaultExecutor` (real), `MockExecutor` (test)
- Location: `internal/claude/client.go`
- Pattern: Interface + implementations for testability

**Parser:**

- Purpose: Parse streaming JSON into Event structs
- Examples: `DefaultParser`
- Location: `internal/claude/parser.go`
- Pattern: Channel-based event emission

**Printer:**

- Purpose: Format terminal output
- Examples: `DefaultPrinter`
- Location: `internal/output/printer.go`
- Pattern: Interface with Writer injection for testing

**Runner:**

- Purpose: Orchestrate workflow execution
- Examples: `Runner`, `QueueRunner`
- Location: `internal/workflow/`
- Pattern: Dependency injection of Executor, Printer, Config

## Entry Points

**CLI Entry:**

- Location: `cmd/bmad-automate/main.go`
- Triggers: User runs `bmad-automate <command>`
- Responsibilities: Delegate to `cli.Execute()`

**CLI Execution:**

- Location: `internal/cli/root.go`
- Function: `Execute()` → `Run()` → `RunWithConfig()`
- Responsibilities: Load config, wire dependencies, run Cobra

**Commands (7 total):**

- `create-story` - `internal/cli/create_story.go`
- `dev-story` - `internal/cli/dev_story.go`
- `code-review` - `internal/cli/code_review.go`
- `git-commit` - `internal/cli/git_commit.go`
- `run` - `internal/cli/run.go` (full cycle)
- `queue` - `internal/cli/queue.go` (batch)
- `raw` - `internal/cli/raw.go` (custom prompt)

## Error Handling

**Strategy:** Custom ExitError type for testable exit codes

**Patterns:**

- Commands use `RunE` pattern (return error, not os.Exit)
- `ExitError` wraps exit codes for Cobra compatibility
- Errors bubble up to `Execute()` which calls `os.Exit()`

**Location:** `internal/cli/errors.go`

## Cross-Cutting Concerns

**Logging:**

- Console output via `Printer` interface
- Stderr handled separately (configurable handler)
- No structured logging (CLI tool)

**Validation:**

- Cobra handles argument validation
- Config provides sensible defaults
- No complex input validation needed

**Configuration:**

- Viper-based loading with env var overrides
- Go template expansion for prompts
- Priority: env vars > config file > defaults

---

_Architecture analysis: 2026-01-08_
_Update when major patterns change_
