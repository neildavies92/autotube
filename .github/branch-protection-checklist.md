# Branch Protection Checklist

Configure the `main` branch with these GitHub repository settings before merging project work:

- Require a pull request before merging.
- Require at least one approving review before merge.
- Dismiss stale approvals when new commits are pushed.
- Require branches to be up to date before merging once CI exists.
- Require a linear history.
- Allow squash merging.
- Disable merge commits.
- Disable direct pushes to `main`.
- Delete head branches automatically after merge.

Until CI exists, required status checks should be added only after the first workflows are in place.
