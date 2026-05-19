---
name: autotube-ticket-implementation
description: Start and execute AutoTube implementation work from Jira tickets. Use when Codex is asked to begin, plan, or implement an AutoTube Jira issue and needs to align the ticket with Notion documentation, choose the right branch, confirm scope, and carry the work through code, tests, and documentation updates.
---

# AutoTube Ticket Implementation

Use this skill to turn a Jira ticket into a clean piece of delivery work without letting tickets, docs, and code quietly diverge.

## Workflow

1. Fetch the Jira issue by key and read the description, acceptance criteria, parent Epic, labels, and current status.
2. Fetch the relevant canonical Notion pages before coding:
   - always: `Project Home`, `Product Spec`, `Architecture`, and `Decision Log`
   - add domain pages when relevant: `Compliance Playbook`, `Operating Runbook`, `Milestone Reports`, or `MVP Readiness Checklist`
3. Compare the Jira ticket with current Notion decisions.
   - If they agree, continue.
   - If they conflict, treat that as project drift. Reconcile Jira and/or Notion first instead of building from stale instructions.
4. Define the implementation slice in one short paragraph: what changes, what stays out, and how completion will be verified.
5. Apply `autotube-git-workflow` before editing files.
   - Do not work directly on `main`.
   - Use a short-lived branch with the Jira key in the name.
6. Implement the smallest coherent change that satisfies the ticket.
7. Verify at a level proportionate to the change:
   - targeted tests for narrow work
   - broader checks when shared contracts, schema, workflows, or user-facing behavior change
8. Use `autotube-doc-sync` when implementation changes architecture, requirements, workflow, setup, or decisions.
9. Update Jira only after the work and supporting docs reflect reality.

## Working Rules

- Treat Jira as the execution record and Notion as the project knowledge record.
- Do not expand a ticket silently just because nearby improvements are tempting.
- Prefer current project decisions over old ticket wording; fix the stale source instead of preserving confusion.
- Keep ticket-sized work small enough that the PR can be reviewed as one idea.
- When requirements are missing but the intended path is obvious from current docs, make the conservative choice and record it if it becomes a durable decision.

## Practical Checks

- Does the ticket still match the current stack and project decisions?
- Is the work attached to the correct milestone and Epic?
- Is there a Notion decision or architecture page that must change with the code?
- Would another engineer understand what “done” means from the ticket plus docs alone?

## Examples

- `Implement AT-13` -> fetch AT-13, reconcile it with the Architecture page, create the correct branch, then implement the scoped foundation work.
- `Start work on the next Milestone 1 ticket` -> identify the Jira issue first, confirm the docs are current, then begin the ticket rather than starting from memory.
