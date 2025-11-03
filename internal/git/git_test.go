package git

import "testing"

func TestSelectBranchesToDelete(t *testing.T) {
    branches := []string{"a", "b", "c", "d", "e"}
    // provide timestamps so that 'e' and 'd' are newest
    ts := map[string]int64{
        "a": 10,
        "b": 20,
        "c": 30,
        "d": 40,
        "e": 50,
    }
    toDel := SelectBranchesToDeleteByTimestamps(branches, ts, 2)
    if len(toDel) != 3 {
        t.Fatalf("expected 3, got %d", len(toDel))
    }
    // oldest three should be a,b,c in that order
    if toDel[0] != "a" || toDel[1] != "b" || toDel[2] != "c" {
        t.Fatalf("unexpected deletion list: %v", toDel)
    }
}

func TestParseBranchesFromString(t *testing.T) {
    input := "\nfeature/one\nfeature/two\n\n" 
    got, err := ParseBranchesFromString(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(got) != 2 {
        t.Fatalf("expected 2 branches, got %d", len(got))
    }
}
