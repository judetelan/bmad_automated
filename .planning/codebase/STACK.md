# Technology Stack

**Analysis Date:** 2026-01-08

## Languages

**Primary:**

- Go 1.25.5 - All application code (`go.mod`, `cmd/`, `internal/`)

**Secondary:**

- YAML - Configuration files (`config/workflows.yaml`)

## Runtime

**Environment:**

- Go 1.25.5 (specified in `go.mod`)
- Unix-like systems (macOS, Linux)

**Package Manager:**

- Go Modules
- Lockfile: `go.sum` present

## Frameworks

**Core:**

- Cobra v1.10.2 - CLI command framework (`github.com/spf13/cobra`)
- Viper v1.21.0 - Configuration management (`github.com/spf13/viper`)

**Testing:**

- Go standard `testing` package - Unit tests
- Testify v1.11.1 - Assertions and mocking (`github.com/stretchr/testify`)

**Build/Dev:**

- `just` command runner - Task automation (`justfile`)
- golangci-lint - Code linting (`.golangci.yml`)
- `go build` - Native compilation

## Key Dependencies

**Critical:**

- `github.com/spf13/cobra` v1.10.2 - CLI structure and command routing
- `github.com/spf13/viper` v1.21.0 - YAML config loading, env var overrides
- `github.com/charmbracelet/lipgloss` v1.1.0 - Terminal styling and output formatting

**Infrastructure:**

- Go standard library - `os/exec` for subprocess, `encoding/json` for parsing
- `text/template` - Go template expansion for prompts

## Configuration

**Environment:**

- Environment variables with `BMAD_` prefix for overrides
- Key vars: `BMAD_CONFIG_PATH`, `BMAD_CLAUDE_PATH`

**Build:**

- `go.mod` - Module definition
- `justfile` - Task runner configuration
- `.golangci.yml` - Linter configuration (5m timeout)

## Platform Requirements

**Development:**

- macOS/Linux/Windows (any platform with Go)
- `just` command runner (optional, for task automation)
- golangci-lint (optional, for linting)

**Production:**

- Compiled binary (`./bmad-automate`)
- Requires Claude CLI installed and in PATH
- No external runtime dependencies

---

_Stack analysis: 2026-01-08_
_Update after major dependency changes_
