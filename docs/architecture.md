# Brancher – Project Design Document

## Overview

**Brancher** is a CLI tool that creates Git branches from tasks in project management systems such as GitHub Issues, Jira, GitLab, etc.

The goal is to streamline the developer workflow by allowing developers to start work on tasks directly from their project board while automatically generating clean branch names.

Example:

```
brancher start 42
```

Result:

```
git checkout -b something-does-not-work
```

Branch names are derived from the task title using a configurable strategy.

---

# Goals

### Primary Goals

* Simple CLI workflow
* Clean branch naming from issue titles
* Extensible provider system (GitHub, Jira, GitLab, etc.)
* Team-friendly repository configuration
* Secure token handling

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

# CLI Commands

## Initialize Repository

```
brancher init
```

Creates repository configuration in:

```
.brancher/config.yml
```
Config structure to be defined.

---

## List Tasks

```
brancher tasks
```

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

Example:

```
brancher start 15
```

Result:

```
git checkout -b something-does-not-work
```

---

# Repository Configuration

Brancher stores project configuration inside the repository.

```
.brancher/
  config.yml
```

Example (still to be defined):

```yaml
provider: github

project:
  owner: myorg
  repo: myrepo

branch:
  strategy: title
```

This configuration is safe to commit and intended to be shared across teams.

---

# Token Handling

Brancher **does not store tokens in configuration files**.

Authentication is performed exclusively via environment variables.

Example:

```
export BRANCHER_GITHUB_TOKEN=xxxx
```

Provider variables follow this pattern:

```
BRANCHER_<Provider>_TOKEN
```

Example:

```
BRANCHER_GITHUB_TOKEN
BRANCHER_JIRA_TOKEN
BRANCHER_GITLAB_TOKEN
```

If a required variable is missing, Brancher will exit with an error.

Example:

```
Missing environment variable: BRANCHER_GITHUB_TOKEN
```

---

# Configuration Philosophy

Configuration is separated into two concerns:

| Type                     | Location               | Purpose                       |
| ------------------------ | ---------------------- | ----------------------------- |
| Repository configuration | `.brancher/config.yml` | project settings              |
| Providers configuration  | `.brancher/provider.yml| provider settings
| Secrets                  | Environment variables  | Authentication                |

This prevents accidental credential leaks while allowing repository configuration to be versioned.

---

# Provider Architecture

Brancher uses a **provider abstraction** to support multiple project management systems.

Providers convert their API data into a common internal representation.

Example providers:

* GitHub
* Jira
* GitLab
* Linear
* Trello

---

# Task Model

All providers normalize their tasks into a common structure.
(Still to be defined)

Example:

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

---

# Provider Interface (Conceptual)

Each provider implements the following behavior:

```
Provider
 ├─ Name()
 ├─ RequiredEnv()
 ├─ ListTasks()
 └─ GetTask()
```

Responsibilities:

| Method      | Purpose                                           |
| ----------- | ------------------------------------------------- |
| Name        | Provider identifier                               |
| RequiredEnv | Environment variables required for authentication |
| ListTasks   | Fetch available tasks                             |
| GetTask     | Fetch a specific task                             |

---

# Branch Naming

Branch names are generated from the task title.

Example:

```
[InternalCode-XX] Something does not work
```

Becomes:

```
internal-code-xx-something-does-not-work
```

Generation rules:

1. Convert to lowercase
2. Remove punctuation
3. Replace spaces with hyphens
4. Normalize unicode characters

Future strategies may include:

```
id-title
label-title
issue-id-title
```

---

# Git Integration

Brancher uses the `go-git` library for Git operations.

Example command:

```
git checkout -b <branch>
```

Using `go-git` provides a pure Go implementation with no external dependencies on system Git binaries.

---

# Project Structure

```
brancher/

cmd/
  brancher/

internal/

  cli/
    root.go
    init.go
    tasks.go
    start.go

  config/
    config.go

  git/
    repo.go
    branch.go

  branch/
    generator.go
    slug.go

  task/
    task.go

providers/

  github/
    provider.go

  jira/
    provider.go
```

---

# Dependencies

Brancher aims to keep dependencies minimal.

Expected dependencies:

| Dependency | Purpose           |
| ---------- | ----------------- |
| cobra      | CLI framework     |
| Viper      | config framework  |
| go-github  | GitHub API client |
| go-git     | Git operations    |

Git operations are executed using the `go-git` library rather than system Git binary calls.

---

# Initialization Workflow

```
brancher init
```

Steps:

1. Detect Git repository
2. Ask user for provider
3. Configure project details
4. Create `.brancher/config.yml`
5. Create `.brancher/provider.yml`

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

# Future Enhancements

Potential improvements:

### Provider Expansion

Support additional providers:

* Jira
* GitLab
* Linear
* Trello

---

### Branch Strategies

Additional naming strategies:

```
feat/login
fix/payment-timeout
```

---

### Interactive Task Selection

Integration with fuzzy search tools such as `fzf`.

```
brancher start
```

Displays interactive task selector.

---

### Smart Commit Messages

Auto-generate commit messages from tasks.

Example:

```
fix: login timeout (#42)
```

---

### Authentication Helpers

Possible future command:

```
brancher auth login github
```

Using OAuth or OS keychains.

---

# Design Principles

Brancher is built around the following principles:

* **Simplicity first**
* **Extensible architecture**
* **Secure credential handling**
* **Minimal dependencies**
* **Unix-style CLI design**

---

# Status

Design phase complete.

Next steps:

1. Implement CLI skeleton
2. Implement configuration loading
3. Implement GitHub provider
4. Implement branch generation
5. Implement `start` command
