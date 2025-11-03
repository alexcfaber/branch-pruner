package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/alexcfaber/branch-pruner/internal/git"
)

func main() {
    dryRun := flag.Bool("dry-run", true, "Show branches that would be deleted but don't delete them")
    remote := flag.String("remote", "origin", "Remote name to prune branches from")
    keep := flag.Int("keep", 5, "Number of most recent branches to keep (based on commit timestamp)")
    prefix := flag.String("prefix", "", "Only consider branches with this prefix")
    yes := flag.Bool("yes", false, "Assume yes for all deletions (non-interactive)")
    interactive := flag.Bool("interactive", false, "Prompt before deleting each branch")
    flag.Parse()

    fmt.Println("branch-pruner — prune old git branches by commit date")
    fmt.Println("Selection heuristic: branches are ranked by the timestamp of their latest commit; the newest 'keep' branches are preserved and older branches are selected for deletion.")
    fmt.Println()

    g := git.NewClient(*remote)

    branches, err := g.ListLocalBranches()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error listing branches: %v\n", err)
        os.Exit(1)
    }

    // filter by prefix if given
    if *prefix != "" {
        filtered := []string{}
        for _, b := range branches {
            if strings.HasPrefix(b, *prefix) {
                filtered = append(filtered, b)
            }
        }
        branches = filtered
    }

    toDelete := git.SelectBranchesToDelete(branches, *keep)

    if len(toDelete) == 0 {
        fmt.Println("No branches to delete")
        return
    }

    fmt.Printf("Branches to delete (%d):\n", len(toDelete))
    for _, b := range toDelete {
        fmt.Println(" - ", b)
    }

    if *dryRun {
        fmt.Println("Dry run: no branches were deleted")
        return
    }

    for _, b := range toDelete {
        // decide whether to delete based on flags
        doDelete := true
        if *yes {
            doDelete = true
        } else if *interactive {
            // prompt
            fmt.Printf("Delete branch %s? [y/N]: ", b)
            var resp string
            if _, err := fmt.Scanln(&resp); err != nil {
                // treat as no
                doDelete = false
            } else {
                resp = strings.TrimSpace(strings.ToLower(resp))
                doDelete = (resp == "y" || resp == "yes")
            }
        } else if *dryRun {
            doDelete = false
        }

        if !doDelete {
            fmt.Printf("skipped %s\n", b)
            continue
        }

        if err := g.DeleteBranch(b); err != nil {
            fmt.Fprintf(os.Stderr, "failed to delete %s: %v\n", b, err)
        } else {
            fmt.Printf("deleted %s\n", b)
        }
    }
}
