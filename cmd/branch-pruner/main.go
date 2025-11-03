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
    keep := flag.Int("keep", 5, "Number of most recent branches to keep per author")
    prefix := flag.String("prefix", "", "Only consider branches with this prefix")
    flag.Parse()

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
        if err := g.DeleteBranch(b); err != nil {
            fmt.Fprintf(os.Stderr, "failed to delete %s: %v\n", b, err)
        } else {
            fmt.Printf("deleted %s\n", b)
        }
    }
}
