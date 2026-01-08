# User Guide

A practical guide to using `bmad-automate` for automating development workflows.

## Overview

`bmad-automate` is a CLI tool that orchestrates Claude AI to automate repetitive development tasks. It handles:

- Creating story definitions from story keys
- Implementing features based on story requirements
- Running code reviews
- Committing and pushing changes
- Processing multiple stories in batch
- Managing entire epics

## Quick Start

### Prerequisites

1. **Go 1.21+** installed
2. **Claude CLI** installed and configured ([installation guide](https://github.com/anthropics/claude-code))
3. **just** command runner (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/bmad-automate.git
cd bmad-automate

# Build the binary
just build
# OR
go build -o bmad-automate ./cmd/bmad-automate

# (Optional) Install globally
go install ./cmd/bmad-automate
```

### Verify Installation

```bash
bmad-automate --help
```

## Basic Usage

### Working with Individual Stories

#### 1. Create a Story

Generate a story definition from a story key:

```bash
bmad-automate create-story PROJ-123
```

This runs Claude with a prompt to create the story definition based on your configured template.

#### 2. Implement the Story

Run the development workflow:

```bash
bmad-automate dev-story PROJ-123
```

Claude will:

- Read the story requirements
- Implement the feature
- Run tests after each change

#### 3. Review the Code

Run code review on your changes:

```bash
bmad-automate code-review PROJ-123
```

Claude reviews the code and automatically fixes issues.

#### 4. Commit and Push

Create a commit and push to remote:

```bash
bmad-automate git-commit PROJ-123
```

### Status-Based Automation

Instead of manually running each step, let `bmad-automate` determine what to do:

```bash
bmad-automate run PROJ-123
```

The tool reads your story's status from the sprint status file and runs the appropriate workflow:

| Story Status    | Action Taken             |
| --------------- | ------------------------ |
| `backlog`       | Creates story definition |
| `ready-for-dev` | Implements story         |
| `in-progress`   | Continues implementation |
| `review`        | Runs code review         |
| `done`          | Skips (already complete) |

### Batch Processing

Process multiple stories at once:

```bash
bmad-automate queue PROJ-123 PROJ-124 PROJ-125
```

The queue:

- Processes each story based on its status
- Skips completed stories
- Stops on first failure
- Shows a summary at the end

### Processing an Epic

Run all stories in an epic:

```bash
bmad-automate epic 05
```

This finds all stories matching the pattern `05-{N}-*` (e.g., `05-01-auth`, `05-02-dashboard`) and processes them in order.

### Ad-Hoc Prompts

Run any prompt directly:

```bash
bmad-automate raw "List all TODO comments in the codebase"
bmad-automate raw "What tests are missing coverage?"
```

## Configuration

### Config File Location

By default, configuration is loaded from `config/workflows.yaml`.

Override with environment variable:

```bash
export BMAD_CONFIG_PATH=/path/to/custom/config.yaml
bmad-automate run PROJ-123
```

### Customizing Workflows

Edit `config/workflows.yaml` to customize workflow prompts:

```yaml
workflows:
  create-story:
    prompt_template: |
      Create a detailed story definition for {{.StoryKey}}.
      Include acceptance criteria and technical requirements.
      Do not ask clarifying questions.

  dev-story:
    prompt_template: |
      Implement story {{.StoryKey}}.
      Follow existing code patterns.
      Run tests after each change.
      Do not ask questions - use best judgment.

  code-review:
    prompt_template: |
      Review changes for story {{.StoryKey}}.
      Check for:
      - Code quality issues
      - Missing tests
      - Security vulnerabilities
      Auto-fix all issues immediately.

  git-commit:
    prompt_template: |
      Commit changes for {{.StoryKey}}.
      Use conventional commit format.
      Push to current branch.
```

### Template Variables

| Variable        | Description                         |
| --------------- | ----------------------------------- |
| `{{.StoryKey}}` | The story key passed to the command |

### Output Settings

Control how much output is displayed:

```yaml
output:
  truncate_lines: 20 # Max lines for tool output
  truncate_length: 60 # Max chars for command headers
```

### Claude Settings

Customize Claude execution:

```yaml
claude:
  binary_path: claude # Path to Claude binary
  output_format: stream-json # Output format (don't change)
```

Or use environment variables:

```bash
export BMAD_CLAUDE_PATH=/usr/local/bin/claude
```

## Sprint Status File

### File Location

The tool reads story status from:

```
_bmad-output/implementation-artifacts/sprint-status.yaml
```

### File Format

```yaml
development_status:
  PROJ-123: ready-for-dev
  PROJ-124: in-progress
  PROJ-125: review
  PROJ-126: done
  05-01-auth: backlog
  05-02-dashboard: ready-for-dev
```

### Valid Status Values

| Status          | Meaning                                 |
| --------------- | --------------------------------------- |
| `backlog`       | Not started, needs story creation       |
| `ready-for-dev` | Story created, ready for implementation |
| `in-progress`   | Currently being implemented             |
| `review`        | Implementation done, needs review       |
| `done`          | Complete                                |

## Workflow Patterns

### Pattern 1: Sequential Development

Run each step manually for full control:

```bash
bmad-automate create-story PROJ-123
# Review the story definition
bmad-automate dev-story PROJ-123
# Test the implementation manually
bmad-automate code-review PROJ-123
# Verify fixes
bmad-automate git-commit PROJ-123
```

### Pattern 2: Status-Driven

Let the tool figure out what to do:

```bash
# Run whatever step is needed next
bmad-automate run PROJ-123

# Run again after updating status
bmad-automate run PROJ-123
```

### Pattern 3: Batch Sprint

Process an entire sprint's stories:

```bash
bmad-automate queue SPRINT-1 SPRINT-2 SPRINT-3 SPRINT-4 SPRINT-5
```

### Pattern 4: Epic Processing

Process all stories in an epic:

```bash
bmad-automate epic 05
```

### Pattern 5: Investigation

Use raw prompts for ad-hoc tasks:

```bash
# Understand the codebase
bmad-automate raw "Explain the authentication flow"

# Find issues
bmad-automate raw "What tests have the most failures?"

# Generate reports
bmad-automate raw "Create a summary of recent changes"
```

## Understanding Output

### Tool Invocations

When Claude uses tools, you'll see formatted output:

```
┌─ Bash ─────────────────────────────────────────────────────────
│  List project files
│  $ ls -la
└────────────────────────────────────────────────────────────────

total 48
drwxr-xr-x  12 user  staff   384 Jan  8 10:00 .
drwxr-xr-x   5 user  staff   160 Jan  8 09:00 ..
...
```

### Progress Indicators

| Symbol | Meaning     |
| ------ | ----------- |
| ●      | In progress |
| ✓      | Success     |
| ✗      | Failure     |
| ○      | Skipped     |

### Queue Summary

After processing multiple stories:

```
Summary:
  PROJ-123  ✓  1m 23s
  PROJ-124  ✓  2m 45s
  PROJ-125  ○  skipped (done)
  PROJ-126  ✗  failed at dev-story (45s)
```

## Error Handling

### Exit Codes

| Code | Meaning                             |
| ---- | ----------------------------------- |
| 0    | Success                             |
| 1    | General error                       |
| N    | Claude's exit code (passed through) |

### Common Issues

**Claude not found:**

```
Error: failed to start claude: exec: "claude": executable file not found in $PATH
```

Solution: Install Claude CLI or set `BMAD_CLAUDE_PATH`.

**Config not found:**

```
Error: error loading config: open config/workflows.yaml: no such file or directory
```

Solution: Create the config file or set `BMAD_CONFIG_PATH`.

**Story not in status file:**

```
Error: story PROJ-999 not found in sprint status
```

Solution: Add the story to `sprint-status.yaml`.

**Unknown status:**

```
Error: unknown status value: invalid-status
```

Solution: Use a valid status: `backlog`, `ready-for-dev`, `in-progress`, `review`, or `done`.

## Tips and Best Practices

### 1. Start Small

Begin with single story commands to understand the workflow:

```bash
bmad-automate create-story TEST-001
```

### 2. Review Output

Always review Claude's output before moving to the next step. The streaming output shows exactly what Claude is doing.

### 3. Use Descriptive Story Keys

Story keys appear in commits and prompts. Use meaningful identifiers:

- Good: `AUTH-001`, `FEAT-user-profile`, `BUG-login-fix`
- Bad: `x`, `test`, `123`

### 4. Customize Prompts for Your Project

Edit the prompt templates to match your project's conventions, coding standards, and requirements.

### 5. Keep Status File Updated

The status file is the source of truth for automation. Keep it current as stories progress.

### 6. Handle Failures Gracefully

When a queue stops on failure:

1. Review the error
2. Fix the issue manually or adjust the story
3. Update the status file
4. Re-run the queue (completed stories will be skipped)

### 7. Use Raw for Exploration

Before starting a story, use raw prompts to understand the codebase:

```bash
bmad-automate raw "What files would I need to change to add user authentication?"
```
