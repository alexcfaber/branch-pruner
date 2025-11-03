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
  -keep       number of most recent branches to keep per author (default 5)
  -prefix     only consider branches with this prefix
