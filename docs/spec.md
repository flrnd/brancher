Status: authoritative design document
All architectural changes must update this document.

# Brancher — Technical Specification

## Overview

Brancher is a CLI tool that creates Git branches from tasks managed in external systems such as GitHub Issues, Jira, GitLab, and similar providers.

The tool retrieves tasks from a provider and generates a branch name derived from provider task data.

Example:

```text
brancher start 42
```

Current result:

```text
branch 42-something-does-not-work created
```

Branch names are generated using configurable strategies.

## Goals

### Primary Goals

- Simple CLI workflow
- Clean branch naming from provider tasks
- Extensible provider system
- Team-friendly repository configuration
- Secure token handling via environment variables

### Non-Goals (v1)

- Complex project board workflows
- OAuth authentication
- Plugin systems
- Interactive TUI interfaces

## Core Features (v1)

- Initialize repository configuration
- Fetch tasks from a provider
- List available tasks
- Create a branch from provider task data
- Provider abstraction for future integrations

## Architecture Overview

The application is composed of several layers:

```text
CLI
 ↓
Config Loader
 ↓
Provider Factory
 ↓
Provider Registry
 ↓
Provider Implementation
 ↓
Git Driver
```

Each layer has a clear responsibility and avoids cross-layer coupling.

## Repository Structure

```text
brancher/
├─ cmd/
│  └─ brancher/
│     └─ main.go
├─ docs/
│  ├─ current.md
│  └─ spec.md
├─ internal/
│  ├─ branch/
│  │  ├─ generator.go
│  │  ├─ strategy.go
│  │  └─ title_strategy.go
│  ├─ cli/
│  │  ├─ init.go
│  │  ├─ root.go
│  │  ├─ start.go
│  │  ├─ task.go
│  │  ├─ input/
│  │  └─ output/
│  ├─ config/
│  │  └─ config.go
│  ├─ env/
│  │  └─ env.go
│  ├─ git/
│  │  ├─ driver.go
│  │  └─ repo.go
│  ├─ provider/
│  │  ├─ github/
│  │  │  ├─ client.go
│  │  │  └─ provider.go
│  │  └─ types.go
│  └─ task/
│     ├─ provider.go
│     ├─ provider_factory.go
│     ├─ provider_registry.go
│     └─ task.go
└─ pkg/
   └─ slug/
      └─ slug.go
```

## CLI Commands

### Initialize Repository

```text
brancher init
```

Creates repository configuration in `.brancher/config.yml`.

Current flow:

1. Detect Git repository
2. Detect `origin` remote
3. Parse owner and repo from the remote URL
4. Prompt for provider, owner, and repo
5. Write `.brancher/config.yml`

Current implementation note:
- remote autodetection is required today; clean fallback to fully manual owner/repo entry is still future work

### List Tasks

```text
brancher tasks
```

Lists tasks from the configured provider.

Example:

```text
12  Fix login bug
15  Something does not work
22  Improve caching
```

### Start Work on Task

```text
brancher start <task-id>
```

Creates a branch from the selected task.

Example:

```text
brancher start 15
```

Current result:

```text
15-something-does-not-work
```

Current implementation note:
- `start` currently creates the branch ref but does not check out the branch yet

## Configuration

Configuration is stored in the repository:

```text
.brancher/config.yml
```

Example:

```yaml
provider: github

project:
  owner: myorg
  repo: myrepo

branch:
  strategy: title
```

This configuration is safe to commit and intended to be shared across teams.

### Configuration Philosophy

Configuration is separated into two concerns:

| Type | Location | Purpose |
| --- | --- | --- |
| Repository configuration | `.brancher/config.yml` | Project settings |
| Secrets | Environment variables | Authentication |

Tokens are never stored in configuration files.

### Config Loading Flow

```text
Load()
  ↓
Locate .brancher/config.yml
  ↓
Read file
  ↓
Parse YAML
  ↓
Validate()
  ↓
Return Config
```

Validation fails fast if required fields are missing.

### Config Validation

Validation is rule-based and value-based.

Current required fields:

- `provider`
- `project.owner`
- `project.repo`
- `branch.strategy`

## Authentication

Tokens are provided via environment variables.

Example:

```text
export BRANCHER_GITHUB_TOKEN=xxxx
```

Provider token variables follow this pattern:

```text
BRANCHER_<PROVIDER>_TOKEN
```

Defined today:

- `BRANCHER_GITHUB_TOKEN`
- `BRANCHER_JIRA_TOKEN`
- `BRANCHER_GITLAB_TOKEN`

Current implementation notes:
- environment loading is centralized in `internal/env`
- `.env` files are loaded via `godotenv`
- only GitHub is currently wired through `env.ProviderToken(...)`

If a required variable is missing, provider construction fails with an error before the provider is created.

## Provider System

Providers supply tasks from external systems.

Examples:

- GitHub
- Jira
- GitLab

Providers convert their API data into a common internal representation.

### Task Model

All providers normalize their tasks into a common structure:

```text
Task
 ├─ ID
 ├─ Title
 ├─ Labels
 ├─ State
 └─ URL
```

Example:

```text
ID: 42
Title: Something does not work
Labels: bug
State: open
URL: https://github.com/myorg/myrepo/issues/42
```

### Provider Interface

Providers implement the interface:

```go
type Provider interface {
    Name() provider.Name
    RequiredEnv() []string
    ListTasks(ctx context.Context) ([]Task, error)
    GetTask(ctx context.Context, id string) (Task, error)
}
```

Responsibilities:

| Method | Purpose |
| --- | --- |
| `Name()` | Provider identifier |
| `RequiredEnv()` | Environment variables required for authentication |
| `ListTasks()` | Fetch available tasks |
| `GetTask()` | Fetch a specific task |

### Provider Registry

Providers register themselves via a registry.

```go
RegisterProvider(Definition)
```

Definition:

```go
type Definition struct {
    Name     provider.Name
    Required []string
    New      func(*config.Config) (Provider, error)
}
```

The registry stores provider metadata and constructors.

### Provider Factory

The provider factory creates providers dynamically.

Flow:

```text
config.provider
      ↓
lookup provider in registry
      ↓
validate required environment variables
      ↓
call constructor
      ↓
return provider instance
```

The factory does not use provider-specific switch statements.

### Provider Naming

Provider identifiers use a strongly typed name:

```go
type Name string
```

Example:

```go
const (
    GitHub Name = "github"
)
```

### Current GitHub Provider

GitHub is the only implemented provider today.

Current behavior:

- uses the GitHub Issues API
- lists open issues
- fetches an issue by ID
- maps GitHub issues into internal `task.Task`
- filters pull requests out of `ListTasks`
- rejects pull requests in `GetTask`

Current limits:

- list pagination is capped at the first 100 open issues
- no dedicated integration or e2e provider tests yet

## Git Driver

The Git driver abstracts Git operations.

Brancher uses the `go-git` library for Git operations, providing a pure Go implementation with no external dependency on the system `git` binary.

Interface:

```go
type Driver interface {
    CreateBranch(name string) error
    CreateAndCheckoutBranch(name string) error
    DeleteBranch(name string) error
    CurrentBranch() (Branch, error)
}
```

Methods:

| Method | Description |
| --- | --- |
| `CreateBranch()` | Create a local branch ref from `HEAD` |
| `CreateAndCheckoutBranch()` | Create a local branch ref from `HEAD` and switch the worktree to it |
| `DeleteBranch()` | Delete a local branch ref |
| `CurrentBranch()` | Return the current `HEAD` branch |

Initial implementation uses:

```text
GoGitDriver
```

Current implementation note:
- `CreateBranch()` is the low-level ref-only operation
- `CreateAndCheckoutBranch()` is used by `brancher start` to provide checkout behavior similar to `git checkout -b`

## Branch Generation

Branch generation must be provider-agnostic.

Different providers expose task metadata differently:

- GitHub often uses numeric task IDs such as `42`
- Jira often uses prefixed task IDs such as `PROJ-123`
- Some teams use structured titles such as `type(scope): summary`
- Other teams use plain-text titles with no enforced format

Brancher should support all of these cases.

### MVP Goal

The MVP branch naming behavior should:

- work across providers without provider-specific parsing rules
- preserve task identity in the generated branch name
- generate readable branch names from plain-text titles
- avoid requiring teams to rename or restructure existing tasks

### MVP Default Strategy

The default branch strategy is provider-agnostic and uses:

```text
<task-id>-<title-slug>
```

Examples:

```text
Task.ID    = 42
Task.Title = Something does not work
Branch     = 42-something-does-not-work
```

```text
Task.ID    = PROJ-123
Task.Title = Implement GitHub task provider
Branch     = proj-123-implement-github-task-provider
```

This default is the safest baseline because it:

- works for GitHub, Jira, GitLab, and similar providers
- preserves useful provider task identity
- supports Jira-style prefixed ticket keys naturally
- does not depend on structured title conventions

### Slug Rules

The `<title-slug>` portion must:

1. Convert to lowercase
2. Normalize unicode characters
3. Treat common separators as word boundaries
4. Collapse repeated separators into a single `-`
5. Trim leading and trailing separators

The `<task-id>` portion is also normalized to lowercase before being included in the final branch name.

Common separators currently include characters such as:

- `-`
- `_`
- `/`
- `\`
- `.`
- `,`
- `:`
- `;`
- `(`
- `)`
- `[`
- `]`
- `{`
- `}`

Example:

```text
Task.ID    = 28
Task.Title = bug(cli): Start command doesnt parse issue naming properly
Branch     = 28-bug-cli-start-command-doesnt-parse-issue-naming-properly
```

### Recommended Task Title Convention

Brancher does not require a specific task title convention.

However, for teams that want more expressive and consistent naming, the recommended title format is:

```text
type(scope): summary
```

Examples:

```text
feat(provider): implement GitHub task provider
fix(cli): print each task on its own line
refactor(branch): separate parsing from slug generation
```

This format is recommended because it:

- is familiar to teams already using Conventional Commit-style naming
- captures change intent
- preserves domain or ownership context
- can support richer branch naming strategies in the future

### Future Structured Strategies

Structured title parsing is not part of the MVP default behavior.

In future versions, Brancher may support additional branch strategies that interpret structured titles such as:

```text
type(scope): summary
```

Example input:

```text
Task.ID    = 42
Task.Title = feat(provider): implement GitHub task provider
```

Potential future output:

```text
feature/provider/42-implement-github-task-provider
```

If structured strategies are introduced later, they must fall back cleanly to the default format:

```text
<task-id>-<title-slug>
```

when the title does not match the expected structure.

## Testing and CI

The repository currently validates core behavior through unit and command tests run with:

```text
make test
```

Current CI behavior:

- GitHub Actions runs on pull requests
- GitHub Actions runs on pushes to `main`
- docs-only changes are skipped via path filters
- the CI workflow currently runs `make build` and `make test`

## Non-Goals for MVP

The MVP does not include:

- provider-specific branch naming rules
- parsing structured titles by default
- configurable branch templates
- nested scopes such as `provider/github`
- labels or metadata affecting branch names
- automatic checkout during `start`
