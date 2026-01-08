# Codebase Structure

**Analysis Date:** 2026-01-08

## Directory Layout

```
bmad_automated/
├── cmd/
│   └── bmad-automate/     # Application entry point
│       └── main.go
├── config/                 # Default configuration
│   └── workflows.yaml
├── internal/               # Private packages
│   ├── cli/               # CLI commands and wiring
│   ├── claude/            # Claude integration
│   ├── config/            # Configuration loading
│   ├── output/            # Terminal formatting
│   └── workflow/          # Workflow orchestration
├── .planning/             # Project planning docs
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── justfile               # Task runner
├── .golangci.yml          # Linter config
├── CLAUDE.md              # Claude Code instructions
├── README.md              # User documentation
└── CONTRIBUTING.md        # Contribution guide
```

## Directory Purposes

**cmd/bmad-automate/**

- Purpose: Application entry point
- Contains: `main.go` - minimal main function
- Key files: `main.go` delegates to `cli.Execute()`
- Subdirectories: None

**internal/cli/**

- Purpose: CLI command definitions and dependency wiring
- Contains: Cobra commands, App struct, error types
- Key files:
  - `root.go` - App struct, NewRootCommand, Execute functions
  - `create_story.go`, `dev_story.go`, `code_review.go`, `git_commit.go` - Workflow commands
  - `run.go` - Full cycle command
  - `queue.go` - Batch processing command
  - `raw.go` - Custom prompt command
  - `errors.go` - ExitError type
  - `cli_test.go` - Tests
- Subdirectories: None

**internal/claude/**

- Purpose: Claude CLI integration and JSON parsing
- Contains: Executor, Parser, Event types
- Key files:
  - `client.go` - Executor interface and DefaultExecutor, MockExecutor
  - `parser.go` - Parser interface and DefaultParser
  - `types.go` - Event, StreamEvent, ContentBlock types
  - `*_test.go` - Tests
- Subdirectories: None

**internal/config/**

- Purpose: Configuration loading and prompt expansion
- Contains: Config struct, Loader, defaults
- Key files:
  - `config.go` - Loader, GetPrompt template expansion
  - `types.go` - Config, WorkflowConfig structs
  - `config_test.go` - Tests
- Subdirectories: None

**internal/output/**

- Purpose: Terminal output formatting
- Contains: Printer interface, Lipgloss styles
- Key files:
  - `printer.go` - Printer interface and DefaultPrinter
  - `styles.go` - Lipgloss style definitions, icons
  - `printer_test.go` - Tests
- Subdirectories: None

**internal/workflow/**

- Purpose: Workflow orchestration logic
- Contains: Runner, QueueRunner, Step types
- Key files:
  - `workflow.go` - Runner for single/full cycle execution
  - `queue.go` - QueueRunner for batch processing
  - `steps.go` - Step, StepResult types
  - `workflow_test.go` - Tests
- Subdirectories: None

**config/**

- Purpose: Default workflow configuration
- Contains: `workflows.yaml` with prompt templates
- Subdirectories: None

## Key File Locations

**Entry Points:**

- `cmd/bmad-automate/main.go` - Process entry point
- `internal/cli/root.go` - CLI initialization

**Configuration:**

- `go.mod` - Module definition (Go 1.25.5)
- `justfile` - Build/test tasks
- `.golangci.yml` - Linter rules
- `config/workflows.yaml` - Workflow prompts

**Core Logic:**

- `internal/cli/root.go` - Dependency wiring (App struct)
- `internal/workflow/workflow.go` - Execution orchestration
- `internal/claude/client.go` - Subprocess execution
- `internal/claude/parser.go` - JSON stream parsing

**Testing:**

- `internal/cli/cli_test.go` - CLI integration tests
- `internal/workflow/workflow_test.go` - Workflow tests
- `internal/claude/client_test.go` - Executor tests
- `internal/claude/parser_test.go` - Parser tests
- `internal/config/config_test.go` - Config tests
- `internal/output/printer_test.go` - Printer tests

**Documentation:**

- `README.md` - User guide
- `CONTRIBUTING.md` - Developer guide
- `CLAUDE.md` - Claude Code instructions

## Naming Conventions

**Files:**

- snake_case for Go files (`create_story.go`, `client_test.go`)
- UPPERCASE for important files (`README.md`, `CLAUDE.md`)
- Dot-prefixed for config (`.golangci.yml`, `.editorconfig`)

**Directories:**

- lowercase single-word (`cli`, `claude`, `config`, `output`, `workflow`)
- `internal/` for private packages (Go convention)
- `cmd/` for executables (Go convention)

**Special Patterns:**

- `*_test.go` - Test files (Go convention)
- `types.go` - Type definitions within a package
- `errors.go` - Error types within a package

## Where to Add New Code

**New CLI Command:**

- Primary code: `internal/cli/{command_name}.go`
- Tests: `internal/cli/cli_test.go` (add to existing)
- Registration: Add to `NewRootCommand()` in `internal/cli/root.go`

**New Workflow:**

- Add to `config/workflows.yaml` with prompt template
- Add constant/step in `internal/config/types.go` if needed
- No code changes needed if just adding prompts

**New Event Type:**

- Implementation: `internal/claude/types.go`
- Tests: `internal/claude/types_test.go`

**New Output Format:**

- Implementation: `internal/output/printer.go`
- Styles: `internal/output/styles.go`
- Tests: `internal/output/printer_test.go`

**Utilities:**

- Add to appropriate `internal/` package
- Avoid creating new packages unless clearly needed

## Special Directories

**.planning/**

- Purpose: Project planning documentation (GSD)
- Source: Created by planning process
- Committed: Yes

**coverage.html, coverage.out**

- Purpose: Test coverage reports
- Source: Generated by `just test-coverage`
- Committed: No (gitignored)

---

_Structure analysis: 2026-01-08_
_Update when directory structure changes_
