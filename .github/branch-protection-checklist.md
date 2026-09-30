# Branch protection

Configure `main` to require pull requests, an approving review, linear history and the `checks` job from the CI workflow. Dismiss stale reviews and require checks on the current PR revision. Prefer squash merges after review and disable direct pushes to `main`.

This is repository administration guidance; SWA-46 adds CI but does not change GitHub protection settings or merge any PR.
