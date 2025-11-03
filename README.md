# branch-pruner

Simple CLI to help prune local and remote Git branches.

Usage:

Build:

    go build ./cmd/branch-pruner

Run (dry-run):

    ./branch-pruner -dry-run -keep=3

Flags:
  -dry-run    show branches that would be deleted (default true)
# branch-pruner

Simple CLI to help prune local and remote Git branches.

## Build

Build the CLI executable:

```bash
go build ./cmd/branch-pruner
```

## Quick run (dry-run)

Show branches that would be deleted (no deletions):

```bash
./branch-pruner -dry-run -keep=3
```

## Flags

- `-dry-run` (default: true): show branches that would be deleted but don't delete them.
- `-remote` (default: `origin`): remote name to prune branches from.
- `-keep` (default: 5): number of most recent branches to keep (based on commit timestamp).
- `-prefix` (default: none): only consider branches with this prefix.
- `-yes`: assume yes for all deletions (non-interactive).
- `-interactive`: prompt before deleting each branch.
- `-confirm`: prompt once to confirm all deletions before they run.

## How selection works

`branch-pruner` ranks branches by the unix timestamp of their latest commit (the most recent commit on that branch). It preserves the `keep` branches with the newest timestamps and selects older branches for deletion. This is deterministic and useful when branch names don't reflect recency.

Notes:

- If the tool cannot read a branch's latest commit timestamp (git error), that branch is treated as very old and will be selected for deletion. This conservative behavior can be changed on request.
- Remote deletion is attempted via `git push <remote> :<branch>` but remote errors are currently ignored by default; local deletion uses `git branch -D <branch>`.

## Examples

List branches that would be deleted, keeping the 3 most recently updated branches:

```bash
./branch-pruner -dry-run -keep=3
```

Only consider branches with prefix `feature/`:

```bash
./branch-pruner -dry-run -prefix=feature/
```

### Safety flags

- `-dry-run` (default true): show branches that would be deleted but do not delete them.
- `-yes`: assume yes to all deletions; useful for scripting non-interactive runs.
- `-interactive`: prompt before deleting each branch.
- `-confirm`: prompt once to confirm all deletions before they run.

Interactive per-branch confirmation:

```bash
./branch-pruner -dry-run=false -interactive
```

Non-interactive destructive run (use with caution):

```bash
./branch-pruner -dry-run=false -yes
```

Bulk confirm once for all deletions:

```bash
./branch-pruner -dry-run=false -confirm
```


```bash

