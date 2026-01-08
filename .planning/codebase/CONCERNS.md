# Codebase Concerns

**Analysis Date:** 2026-01-08

## Tech Debt

**Scanner error intentionally ignored:**

- Issue: `scanner.Err()` is not checked after scanning in parser
- Location: `internal/claude/parser.go:63-64`
- Why: Comment states "intentionally not checked... to gracefully handle EOF and pipe closure"
- Impact: Potential silent failures if scanner encounters non-EOF errors
- Fix approach: Log scanner errors at debug level, or return them via a separate error channel

**Duplicate StepResult type definition:**

- Issue: `StepResult` is defined in both `internal/workflow/steps.go` and `internal/output/printer.go`
- Files: `internal/workflow/steps.go:13-18`, `internal/output/printer.go:11-16`
- Why: Likely evolved independently
- Impact: Type conversion needed between layers, potential confusion
- Fix approach: Use single definition in one package, import in other

## Known Bugs

**None currently identified:**

- Codebase is relatively new and clean
- Test coverage exists for critical paths

## Security Considerations

**`--dangerously-skip-permissions` flag:**

- Risk: Claude CLI runs without permission prompts
- Location: `internal/claude/client.go` (hardcoded in command args)
- Current mitigation: This is intentional for automation - prompts would break batch processing
- Recommendations: Document this clearly for users, consider making it configurable

**Prompt injection potential:**

- Risk: User-provided story keys are embedded in prompts without sanitization
- Location: `internal/config/config.go` → `GetPrompt()`
- Current mitigation: Go template escaping provides some protection
- Recommendations: Story keys should be validated (alphanumeric, limited length)

## Performance Bottlenecks

**None significant:**

- CLI tool with straightforward execution
- JSON parsing is streaming (not loading entire output into memory)
- No database queries or network calls (beyond Claude subprocess)

**Potential concern - large buffer allocation:**

- Location: `internal/claude/parser.go:44`
- Current: 1MB initial buffer, 10MB max
- Impact: Memory allocation per command execution
- Note: Appropriate for expected Claude output sizes

## Fragile Areas

**Event type conversion:**

- Location: `internal/claude/types.go` → `NewEventFromStream()`
- Why fragile: Complex switch on nested JSON structure from Claude
- Common failures: Claude output format changes would break parsing
- Safe modification: Add tests for edge cases before changing
- Test coverage: Good (`internal/claude/types_test.go`)

**Configuration loading order:**

- Location: `internal/config/config.go` → `Load()`
- Why fragile: Multiple config sources (file, env, defaults) with priority
- Common failures: Unexpected config values when env vars set
- Safe modification: Test with various env var combinations
- Test coverage: Good (`internal/config/config_test.go`)

## Scaling Limits

**Not applicable:**

- Single-user CLI tool
- No concurrent usage scenarios
- Claude CLI handles actual AI processing

## Dependencies at Risk

**No concerning dependencies:**

- Cobra, Viper, Testify: Well-maintained, stable
- Lipgloss: Active Charmbracelet project
- All dependencies are popular Go ecosystem choices

## Missing Critical Features

**No graceful shutdown on context cancellation:**

- Problem: Context cancellation during Claude execution may leave subprocess orphaned
- Location: `internal/claude/client.go`
- Current workaround: Process typically completes or user kills manually
- Blocks: Clean interruption with Ctrl+C
- Implementation complexity: Low (need to kill subprocess on ctx.Done())

**No retry mechanism:**

- Problem: Transient failures (Claude CLI startup, network) cause immediate failure
- Current workaround: User re-runs command
- Blocks: Reliable automation in batch mode
- Implementation complexity: Medium (exponential backoff, configurable retries)

## Test Coverage Gaps

**CLI integration tests limited:**

- What's not tested: Full end-to-end command execution with real Cobra
- Files: `internal/cli/cli_test.go` uses mocks throughout
- Risk: Cobra command registration issues not caught
- Priority: Low (Cobra is well-tested)
- Difficulty to test: Would need to capture stdout/stderr from subprocess

**No integration tests with real Claude CLI:**

- What's not tested: Actual subprocess execution, real JSON parsing
- Risk: Claude output format changes break parsing silently
- Priority: Medium
- Difficulty to test: Requires Claude CLI installed, slow execution

---

_Concerns audit: 2026-01-08_
_Update as issues are fixed or new ones discovered_
