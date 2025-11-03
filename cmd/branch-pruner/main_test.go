package main

import (
    "os"
    "path/filepath"
    "testing"
)

func TestMergeWithConfig(t *testing.T) {
    tmpdir, err := os.MkdirTemp("", "bpconfig")
    if err != nil {
        t.Fatal(err)
    }
    defer func() {
        if err := os.RemoveAll(tmpdir); err != nil {
            t.Logf("warning: failed to remove temp dir %s: %v", tmpdir, err)
        }
    }()

    cfg := `keep: 2
prefix: feature/
remote: origin
older-than: 7d
`
    cfgPath := filepath.Join(tmpdir, ".branch-pruner.yaml")
    if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
        t.Fatal(err)
    }

    p := ""
    r := "origin"
    k := 5
    o := ""

    if err := mergeWithConfig(&p, &r, &k, &o, cfgPath); err != nil {
        t.Fatalf("merge failed: %v", err)
    }

    if p != "feature/" {
        t.Fatalf("expected prefix from config, got %s", p)
    }
    if k != 2 {
        t.Fatalf("expected keep=2 from config, got %d", k)
    }
    if o != "7d" {
        t.Fatalf("expected older-than from config, got %s", o)
    }
}
