package cmd

import (
	"testing"
)

func TestExtractKey(t *testing.T) {
	cases := []struct {
		line string
		want string
	}{
		{"KEY=value", "KEY"},
		{"MY_VAR=hello world", "MY_VAR"},
		{"NO_EQUALS", ""},
		{" SPACED = val", "SPACED"},
	}
	for _, tc := range cases {
		got := extractKey(tc.line)
		if got != tc.want {
			t.Errorf("extractKey(%q) = %q, want %q", tc.line, got, tc.want)
		}
	}
}

func TestIsSorted(t *testing.T) {
	if !isSorted([]string{"ALPHA=1", "BETA=2", "GAMMA=3"}) {
		t.Error("expected sorted slice to be detected as sorted")
	}
	if isSorted([]string{"ZEBRA=1", "ALPHA=2"}) {
		t.Error("expected unsorted slice to be detected as unsorted")
	}
	if !isSorted([]string{}) {
		t.Error("empty slice should be considered sorted")
	}
	if !isSorted([]string{"ONLY=1"}) {
		t.Error("single element should be considered sorted")
	}
}

func TestDeduplicateKeys(t *testing.T) {
	lines := []string{"KEY=first", "OTHER=val", "KEY=second"}
	result := deduplicateKeys(lines)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines after dedup, got %d", len(result))
	}
0]) != "OTHER" {
		t.Errorf("expected OTHER, got %s", result[0])
	}
	if result[1] != "KEY=second" {
		t.Errorf("expected last KEY value, got %s", result[1])
	}
}

func TestDeduplicateKeysNoDuplicates(t *testing.T) {
	lines := []string{"A=1", "B=2", "C=3"}
	result := deduplicateKeys(lines)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
}
