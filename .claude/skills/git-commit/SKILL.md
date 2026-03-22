---
name: git-commit
description: Execute git status/diff and commit changes with a well-structured commit message
userInvocable: true
---

# Git Commit Workflow

Automatically check git status, review changes, and create a well-structured commit message.

## Workflow

### 1. Check Git Status

```bash
git status
```

### 2. Review Changes

```bash
# Show staged and unstaged changes summary
git diff --stat

# Show detailed changes if needed
git diff
```

### 3. Stage Changes

```bash
# Stage all changes
git add -A

# Or stage specific files
git add <files>
```

### 4. Create Commit Message

Follow conventional commits format:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Code style (formatting, semicolons) |
| `refactor` | Code change without feature/fix |
| `perf` | Performance improvement |
| `test` | Adding/updating tests |
| `chore` | Build, CI, dependencies |

### 5. Commit

```bash
git commit -m "$(cat <<'EOF'
<type>(<scope>): <description>

## Changes
- Change 1
- Change 2
EOF
)"
```

### 6. Verify

```bash
git status
git log --oneline -3
```

## Examples

### Feature Commit

```
feat(product): add SKU management endpoints

## Backend Changes
- Add SKU CRUD operations
- Add stock lock/deduct operations
- Add SKU repository

## API Endpoints
- POST /api/v1/skus
- GET /api/v1/products/:id/skus
```

### Bug Fix Commit

```
fix(inventory): correct stock deduction logic

- Fix negative stock issue when quantity > available
- Add proper error handling for insufficient stock
```

### Documentation Commit

```
docs: add API reference documentation

- Add OpenAPI specification
- Add database schema documentation
- Add developer guide
```

## Execution Steps

When invoked, execute the following:

1. Run `git status` to see current changes
2. Run `git diff --stat` to see change summary
3. Analyze changes and determine commit type
4. Stage appropriate files with `git add`
5. Generate commit message based on changes
6. Execute `git commit`
7. Run `git status` to verify clean state
8. Show recent commits with `git log --oneline -3`

## Commit Message Guidelines

- Use imperative mood ("add" not "added")
- First line should be 50 chars or less
- Reference issues/tickets if applicable
- Group related changes in single commit
- Avoid mixing unrelated changes