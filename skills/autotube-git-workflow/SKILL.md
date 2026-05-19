---
name: autotube-git-workflow
description: Apply the AutoTube repository Git workflow. Use when Codex is working in the AutoTube project and needs to create or switch branches, make commits, open pull requests, decide how to update a feature branch, or choose a merge strategy.
---

# AutoTube Git Workflow

Use trunk-based development with `main` as the protected integration branch.

## Core Rules

- Never commit or push directly to `main`.
- Create a short-lived branch before changing project files.
- Prefer branch names in the form `feat/JIRAKEY-short-kebab-description`.
- Use another conventional prefix when it fits better, such as `fix/`, `docs/`, or `chore/`.
- Keep commits in Conventional Commit style, without repeating the Jira key when it is already in the branch and PR title.
- Keep feature branches current by rebasing onto `main`; do not merge `main` into a feature branch.
- Merge pull requests into `main` with squash merge only.
- Delete the feature branch after merge.

## Pull Requests

- Require a pull request for every change that lands on `main`.
- Title pull requests as `JIRAKEY: Description of PR content`.
- Include a brief PR description covering:
  - what changed
  - why it changed
  - relevant verification performed
- Before merging, confirm:
  - CI is green
  - the PR description is present
  - Jira has been updated where relevant
  - Notion/docs have been updated where relevant
  - no secrets are included
- Use `autotube-pr-prep` for the final pre-PR review and PR body drafting workflow.

## Working Sequence

1. Confirm the current branch before editing files.
2. If on `main`, create the appropriate feature branch first.
3. Make focused commits using Conventional Commit messages.
4. Rebase the branch onto the latest `main` before opening or finalizing the PR when needed.
5. Open a PR with the required title and short description.
6. Use squash merge after review and passing checks.

## After A PR Is Merged

Clean up locally after GitHub shows the PR as merged:

```bash
git fetch --prune origin
git switch main
git pull --rebase origin main
git branch -d <merged-branch>
```

What each command is doing:

- `git fetch --prune origin` updates local knowledge of GitHub and removes stale `origin/*` remote-tracking refs for branches that no longer exist remotely.
- `git switch main` returns to the protected integration branch.
- `git pull --rebase origin main` updates local `main` from GitHub while preserving a linear local history if local commits exist.
- `git branch -d <merged-branch>` deletes the local feature branch only when Git considers it safely merged.

Use `git branch -D` only when the user explicitly confirms that an unmerged local branch should be discarded.

## Examples

- Branch: `feat/AT-13-scaffold-go-api-web`
- Commit: `feat: scaffold Go API health endpoint`
- PR title: `AT-13: Scaffold Go API and web foundation`
