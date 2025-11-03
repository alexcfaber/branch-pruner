package git

import (
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "testing"
    "time"
)

// runGit runs git commands in dir and fails the test on error
func runGit(t *testing.T, dir string, env []string, args ...string) {
    t.Helper()
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    cmd.Env = append(os.Environ(), env...)
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("git %v failed: %v\noutput: %s", args, err, string(out))
    }
}

func TestIntegration_SelectBranchesByTimestamp(t *testing.T) {
    // create temp dir
    dir, err := os.MkdirTemp("", "branch-pruner-integ-")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(dir)

    // init repo
    runGit(t, dir, nil, "init")
    runGit(t, dir, nil, "config", "user.name", "Test User")
    runGit(t, dir, nil, "config", "user.email", "test@example.com")

    // create an initial file and commit at t0
    file := filepath.Join(dir, "file.txt")
    if err := os.WriteFile(file, []byte("initial"), 0o644); err != nil {
        t.Fatal(err)
    }
    runGit(t, dir, nil, "add", "file.txt")
    env := []string{"GIT_AUTHOR_DATE=1970-01-01T00:00:10Z", "GIT_COMMITTER_DATE=1970-01-01T00:00:10Z"}
    runGit(t, dir, env, "commit", "-m", "initial")

    // create branches with controlled commit timestamps
    branches := map[string]string{
        "old-1": "1970-01-01T00:00:20Z",
        "old-2": "1970-01-01T00:00:30Z",
        "new-1": "1970-01-01T00:01:00Z",
        "new-2": "1970-01-01T00:02:00Z",
    }

    for name, date := range branches {
        // create branch and commit a change with the given date
        runGit(t, dir, nil, "checkout", "-b", name)
        if err := os.WriteFile(file, []byte(name+"\n"), 0o644); err != nil {
            t.Fatal(err)
        }
        runGit(t, dir, nil, "add", "file.txt")
        env := []string{"GIT_AUTHOR_DATE=" + date, "GIT_COMMITTER_DATE=" + date}
        runGit(t, dir, env, "commit", "-m", "commit for "+name)
        // small sleep to ensure git internal ordering if needed
        time.Sleep(10 * time.Millisecond)
    }

    // list branches using the client
    c := NewClient("origin")
    // change working dir for commands by temporarily switching cwd
    origWd, _ := os.Getwd()
    defer os.Chdir(origWd)
    if err := os.Chdir(dir); err != nil {
        t.Fatal(err)
    }

    gotBranches, err := c.ListLocalBranches()
    if err != nil {
        t.Fatalf("ListLocalBranches failed: %v", err)
    }
    if len(gotBranches) == 0 {
        t.Fatalf("expected branches, got none")
    }

    // build timestamps map
    ts := map[string]int64{}
    for _, b := range gotBranches {
        tval, err := GetBranchCommitUnixTime(b)
        if err != nil {
            t.Fatalf("failed to get timestamp for %s: %v", b, err)
        }
        ts[b] = tval
    }

    // choose keep=2; compute expected: oldest len(branches)-2 branches
    keep := 2
    // compute expected by sorting
    type it struct{ name string; ts int64 }
    items := make([]it, 0, len(gotBranches))
    for _, b := range gotBranches {
        items = append(items, it{b, ts[b]})
    }
    sort.Slice(items, func(i, j int) bool {
        if items[i].ts == items[j].ts { return items[i].name < items[j].name }
        return items[i].ts < items[j].ts
    })
    numDel := len(items) - keep
    if numDel < 0 { numDel = 0 }
    expected := map[string]struct{}{}
    for i := 0; i < numDel; i++ {
        expected[items[i].name] = struct{}{}
    }

    gotDel := SelectBranchesToDelete(gotBranches, keep)
    // convert to map for comparison
    gotMap := map[string]struct{}{}
    for _, b := range gotDel { gotMap[b] = struct{}{} }

    if len(gotMap) != len(expected) {
        t.Fatalf("expected %d deletions, got %d; expected=%v got=%v", len(expected), len(gotMap), expected, gotMap)
    }
    for k := range expected {
        if _, ok := gotMap[k]; !ok {
            t.Fatalf("expected branch %s to be selected for deletion but it was not; expected=%v got=%v", k, expected, gotMap)
        }
    }
}
