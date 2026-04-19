package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupAuditDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestWriteAndReadAuditLog(t *testing.T) {
	dir := setupAuditDir(t)
	writeAuditEntry(dir, "switch", "production")
	writeAuditEntry(dir, "switch", "staging")

	entries, err := readAuditLog(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !strings.Contains(entries[0], "switch") || !strings.Contains(entries[0], "production") {
		t.Errorf("unexpected entry: %s", entries[0])
	}
}

func TestReadAuditLogEmpty(t *testing.T) {
	dir := setupAuditDir(t)
	entries, err := readAuditLog(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestReadAuditLogNoDir(t *testing.T) {
	dir := t.TempDir()
	entries, err := readAuditLog(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries, got %v", entries)
	}
}

func TestAuditEntryFormat(t *testing.T) {
	dir := setupAuditDir(t)
	writeAuditEntry(dir, "delete", "dev")

	data, err := os.ReadFile(auditLogPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	line := string(data)
	if !strings.Contains(line, "[2") {
		t.Errorf("expected timestamp, got: %s", line)
	}
	if !strings.Contains(line, "profile=dev") {
		t.Errorf("expected profile=dev, got: %s", line)
	}
}
