---
name: autotube-pr-prep
description: Prepare AutoTube pull requests for review. Use when Codex is asked to open, draft, finalize, or check a PR for the AutoTube repository and needs to apply the Git workflow, verify tests and documentation, confirm Jira alignment, and produce the required PR title and description.
---

# AutoTube PR Prep

Use this skill at the point where implementation becomes reviewable project history.

## Before The PR

1. Apply `autotube-git-workflow`.
2. Confirm the branch is not `main` and follows the agreed naming pattern.
3. Inspect the diff and make sure the branch contains one coherent unit of work.
4. Confirm the linked Jira issue is correct and no requirement drift remains unresolved.
5. Confirm relevant Notion pages were updated when architecture, requirements, workflow, setup, or durable decisions changed.
6. Run the verification appropriate to the change and capture what was run.
7. Check that no secrets or accidental local files are included.

## PR Shape

- Title format: `JIRAKEY: Description of PR content`
- Body should briefly cover:
  - what changed
  - why it changed
  - verification performed
  - documentation updates, when relevant
  - linked Jira issue
- Prefer a draft PR until the branch is genuinely ready for review.

## Ready-To-Open Checklist

- Branch name follows project convention.
- Commit messages follow Conventional Commit style.
- Relevant tests/checks have run.
- Jira reflects the implemented scope.
- Notion is current where it should be.
- PR title and body are present and specific.
- The branch is rebased on current `main` when needed.

## After Review

- Use squash merge into `main`.
- Delete the feature branch after merge.
- If review materially changes behavior or decisions, re-run `autotube-doc-sync` before merge.

## Example

- Branch: `feat/AT-13-scaffold-go-api-web`
- Commit: `feat: scaffold Go API health endpoint`
- PR title: `AT-13: Scaffold Go API and web foundation`
