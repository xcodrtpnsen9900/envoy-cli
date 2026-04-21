package cmd

import (
	"testing"
)

func TestComputeWatchDiffAdded(t *testing.T) {
	old := "FOO=bar\n"
	new_ := "FOO=bar\nBAZ=qux\n"
	d := computeWatchDiff(old, new_)
	if len(d.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(d.Added))
	}
	if d.Added["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", d.Added["BAZ"])
	}
}

func TestComputeWatchDiffRemoved(t *testing.T) {
	old := "FOO=bar\nBAZ=qux\n"
	new_ := "FOO=bar\n"
	d := computeWatchDiff(old, new_)
	if len(d.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(d.Removed))
	}
	if d.Removed["BAZ"] != "qux" {
		t.Errorf("expected removed BAZ=qux, got %s", d.Removed["BAZ"])
	}
}

func TestComputeWatchDiffChanged(t *testing.T) {
	old := "FOO=bar\n"
	new_ := "FOO=baz\n"
	d := computeWatchDiff(old, new_)
	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(d.Changed))
	}
	pair := d.Changed["FOO"]
	if pair[0] != "bar" || pair[1] != "baz" {
		t.Errorf("expected bar->baz, got %v", pair)
	}
}

func TestComputeWatchDiffIdentical(t *testing.T) {
	content := "FOO=bar\nBAZ=qux\n"
	d := computeWatchDiff(content, content)
	if hasWatchChanges(d) {
		t.Error("expected no changes for identical content")
	}
}

func TestParseEnvContentSkipsComments(t *testing.T) {
	content := "# comment\nFOO=bar\n"
	m := parseEnvContent(content)
	if _, ok := m["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", m["FOO"])
	}
}

func TestParseEnvContentSkipsMalformed(t *testing.T) {
	content := "NOEQUALS\nFOO=bar\n"
	m := parseEnvContent(content)
	if _, ok := m["NOEQUALS"]; ok {
		t.Error("malformed line should be skipped")
	}
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", m["FOO"])
	}
}

func TestHasWatchChanges(t *testing.T) {
	empty := watchDiffResult{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
	}
	if hasWatchChanges(empty) {
		t.Error("expected no changes for empty diff")
	}
	empty.Added["X"] = "1"
	if !hasWatchChanges(empty) {
		t.Error("expected changes when Added is non-empty")
	}
}
