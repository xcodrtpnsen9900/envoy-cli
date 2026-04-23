package cmd

import (
	"reflect"
	"testing"
)

func TestReorderLinesBasic(t *testing.T) {
	lines := []string{"APP=foo", "PORT=3000", "DEBUG=true"}
	result := reorderLines(lines, []string{"PORT", "APP"})
	if result[0] != "PORT=3000" {
		t.Errorf("expected PORT first, got %q", result[0])
	}
	if result[1] != "APP=foo" {
		t.Errorf("expected APP second, got %q", result[1])
	}
	if result[2] != "DEBUG=true" {
		t.Errorf("expected DEBUG third, got %q", result[2])
	}
}

func TestReorderLinesSkipsComments(t *testing.T) {
	lines := []string{"# header", "APP=foo", "PORT=3000"}
	result := reorderLines(lines, []string{"PORT"})
	if result[0] != "PORT=3000" {
		t.Errorf("expected PORT first, got %q", result[0])
	}
	found := false
	for _, l := range result {
		if l == "# header" {
			found = true
		}
	}
	if !found {
		t.Error("expected comment to be preserved")
	}
}

func TestReorderLinesMissingKeyIgnored(t *testing.T) {
	lines := []string{"APP=foo"}
	result := reorderLines(lines, []string{"MISSING"})
	if len(result) != 1 || result[0] != "APP=foo" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestKeysInOrder(t *testing.T) {
	lines := []string{"# comment", "APP=foo", "PORT=3000", "DEBUG=true"}
	keys := keysInOrder(lines)
	expected := []string{"APP", "PORT", "DEBUG"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestKeysInOrderEmpty(t *testing.T) {
	keys := keysInOrder([]string{"# only comments", ""})
	if len(keys) != 0 {
		t.Errorf("expected no keys, got %v", keys)
	}
}
