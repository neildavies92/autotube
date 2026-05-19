---
name: autotube-readme-maintenance
description: Keep the AutoTube repository README accurate, concise, and useful. Use when Codex creates or updates `README.md`, when project setup or public-facing repository information changes, or when the README should be checked against current GitHub guidance and AutoTube project documentation.
---

# AutoTube README Maintenance

Use this skill to keep `README.md` as the repository front door rather than a second, stale wiki.

## Source Of Truth

- Verify GitHub's current README guidance from the official GitHub Docs page before making material structure changes:
  - `https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes`
- Use current AutoTube Notion pages for project facts:
  - `Project Home`
  - `Product Spec`
  - `Architecture`
  - `Decision Log`
- Use repository files for commands, layout, and current setup details once implementation exists.

## README Scope

Keep the README focused on the information a repository visitor needs first:
- what AutoTube is
- why it exists
- current project status
- how to get started
- repository structure
- working conventions
- where fuller documentation lives
- who maintains the project

Keep long-form product, architecture, compliance, runbook, and decision history in Notion. Link outward rather than duplicating whole sections into the README.

## Workflow

1. Fetch current GitHub README guidance when the README structure may change.
2. Read the current README, current repository layout, and relevant Notion pages.
3. Identify stale or missing public-facing information.
4. Update the smallest set of sections needed.
5. Prefer relative repository links for in-repo files.
6. Re-check claims against the codebase and Notion before finishing.
7. Use `autotube-doc-sync` if the README update reveals stale Notion documentation rather than silently papering over the mismatch.

## Maintenance Triggers

- stack or architecture changes
- setup command changes
- repository layout changes
- a new major workflow becomes usable
- support or contribution guidance changes
- the README no longer reflects the project phase

## Rules

- Do not promise features that do not yet exist.
- Do not let the README become a task tracker.
- Do not duplicate long-form Notion content just to make the README look fuller.
- Keep headings stable and scannable.
- Update examples and commands whenever they become executable reality.

## Examples

- `We added the first runnable API service` -> replace any placeholder status text with real setup and run commands.
- `The frontend moved from Next.js to React/Vite` -> update the stack section and confirm the matching Notion architecture page is also current.
