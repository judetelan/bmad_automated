# CLI Reference

Complete command-line interface reference for `bmad-automate`.

## Synopsis

```
bmad-automate [command] [arguments] [flags]
```

## Description

BMAD Automation CLI orchestrates Claude AI to run development workflows including story creation, implementation, code review, and git operations.

## Global Behavior

All commands:

- Load configuration from `config/workflows.yaml` (or `BMAD_CONFIG_PATH`)
- Execute Claude CLI with `--dangerously-skip-permissions` and `--output-format stream-json`
- Display styled terminal output with progress indicators
- Return appropriate exit codes (0 for success, non-zero for failure)

---

## Commands

### create-story

Create a story definition from a story key.

**Usage:**

```bash
bmad-automate create-story <story-key>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | The story identifier (e.g., `PROJ-123`) |

**Example:**

```bash
bmad-automate create-story PROJ-123
```

**Behavior:**

1. Loads `create-story` workflow prompt from configuration
2. Expands `{{.StoryKey}}` template with provided story key
3. Executes Claude with the expanded prompt
4. Displays streaming output

---

### dev-story

Implement a story by running the development workflow.

**Usage:**

```bash
bmad-automate dev-story <story-key>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | The story identifier |

**Example:**

```bash
bmad-automate dev-story PROJ-123
```

**Behavior:**

1. Loads `dev-story` workflow prompt
2. Executes Claude to implement the story
3. Claude runs tests after each implementation step

---

### code-review

Run code review on a story's changes.

**Usage:**

```bash
bmad-automate code-review <story-key>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | The story identifier |

**Example:**

```bash
bmad-automate code-review PROJ-123
```

**Behavior:**

1. Loads `code-review` workflow prompt
2. Executes Claude to review code changes
3. Automatically applies fixes when issues are found

---

### git-commit

Commit and push changes for a story.

**Usage:**

```bash
bmad-automate git-commit <story-key>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | The story identifier |

**Example:**

```bash
bmad-automate git-commit PROJ-123
```

**Behavior:**

1. Loads `git-commit` workflow prompt
2. Executes Claude to create a commit with conventional commit format
3. Pushes to the current branch

---

### run

Execute a workflow based on the story's current status in sprint-status.yaml.

**Usage:**

```bash
bmad-automate run <story-key>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | The story identifier |

**Example:**

```bash
bmad-automate run PROJ-123
```

**Status Routing:**

| Story Status    | Workflow Executed          |
| --------------- | -------------------------- |
| `backlog`       | `create-story`             |
| `ready-for-dev` | `dev-story`                |
| `in-progress`   | `dev-story`                |
| `review`        | `code-review`              |
| `done`          | No action (story complete) |

**Behavior:**

1. Reads story status from `_bmad-output/implementation-artifacts/sprint-status.yaml`
2. Routes to appropriate workflow based on status
3. Skips if story is already `done`

---

### queue

Process multiple stories in batch using status-based routing.

**Usage:**

```bash
bmad-automate queue <story-key> [story-key...]
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes | One or more story identifiers |

**Example:**

```bash
bmad-automate queue PROJ-123 PROJ-124 PROJ-125
```

**Behavior:**

1. Processes each story using the `run` command logic
2. Skips stories with status `done`
3. Stops on first failure
4. Displays summary with timing for each story

**Output:**

```
Queue: 3 stories [PROJ-123, PROJ-124, PROJ-125]

[1/3] PROJ-123
  ... workflow output ...

[2/3] PROJ-124
  ... workflow output ...

Summary:
  PROJ-123  ✓  1m 23s
  PROJ-124  ✓  2m 45s
  PROJ-125  ○  skipped (done)
```

---

### epic

Process all stories in an epic.

**Usage:**

```bash
bmad-automate epic <epic-id>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| epic-id | Yes | The epic identifier |

**Example:**

```bash
bmad-automate epic 05
```

**Story Discovery:**

Stories are discovered from `sprint-status.yaml` using the pattern:

```
{epic-id}-{story-number}-*
```

For epic `05`, this matches:

- `05-01-implement-auth`
- `05-02-add-dashboard`
- `05-03-fix-navigation`

Stories are sorted by story number and processed in order.

**Behavior:**

1. Finds all stories matching the epic pattern
2. Sorts by story number
3. Processes using queue logic (status-based routing)
4. Stops on first failure

---

### raw

Execute an arbitrary prompt with Claude.

**Usage:**

```bash
bmad-automate raw <prompt>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| prompt | Yes | The prompt text (can be multiple words) |

**Example:**

```bash
bmad-automate raw "List all Go files in the project"
bmad-automate raw Explain the architecture of this codebase
```

**Behavior:**

1. Joins all arguments into a single prompt
2. Executes Claude directly with the prompt
3. Does not use any workflow templates

---

## Exit Codes

| Code | Meaning                                              |
| ---- | ---------------------------------------------------- |
| 0    | Success                                              |
| 1    | General error (config load failure, unknown command) |
| N    | Claude exit code (passed through from Claude CLI)    |

---

## Environment Variables

| Variable           | Description                | Default                   |
| ------------------ | -------------------------- | ------------------------- |
| `BMAD_CONFIG_PATH` | Path to configuration file | `./config/workflows.yaml` |
| `BMAD_CLAUDE_PATH` | Path to Claude binary      | `claude` (from PATH)      |

---

## Configuration File

The default configuration file is `config/workflows.yaml`:

```yaml
workflows:
  create-story:
    prompt_template: "Create story: {{.StoryKey}}"

  dev-story:
    prompt_template: "Work on story: {{.StoryKey}}"

  code-review:
    prompt_template: "Review story: {{.StoryKey}}"

  git-commit:
    prompt_template: "Commit changes for {{.StoryKey}}"

full_cycle:
  steps:
    - create-story
    - dev-story
    - code-review
    - git-commit

claude:
  output_format: stream-json
  binary_path: claude

output:
  truncate_lines: 20 # Max lines to show for tool output
  truncate_length: 60 # Max chars for command header
```

### Template Variables

| Variable        | Description                         |
| --------------- | ----------------------------------- |
| `{{.StoryKey}}` | The story key passed to the command |

---

## Sprint Status File

The `run`, `queue`, and `epic` commands read story status from:

```
_bmad-output/implementation-artifacts/sprint-status.yaml
```

**Format:**

```yaml
development_status:
  PROJ-123: ready-for-dev
  PROJ-124: in-progress
  PROJ-125: done
```

**Valid Status Values:**

- `backlog` - Story not yet started
- `ready-for-dev` - Story ready for implementation
- `in-progress` - Story being implemented
- `review` - Story in code review
- `done` - Story complete

---

## Examples

### Basic Workflow

```bash
# Step-by-step workflow
bmad-automate create-story PROJ-123
bmad-automate dev-story PROJ-123
bmad-automate code-review PROJ-123
bmad-automate git-commit PROJ-123
```

### Status-Based Automation

```bash
# Let the tool determine the right workflow
bmad-automate run PROJ-123

# Process multiple stories
bmad-automate queue PROJ-123 PROJ-124 PROJ-125

# Process an entire epic
bmad-automate epic 05
```

### Ad-Hoc Tasks

```bash
# Run arbitrary prompts
bmad-automate raw "What is the test coverage?"
bmad-automate raw "Find all TODO comments"
```

### Custom Configuration

```bash
# Use custom config file
BMAD_CONFIG_PATH=/path/to/config.yaml bmad-automate run PROJ-123

# Use custom Claude binary
BMAD_CLAUDE_PATH=/usr/local/bin/claude bmad-automate dev-story PROJ-123
```
