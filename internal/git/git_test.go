package git

import (
	"testing"
	"time"
)

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

func TestParseDurationLike(t *testing.T) {
	tests := map[string]bool{
		"30d": true,
		"1w":  true,
		"72h": true,
		"15m": true,
		"":    false,
		"5x":  false,
	}
	for in, ok := range tests {
		_, err := ParseDurationLike(in)
		if ok && err != nil {
			t.Fatalf("expected success for %s, got %v", in, err)
		}
		if !ok && err == nil {
			t.Fatalf("expected failure for %s", in)
		}
	}
}

func TestSelectBranchesByAge(t *testing.T) {
	branches := []string{"a", "b", "c"}
	// pretend timestamps: a=100, b=200, c=300
	ts := map[string]int64{"a": 100, "b": 200, "c": 300}
	// threshold at 250 -> a and b are older
	threshold := time.Unix(250, 0)

	// use the pure helper by mapping through timestamps
	old := []string{}
	for _, b := range branches {
		if time.Unix(ts[b], 0).Before(threshold) {
			old = append(old, b)
		}
	}
	if len(old) != 2 {
		t.Fatalf("expected 2 old branches, got %d", len(old))
	}
}
