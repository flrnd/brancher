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
