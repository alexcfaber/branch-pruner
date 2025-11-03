package git

import (
    "bufio"
    "errors"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type Client struct{
    remote string
}

func NewClient(remote string) *Client {
    return &Client{remote: remote}
}

// ListLocalBranches returns a slice of branch names (excluding HEAD/remote refs)
func (c *Client) ListLocalBranches() ([]string, error) {
    cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
    out, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(strings.NewReader(string(out)))
    branches := []string{}
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        branches = append(branches, line)
    }
    return branches, nil
}

func (c *Client) DeleteBranch(branch string) error {
    // delete local branch forcefully
    cmd := exec.Command("git", "branch", "-D", branch)
    if err := cmd.Run(); err != nil {
        return err
    }
    // attempt to delete remote branch
    remoteRef := c.remote + "/" + branch
    cmd = exec.Command("git", "push", c.remote, ":"+branch)
    _ = cmd.Run() // ignore remote errors for now
    _ = remoteRef
    return nil
}

// SelectBranchesToDelete picks branches to delete keeping 'keep' most recent by commit date.
// It queries the latest commit timestamp for each branch and keeps the `keep` branches with
// the newest timestamps. Older branches are returned for deletion.
func SelectBranchesToDelete(branches []string, keep int) []string {
    if keep < 0 {
        keep = 0
    }
    if len(branches) <= keep {
        return nil
    }

    // gather timestamps; if we can't get a timestamp for a branch treat it as very old (0)
    ts := map[string]int64{}
    for _, b := range branches {
        t, err := GetBranchCommitUnixTime(b)
        if err != nil {
            ts[b] = 0
            continue
        }
        ts[b] = t
    }
    return SelectBranchesToDeleteByTimestamps(branches, ts, keep)
}

// SelectBranchesToDeleteByTimestamps is a pure function that, given a map of branch->timestamp,
// returns the branches that should be deleted (older ones), keeping `keep` most recent.
func SelectBranchesToDeleteByTimestamps(branches []string, timestamps map[string]int64, keep int) []string {
    if keep < 0 {
        keep = 0
    }
    if len(branches) <= keep {
        return nil
    }

    type item struct{
        name string
        ts   int64
    }
    items := make([]item, 0, len(branches))
    for _, b := range branches {
        t := int64(0)
        if v, ok := timestamps[b]; ok {
            t = v
        }
        items = append(items, item{name: b, ts: t})
    }

    // sort by timestamp ascending (oldest first)
    sort.Slice(items, func(i, j int) bool {
        if items[i].ts == items[j].ts {
            // tie-breaker: alphabetical
            return items[i].name < items[j].name
        }
        return items[i].ts < items[j].ts
    })

    // number to delete
    del := len(items) - keep
    if del <= 0 {
        return nil
    }
    out := make([]string, 0, del)
    for i := 0; i < del; i++ {
        out = append(out, items[i].name)
    }
    return out
}

// GetBranchCommitUnixTime returns the unix timestamp (seconds since epoch) of the latest commit on the branch
func GetBranchCommitUnixTime(branch string) (int64, error) {
    cmd := exec.Command("git", "log", "-1", "--format=%ct", branch)
    out, err := cmd.Output()
    if err != nil {
        return 0, err
    }
    s := strings.TrimSpace(string(out))
    if s == "" {
        return 0, errors.New("no commit timestamp")
    }
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0, err
    }
    return v, nil
}

// ParseBranchesFromString splits git output into branches (helper used by tests)
func ParseBranchesFromString(s string) ([]string, error) {
    if s == "" {
        return nil, errors.New("empty input")
    }
    lines := []string{}
    scanner := bufio.NewScanner(strings.NewReader(s))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" { continue }
        lines = append(lines, line)
    }
    return lines, nil
}
