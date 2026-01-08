# External Integrations

**Analysis Date:** 2026-01-08

## APIs & External Services

**Claude CLI (Primary Integration):**

- Claude AI CLI - Subprocess orchestration for AI-powered workflows
  - Integration method: `os/exec.Command` subprocess execution
  - Binary: `claude` (configurable via `BMAD_CLAUDE_PATH`)
  - Location: `internal/claude/client.go`
  - Protocol: Streaming JSON output (`--output-format stream-json`)
  - Flags: `--dangerously-skip-permissions`, `-p "<prompt>"`
  - Auth: Uses Claude CLI's built-in authentication

**No Other External APIs:**

- No HTTP requests to external services
- No REST/GraphQL APIs consumed
- Single integration point: Claude CLI

## Data Storage

**Databases:**

- None - This is a stateless CLI tool

**File Storage:**

- Configuration file: `config/workflows.yaml` (read-only)
- No file writes at CLI level (Claude handles file operations)

**Caching:**

- None - Each command execution is independent

## Authentication & Identity

**Auth Provider:**

- None managed by this tool
- Relies on Claude CLI's authentication (system-level)

**OAuth Integrations:**

- None

## Monitoring & Observability

**Error Tracking:**

- None - Errors output to stderr

**Analytics:**

- None

**Logs:**

- Console output via `Printer` interface
- Stderr from Claude subprocess passed through
- No structured logging or log aggregation

## CI/CD & Deployment

**Hosting:**

- Local CLI tool - not deployed to server
- Distributed as compiled binary

**CI Pipeline:**

- Not currently configured
- Manual testing via `just test`, `just lint`

## Environment Configuration

**Development:**

- Required env vars: None (all optional)
- Optional env vars:
  - `BMAD_CONFIG_PATH` - Custom config file location
  - `BMAD_CLAUDE_PATH` - Custom Claude binary path
- Mock services: `MockExecutor` for testing without Claude CLI

**Production:**

- Same binary as development
- Requires Claude CLI installed and in PATH
- Configuration via `config/workflows.yaml` or env vars

## Webhooks & Callbacks

**Incoming:**

- None - CLI tool, not a server

**Outgoing:**

- None

## Claude CLI Integration Details

**Subprocess Execution:**

```go
// From internal/claude/client.go
cmd := exec.CommandContext(ctx, e.config.BinaryPath,
    "--dangerously-skip-permissions",
    "-p", prompt,
    "--output-format", e.config.OutputFormat,
)
```

**Streaming JSON Protocol:**

- Each line is a JSON object
- Event types: `system`, `assistant`, `user`, `result`
- Parsed by `internal/claude/parser.go`
- Events emitted via channel for real-time processing

**Event Types:**

```go
// From internal/claude/types.go
EventTypeSystem    EventType = "system"    // Session lifecycle
EventTypeAssistant EventType = "assistant" // Claude responses (text, tool_use)
EventTypeUser      EventType = "user"      // Tool results
EventTypeResult    EventType = "result"    // Session completion
```

**Configuration:**

```yaml
# From config/workflows.yaml
claude:
  binary_path: claude # Path to Claude CLI
  output_format: stream-json
```

## Workflow Prompts

**Defined Workflows:**

- `create-story` - Story creation workflow
- `dev-story` - Development/implementation workflow
- `code-review` - Code review workflow
- `git-commit` - Git commit workflow

**Prompt Template Expansion:**

- Go `text/template` used for `{{.StoryKey}}` substitution
- Location: `internal/config/config.go` → `GetPrompt()`

---

_Integration audit: 2026-01-08_
_Update when adding/removing external services_
