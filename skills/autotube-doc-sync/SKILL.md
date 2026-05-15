---
name: autotube-doc-sync
description: Keep AutoTube project documentation aligned with reality. Use when Codex changes or reviews architecture, requirements, workflows, operations, project decisions, milestone status, or anything else that should be reflected in the AutoTube Notion wiki.
---

# AutoTube Doc Sync

Use this skill whenever the codebase or project decisions move far enough that future-you would be misled by yesterday's docs.

## Canonical Pages

- `Project Home`: current milestone, key links, active risks, documentation discipline.
- `Product Spec`: goals, learning goals, audience, v1 scope, out-of-scope items, success metrics.
- `Architecture`: stack, modules, data flow, storage model, repo layout, integration choices.
- `Compliance Playbook`: originality policy, review gates, rejection rules.
- `Operating Runbook`: weekly process, release process, review cadence.
- `Decision Log`: dated durable decisions with rationale.
- `Milestone Reports`: fortnightly summaries and upcoming focus.
- `MVP Readiness Checklist`: readiness criteria that should remain true before real operation.

## Workflow

1. Identify what changed:
   - decision
   - requirement
   - architecture
   - workflow
   - operating practice
   - milestone outcome
2. Fetch the current canonical page or pages before editing.
3. Update the smallest set of pages needed to keep the wiki truthful.
4. Add a `Decision Log` entry whenever a durable choice was made or reversed.
5. Update `Milestone Reports` when a fortnightly checkpoint closes or the milestone plan meaningfully changes.
6. Re-check for stale terms, duplicated guidance, or contradictions after editing.
7. Prefer one canonical answer over several nearly matching pages.

## Mapping Changes To Pages

- Stack or module change -> `Architecture` and usually `Decision Log`
- Scope or success metric change -> `Product Spec`
- Review gate or originality rule change -> `Compliance Playbook` and usually `Decision Log`
- Weekly delivery process change -> `Operating Runbook`
- Milestone completion or replanning -> `Milestone Reports` and maybe `Project Home`
- Readiness requirement change -> `MVP Readiness Checklist`

## Rules

- Jira remains the task tracker; do not mirror ticket status into Notion unless summarizing a milestone.
- Notion remains the durable source of truth for project knowledge; do not leave known contradictions in place.
- If a major implementation choice changes, update docs in the same delivery cycle rather than promising to “come back later.”
- Keep current pages canonical. Archive or clearly mark superseded duplicates instead of letting two live copies drift apart.

## Examples

- `We changed from SQLite to PostgreSQL` -> update `Architecture`, add a `Decision Log` entry, and re-check for stale stack references.
- `We changed how weekly review works` -> update `Operating Runbook`, possibly `Project Home`, and add a decision if the rule is durable.
- PDF skill: `fill_fillable_fields.py`, `extract_form_field_info.py` - utilities for PDF manipulation
- DOCX skill: `document.py`, `utilities.py` - Python modules for document processing

**Appropriate for:** Python scripts, shell scripts, or any executable code that performs automation, data processing, or specific operations.

**Note:** Scripts may be executed without loading into context, but can still be read by Codex for patching or environment adjustments.

### references/
Documentation and reference material intended to be loaded into context to inform Codex's process and thinking.

**Examples from other skills:**
- Product management: `communication.md`, `context_building.md` - detailed workflow guides
- BigQuery: API reference documentation and query examples
- Finance: Schema documentation, company policies

**Appropriate for:** In-depth documentation, API references, database schemas, comprehensive guides, or any detailed information that Codex should reference while working.

### assets/
Files not intended to be loaded into context, but rather used within the output Codex produces.

**Examples from other skills:**
- Brand styling: PowerPoint template files (.pptx), logo files
- Frontend builder: HTML/React boilerplate project directories
- Typography: Font files (.ttf, .woff2)

**Appropriate for:** Templates, boilerplate code, document templates, images, icons, fonts, or any files meant to be copied or used in the final output.

---

**Not every skill requires all three types of resources.**
