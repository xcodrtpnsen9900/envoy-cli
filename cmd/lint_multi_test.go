package cmd

import (
	"strings"
	"testing"
)

func TestLintMultipleIssues(t *testing.T) {
	dir := setupTempDir(t)
	content := "MY KEY= bad \nNOEQUALS\nOK=fine\n"
	writeProfile(t, dir, "multi", content)
	issues, err := lintProfile(dir, "multi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// MY KEY has whitespace in key + leading/trailing whitespace in value = 2 issues
	// NOEQUALS has missing '=' = 1 issue
	if len(issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d: %v", len(issues), issues)
	}
}

func TestLintTabInValue(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "tabval", "KEY=val\there\n")
	issues, err := lintProfile(dir, "tabval")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issue for tab in value")
	}
	if !strings.Contains(issues[0], "tab") {
		t.Fatalf("expected tab mention in issue, got: %s", issues[0])
	}
}

func TestLintEmptyKey(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "emptykey", "=value\n")
	issues, err := lintProfile(dir, "emptykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issue for empty key")
	}
}
