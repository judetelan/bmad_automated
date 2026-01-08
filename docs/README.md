# bmad-automate Documentation

Comprehensive documentation for the `bmad-automate` CLI tool.

## Documentation Index

### For Users

| Document                          | Description                                   |
| --------------------------------- | --------------------------------------------- |
| [User Guide](USER_GUIDE.md)       | Getting started, installation, usage patterns |
| [CLI Reference](CLI_REFERENCE.md) | Complete command reference with examples      |

### For Developers

| Document                             | Description                            |
| ------------------------------------ | -------------------------------------- |
| [Architecture](ARCHITECTURE.md)      | System design, diagrams, data flow     |
| [Package Documentation](PACKAGES.md) | API reference for all packages         |
| [Development Guide](DEVELOPMENT.md)  | Setup, testing, extending the codebase |

### Quick Links

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands Overview](#commands-overview)
- [Configuration](#configuration)

## Installation

```bash
# Clone and build
git clone https://github.com/yourusername/bmad-automate.git
cd bmad-automate
just build

# Or install globally
go install ./cmd/bmad-automate
```

## Quick Start

```bash
# Process a story based on its status
bmad-automate run PROJ-123

# Process multiple stories
bmad-automate queue PROJ-123 PROJ-124 PROJ-125

# Process an entire epic
bmad-automate epic 05

# Run an arbitrary prompt
bmad-automate raw "What files need tests?"
```

## Commands Overview

| Command        | Purpose                     |
| -------------- | --------------------------- |
| `create-story` | Create story definition     |
| `dev-story`    | Implement a story           |
| `code-review`  | Review code changes         |
| `git-commit`   | Commit and push             |
| `run`          | Auto-route based on status  |
| `queue`        | Batch process stories       |
| `epic`         | Process all stories in epic |
| `raw`          | Execute arbitrary prompt    |

See [CLI Reference](CLI_REFERENCE.md) for complete details.

## Configuration

Default configuration file: `config/workflows.yaml`

```yaml
workflows:
  create-story:
    prompt_template: "Create story: {{.StoryKey}}"
  dev-story:
    prompt_template: "Implement story: {{.StoryKey}}"
  code-review:
    prompt_template: "Review story: {{.StoryKey}}"
  git-commit:
    prompt_template: "Commit changes for: {{.StoryKey}}"
```

See [User Guide](USER_GUIDE.md#configuration) for complete configuration options.

## Architecture Overview

```
cmd/bmad-automate/main.go
         │
         ▼
    internal/cli (Cobra commands)
         │
         ├──► internal/workflow (orchestration)
         │         │
         │         ├──► internal/claude (Claude execution)
         │         └──► internal/output (terminal formatting)
         │
         └──► internal/config (configuration)
```

See [Architecture](ARCHITECTURE.md) for detailed diagrams and explanations.

## Contributing

See [Development Guide](DEVELOPMENT.md) for:

- Development setup
- Testing practices
- Adding new commands
- Code style guidelines
