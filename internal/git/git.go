package git

import (
    "bufio"
    "errors"
    "os/exec"
    "sort"
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

// SelectBranchesToDelete picks branches to delete keeping 'keep' most recent by simple alphabetical order
// This is a placeholder heuristic; real implementation would sort by commit date or author.
func SelectBranchesToDelete(branches []string, keep int) []string {
    if keep < 0 {
        keep = 0
    }
    if len(branches) <= keep {
        return nil
    }

    // simple stable sort
    sort.Strings(branches)

    // keep last `keep` branches
    keepStart := len(branches) - keep
    if keepStart < 0 {
        keepStart = 0
    }
    return append([]string{}, branches[:keepStart]...)
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
