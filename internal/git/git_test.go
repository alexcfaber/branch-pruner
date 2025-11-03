package git

import "testing"

func TestSelectBranchesToDelete(t *testing.T) {
    branches := []string{"a", "b", "c", "d", "e"}
    toDel := SelectBranchesToDelete(branches, 2)
    if len(toDel) != 3 {
        t.Fatalf("expected 3, got %d", len(toDel))
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
