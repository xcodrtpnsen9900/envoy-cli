package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupHealthDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeHealthProfile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, ".envoy", name+".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheckHealthy(t *testing.T) {
	root := setupHealthDir(t)
	writeHealthProfile(t, root, "prod", "HOST=localhost\nPORT=8080\n")

	r, err := healthCheck(root, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.OK {
		t.Errorf("expected healthy, got issues: %v", healthSummary(r))
	}
	if r.TotalKeys != 2 {
		t.Errorf("expected 2 keys, got %d", r.TotalKeys)
	}
}

func TestHealthCheckEmptyValue(t *testing.T) {
	root := setupHealthDir(t)
	writeHealthProfile(t, root, "dev", "HOST=\nPORT=9090\n")

	r, err := healthCheck(root, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.OK {
		t.Error("expected unhealthy due to empty value")
	}
	if len(r.EmptyVals) != 1 || r.EmptyVals[0] != "HOST" {
		t.Errorf("expected EmptyVals=[HOST], got %v", r.EmptyVals)
	}
}

func TestHealthCheckDuplicateKey(t *testing.T) {
	root := setupHealthDir(t)
	writeHealthProfile(t, root, "staging", "HOST=a\nHOST=b\nPORT=80\n")

	r, err := healthCheck(root, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.OK {
		t.Error("expected unhealthy due to duplicate key")
	}
	if len(r.Duplicates) != 1 || r.Duplicates[0] != "HOST" {
		t.Errorf("expected Duplicates=[HOST], got %v", r.Duplicates)
	}
}

func TestHealthCheckNonExistentProfile(t *testing.T) {
	root := setupHealthDir(t)
	_, err := healthCheck(root, "ghost")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestHealthCheckSkipsComments(t *testing.T) {
	root := setupHealthDir(t)
	writeHealthProfile(t, root, "ci", "# comment\nKEY=value\n")

	r, err := healthCheck(root, "ci")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.TotalKeys != 1 {
		t.Errorf("expected 1 key (comment skipped), got %d", r.TotalKeys)
	}
	if !r.OK {
		t.Error("expected healthy")
	}
}

func TestHealthIssueCount(t *testing.T) {
	r := &healthReport{
		EmptyVals:  []string{"A", "B"},
		Duplicates: []string{"C"},
	}
	if n := healthIssueCount(r); n != 3 {
		t.Errorf("expected 3 issues, got %d", n)
	}
}

func TestHealthSummaryHealthy(t *testing.T) {
	r := &healthReport{OK: true}
	if s := healthSummary(r); s != "healthy" {
		t.Errorf("expected 'healthy', got %q", s)
	}
}
