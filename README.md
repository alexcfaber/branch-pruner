# branch-pruner

Simple CLI to help prune local and remote Git branches.

Usage:

Build:

    go build ./cmd/branch-pruner

Run (dry-run):

    ./branch-pruner -dry-run -keep=3

Flags:
  -dry-run    show branches that would be deleted (default true)
  -remote     remote name (default origin)
    -keep       number of most recent branches to keep (based on commit timestamp) (default 5)
    -prefix     only consider branches with this prefix

How selection works
-------------------

branch-pruner ranks branches by the unix timestamp of their latest commit (the most recent commit on that branch). It preserves the `keep` branches with the newest timestamps and selects older branches for deletion. This is deterministic and robust across workflows where branch names aren't indicative of recency.

Notes:
- If the tool cannot read a branch's latest commit timestamp (git error), that branch is treated as very old and will be selected for deletion. This conservative behavior can be changed on request.
- Remote deletion is attempted via `git push <remote> :<branch>` but remote errors are currently ignored by default; local deletion uses `git branch -D <branch>`.

Examples
--------

# List branches that would be deleted, keeping the 3 most recently updated branches
./branch-pruner -dry-run -keep=3

# Only consider branches with prefix `feature/`
./branch-pruner -dry-run -prefix=feature/

Safety flags
------------

- `-dry-run` (default true): show branches that would be deleted but do not delete them.
- `-yes`: assume yes to all deletions; useful for scripting non-interactive runs. When set, branches shown will be deleted without prompting.
- `-interactive`: prompt before deleting each branch. Useful as an extra safety layer when running without `-dry-run`.

- `-confirm`: prompt once to confirm all deletions before they run. Useful when you want a single confirmation rather than per-branch prompts.

Examples:

# Interactive per-branch confirmation
./branch-pruner -dry-run=false -interactive

# Non-interactive destructive run (use with caution)
./branch-pruner -dry-run=false -yes

# Bulk confirm once for all deletions
./branch-pruner -dry-run=false -confirm


