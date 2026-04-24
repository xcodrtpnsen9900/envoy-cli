package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupPlaceholderDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writePlaceholderProfile(t *testing.T, root, name, content string) {
	t.Helper()
	p := profilePath(root, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestFindPlaceholdersEmpty(t *testing.T) {
	dir := setupPlaceholderDir(t)
	writePlaceholderProfile(t, dir, "dev", "DB_HOST=localhost\nDB_PORT=5432\n")
	results, err := findPlaceholders(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFindPlaceholdersDetectsChangeme(t *testing.T) {
	dir := setupPlaceholderDir(t)
	writePlaceholderProfile(t, dir, "dev", "API_KEY=CHANGEME\nDB_HOST=localhost\n")
	results, err := findPlaceholders(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY placeholder, got %+v", results)
	}
}

func TestFindPlaceholdersDetectsEmptyValue(t *testing.T) {
	dir := setupPlaceholderDir(t)
	writePlaceholderProfile(t, dir, "dev", "SECRET=\nHOST=localhost\n")
	results, err := findPlaceholders(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "SECRET" {
		t.Errorf("expected SECRET placeholder, got %+v", results)
	}
}

func TestFindPlaceholdersSkipsComments(t *testing.T) {
	dir := setupPlaceholderDir(t)
	writePlaceholderProfile(t, dir, "dev", "# TODO: fill this in\nAPI_KEY=real-value\n")
	results, err := findPlaceholders(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFindPlaceholdersNonExistentProfile(t *testing.T) {
	dir := setupPlaceholderDir(t)
	_, err := findPlaceholders(dir, "ghost")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestWritePlaceholderReport(t *testing.T) {
	dir := setupPlaceholderDir(t)
	results := []placeholderResult{
		{Line: 1, Key: "API_KEY", Value: "CHANGEME"},
		{Line: 3, Key: "SECRET", Value: ""},
	}
	if err := writePlaceholderReport(dir, "report", results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := findPlaceholders(dir, "report")
	if err != nil {
		t.Fatalf("unexpected error reading report: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 placeholder entries in report, got %d", len(out))
	}
}
