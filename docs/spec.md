Status: authoritative design document
All architectural changes must update this document.

# Brancher — Technical Specification

## Overview

Brancher is a CLI tool that creates Git branches from tasks managed in external systems such as GitHub Issues, Jira, GitLab, etc.

The tool retrieves tasks from a provider and generates a branch name derived from the task title.

Example:

```
brancher start 42
```

Result:

```
git checkout -b something-does-not-work
```

Branch names are generated using configurable strategies.

---

# Goals

### Primary Goals

* Simple CLI workflow
* Clean branch naming from issue titles
* Extensible provider system (GitHub, Jira, GitLab, etc.)
* Team-friendly repository configuration
* Secure token handling via environment variables

### Non-Goals (v1)

* Complex project board workflows
* OAuth authentication
* Plugin systems
* Interactive TUI interfaces

---

# Core Features (v1)

* Initialize repository configuration
* Fetch tasks from a provider
* List available tasks
* Create a branch from a task title
* Provider abstraction for future integrations

---

# Architecture Overview

The application is composed of several layers:

```
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

---

# Repository Structure

```
brancher/

cmd/
  brancher/
    main.go

internal/

  cli/
    root.go
    init.go
    tasks.go
    start.go

  config/
    config.go
    loader.go

  provider/
    types.go

  task/
    task.go
    provider.go
    provider_registry.go
    provider_factory.go

  git/
    driver.go
    repo.go
    branch.go

  branch/
    generator.go
    strategy.go

providers/

  github/
    provider.go
    env.go

pkg/

  slug/
    slug.go
```

---

# CLI Commands

## Initialize Repository

```
brancher init
```

Creates repository configuration in `.brancher/config.yml`.

Steps:

1. Detect Git repository
2. Ask user for provider
3. Configure project details
4. Create `.brancher/config.yml`

---

## List Tasks

```
brancher tasks
```

Lists tasks from the configured provider.

Example:

```
12  Fix login bug
15  Something does not work
22  Improve caching
```

---

## Start Work on Task

```
brancher start <task-id>
```

Creates a branch from the selected task.

Example:

```
brancher start 15
```

Result:

```
git checkout -b something-does-not-work
```

---

# Configuration

Configuration is stored in the repository:

```
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

## Configuration Philosophy

Configuration is separated into two concerns:

| Type                     | Location                  | Purpose                       |
| ------------------------ | ------------------------- | ----------------------------- |
| Repository configuration | `.brancher/config.yml`    | Project settings              |
| Provider configuration   | `.brancher/provider.yml`  | Provider-specific settings    |
| Secrets                  | Environment variables     | Authentication                |

This prevents accidental credential leaks while allowing repository configuration to be versioned.

## Config Loading Flow

```
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

## Config Validation

Validation uses a **rule-based pointer pattern** to avoid repetitive code.

Example rule:

```
{name: "project.owner", field: &c.Project.Owner}
```

Validation loop:

```
for _, rule := range rules {
    if *rule.field == "" {
        return error
    }
}
```

---

# Authentication

Tokens are provided via environment variables.

Example:

```
export BRANCHER_GITHUB_TOKEN=xxxx
```

Provider variables follow this pattern:

```
BRANCHER_<Provider>_TOKEN
```

Examples:

```
BRANCHER_GITHUB_TOKEN
BRANCHER_JIRA_TOKEN
BRANCHER_GITLAB_TOKEN
```

If a required variable is missing, Brancher exits with an error:

```
Missing environment variable: BRANCHER_GITHUB_TOKEN
```

**Tokens are never stored in configuration files.**

---

# Provider System

Providers supply tasks from external systems.

Examples:

* GitHub
* Jira
* GitLab
* Linear
* Trello

Providers convert their API data into a common internal representation.

## Task Model

All providers normalize their tasks into a common structure:

```
Task
 ├─ ID
 ├─ Title
 ├─ Labels
 └─ State
```

Example:

```
ID: 42
Title: Something does not work
Labels: bug
```

## Provider Interface

Providers implement the interface:

```
type Provider interface {
    Name() provider.Name
    RequiredEnv() []string
    ListTasks(ctx context.Context) ([]Task, error)
    GetTask(ctx context.Context, id string) (Task, error)
}
```

Responsibilities:

| Method      | Purpose                                           |
| ----------- | ------------------------------------------------- |
| Name        | Provider identifier                               |
| RequiredEnv | Environment variables required for authentication |
| ListTasks   | Fetch available tasks                             |
| GetTask     | Fetch a specific task                             |

## Provider Registry

Providers register themselves via a registry.

```
RegisterProvider(Definition)
```

Definition:

```
type Definition struct {
    Name     provider.Name
    Required []string
    New      func(*config.Config) (Provider, error)
}
```

Registry stores provider metadata and constructors.

## Provider Factory

The provider factory creates providers dynamically.

Flow:

```
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

No switch statements are used.

## Provider Naming

Provider identifiers use a strongly typed name.

```
type Name string
```

Example:

```
const (
    GitHub Name = "github"
)
```

This avoids string duplication and typos.

## Provider Implementation Example

Example GitHub provider:

```
type GitHubProvider struct {
    token string
    owner string
    repo  string
}
```

Environment variables are defined by the provider:

```
const TokenEnv = "BRANCHER_GITHUB_TOKEN"
```

Providers self-register using `init()`.

---

# Git Driver

The Git driver abstracts Git operations.

Brancher uses the `go-git` library for Git operations, providing a pure Go implementation with no external dependencies on system Git binaries.

Interface:

```
type Driver interface {
    CreateBranch(name string) error
    DeleteBranch(name string) error
    ListLocalBranches() ([]Branch, error)
    ListRemoteBranches() ([]Branch, error)
    ListAllBranches() ([]Branch, error)
    CurrentBranch() (Branch, error)
}
```

Methods:

| Method | Description |
|--------|-------------|
| `ListLocalBranches()` | Returns only local branches |
| `ListRemoteBranches()` | Returns only remote-tracking branches |
| `ListAllBranches()` | Returns both local and remote branches |

Initial implementation uses:

```
GoGitDriver
```

---

# Branch Generation

Branch names are generated from task titles.

Example:

```
Something does not work
```

Becomes:

```
something-does-not-work
```

## Generation Rules

1. Convert to lowercase
2. Remove punctuation
3. Replace spaces with hyphens
4. Normalize unicode characters

## Branch Strategies

Strategy is configurable. Future strategies may include:

```
id-title
label-title
issue-id-title
feat/login
fix/payment-timeout
```

---

# Example Workflow

Developer clones a repository:

```
git clone repo
cd repo
```

Set token:

```
export BRANCHER_GITHUB_TOKEN=xxxx
```

List tasks:

```
brancher tasks
```

Start working:

```
brancher start 42
```

Branch created:

```
something-does-not-work
```

---

# Dependencies

Brancher aims to keep dependencies minimal.

| Dependency | Purpose           |
| ---------- | ----------------- |
| cobra      | CLI framework     |
| yaml.v3    | YAML parsing      |
| go-git     | Git operations    |

---

# Design Principles

Brancher follows these principles:

* **Simplicity first**
* **Extensible architecture**
* **Secure credential handling**
* **Minimal dependencies**
* **Strongly typed identifiers**
* **No magic strings**
* **Early validation of configuration**
* **Unix-style CLI design**

---

# Future Enhancements

## Provider Expansion

Support additional providers:

* Jira
* GitLab
* Linear
* Trello

## Interactive Task Selection

Integration with fuzzy search tools such as `fzf`.

```
brancher start
```

Displays interactive task selector.

## Smart Commit Messages

Auto-generate commit messages from tasks.

Example:

```
fix: login timeout (#42)
```

## Authentication Helpers

Possible future command:

```
brancher auth login github
```

Using OAuth or OS keychains.

## Automatic Branch Cleanup

Track and clean up stale branches created by Brancher.

## Global Configuration

Optional global configuration for user preferences.
