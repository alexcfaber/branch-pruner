package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexcfaber/branch-pruner/internal/config"
	"github.com/alexcfaber/branch-pruner/internal/git"
)

func main() {
	dryRun := flag.Bool("dry-run", true, "Show branches that would be deleted but don't delete them")
	remote := flag.String("remote", "origin", "Remote name to prune branches from")
	keep := flag.Int("keep", 5, "Number of most recent branches to keep (based on commit timestamp)")
	prefix := flag.String("prefix", "", "Only consider branches with this prefix")
	olderThan := flag.String("older-than", "", "Only consider branches with latest commit older than this duration (e.g. 30d, 72h)")
	configPath := flag.String("config", "", "Path to YAML config file (defaults: ~/.branch-pruner.yaml, ./.branch-pruner.yaml)")
	yes := flag.Bool("yes", false, "Assume yes for all deletions (non-interactive)")
	interactive := flag.Bool("interactive", false, "Prompt before deleting each branch")
	confirm := flag.Bool("confirm", false, "Prompt once to confirm all deletions")
	flag.Parse()

	if err := mergeWithConfig(prefix, remote, keep, olderThan, *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

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

	// filter by older-than if provided
	if *olderThan != "" {
		d, err := git.ParseDurationLike(*olderThan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid older-than: %v\n", err)
			os.Exit(1)
		}
		threshold := time.Now().Add(-d)
		older, err := git.FilterBranchesOlderThan(branches, threshold)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error filtering by age: %v\n", err)
			os.Exit(1)
		}
		// use only older branches from now on
		branches = older
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

	// bulk confirmation (single prompt)
	bulkConfirmed := false
	if *yes {
		bulkConfirmed = true
	} else if *confirm {
		fmt.Printf("About to delete %d branches. Continue? [y/N]: ", len(toDelete))
		var resp string
		if _, err := fmt.Scanln(&resp); err != nil {
			bulkConfirmed = false
		} else {
			resp = strings.TrimSpace(strings.ToLower(resp))
			bulkConfirmed = (resp == "y" || resp == "yes")
		}
	}

	for _, b := range toDelete {
		// decide whether to delete based on flags
		doDelete := true

		if *dryRun {
			doDelete = false
		}

		if *yes {
			doDelete = true
		} else if *confirm {
			// if confirm was requested but not accepted, skip
			if !bulkConfirmed {
				doDelete = false
			}
		} else if *interactive {
			// per-branch prompt
			fmt.Printf("Delete branch %s? [y/N]: ", b)
			var resp string
			if _, err := fmt.Scanln(&resp); err != nil {
				doDelete = false
			} else {
				resp = strings.TrimSpace(strings.ToLower(resp))
				doDelete = (resp == "y" || resp == "yes")
			}
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

// mergeWithConfig loads config (if present) and merges into provided flag pointers
// It prefers existing flag values unless they are empty/defaults.
func mergeWithConfig(prefix *string, remote *string, keep *int, olderThan *string, configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}
	if *prefix == "" && cfg.Prefix != "" {
		*prefix = cfg.Prefix
	}
	if *remote == "origin" && cfg.Remote != "" {
		*remote = cfg.Remote
	}
	if *keep == 5 && cfg.Keep != 0 {
		*keep = cfg.Keep
	}
	if *olderThan == "" && cfg.OlderThan != "" {
		*olderThan = cfg.OlderThan
	}
	return nil
}
