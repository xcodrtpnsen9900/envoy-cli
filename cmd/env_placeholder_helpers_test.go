package cmd

import (
	"strings"
	"testing"
)

func TestPlaceholderSummaryEmpty(t *testing.T) {
	summary := placeholderSummary(nil)
	if !strings.Contains(summary, "No placeholders") {
		t.Errorf("expected 'No placeholders' in summary, got %q", summary)
	}
}

func TestPlaceholderSummaryWithResults(t *testing.T) {
	results := []placeholderResult{
		{Line: 2, Key: "TOKEN", Value: "CHANGEME"},
	}
	summary := placeholderSummary(results)
	if !strings.Contains(summary, "TOKEN") {
		t.Errorf("expected TOKEN in summary, got %q", summary)
	}
	if !strings.Contains(summary, "1 placeholder") {
		t.Errorf("expected count in summary, got %q", summary)
	}
}

func TestPlaceholderKeys(t *testing.T) {
	results := []placeholderResult{
		{Line: 1, Key: "A", Value: ""},
		{Line: 2, Key: "B", Value: "FIXME"},
	}
	keys := placeholderKeys(results)
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFilterPlaceholdersByToken(t *testing.T) {
	results := []placeholderResult{
		{Line: 1, Key: "A", Value: "CHANGEME"},
		{Line: 2, Key: "B", Value: "FIXME"},
		{Line: 3, Key: "C", Value: ""},
	}
	filtered := filterPlaceholdersByToken(results, "changeme")
	if len(filtered) != 1 || filtered[0].Key != "A" {
		t.Errorf("expected only A, got %+v", filtered)
	}
}

func TestCountEmptyPlaceholders(t *testing.T) {
	results := []placeholderResult{
		{Line: 1, Key: "A", Value: ""},
		{Line: 2, Key: "B", Value: "CHANGEME"},
		{Line: 3, Key: "C", Value: ""},
	}
	count := countEmptyPlaceholders(results)
	if count != 2 {
		t.Errorf("expected 2 empty placeholders, got %d", count)
	}
}
