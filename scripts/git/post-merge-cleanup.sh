#!/bin/sh
set -eu

remote="${REMOTE:-origin}"
base_branch="${BASE_BRANCH:-main}"
branch="${1:-}"

usage() {
	printf '%s\n' "Usage: make post-merge-cleanup BRANCH=<merged-branch>"
	printf '%s\n' "       sh scripts/git/post-merge-cleanup.sh <merged-branch>"
	printf '\n%s\n' "If BRANCH is omitted, the current branch is used when it is not ${base_branch}."
}

die() {
	printf 'error: %s\n' "$1" >&2
	exit 1
}

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
	die "not inside a Git worktree"
fi

current_branch="$(git branch --show-current)"

if [ -z "$branch" ]; then
	if [ -z "$current_branch" ] || [ "$current_branch" = "$base_branch" ]; then
		usage
		exit 2
	fi

	branch="$current_branch"
fi

if [ "$branch" = "$base_branch" ]; then
	die "refusing to delete ${base_branch}"
fi

if ! git show-ref --verify --quiet "refs/heads/${branch}"; then
	die "local branch ${branch} does not exist"
fi

if [ -n "$(git status --porcelain)" ]; then
	die "worktree has uncommitted changes; commit, stash, or discard them before cleanup"
fi

printf 'Fetching and pruning %s...\n' "$remote"
git fetch --prune "$remote"

printf 'Switching to %s...\n' "$base_branch"
git switch "$base_branch"

printf 'Updating %s from %s/%s...\n' "$base_branch" "$remote" "$base_branch"
git pull --rebase "$remote" "$base_branch"

printf 'Deleting local branch %s...\n' "$branch"
if git branch -d "$branch"; then
	printf 'Post-merge cleanup complete.\n'
	exit 0
fi

printf '\n%s\n' "Git does not consider ${branch} fully merged into local ${base_branch}."
printf '%s\n' "This is common after a squash merge, so checking GitHub before force-deleting..."

if ! command -v gh >/dev/null 2>&1; then
	die "GitHub CLI is not available; verify the PR is merged before running: git branch -D ${branch}"
fi

pr_json="$(gh pr view "$branch" --json baseRefName,mergedAt,state,url 2>/dev/null || true)"

if [ -z "$pr_json" ]; then
	die "could not confirm a merged GitHub PR for ${branch}; verify manually before force deletion"
fi

case "$pr_json" in
	*"\"state\":\"MERGED\""*) pr_is_merged="true" ;;
	*) pr_is_merged="false" ;;
esac

case "$pr_json" in
	*"\"baseRefName\":\"${base_branch}\""*) pr_targets_base="true" ;;
	*) pr_targets_base="false" ;;
esac

if [ "$pr_is_merged" = "true" ] && [ "$pr_targets_base" = "true" ]; then
	printf '%s\n' "$pr_json"
	printf 'GitHub confirms the PR was merged into %s; deleting local branch with -D...\n' "$base_branch"
	git branch -D "$branch"
	printf 'Post-merge cleanup complete.\n'
	exit 0
fi

printf '%s\n' "$pr_json"
die "GitHub did not confirm ${branch} is merged into ${base_branch}; local branch was kept"
