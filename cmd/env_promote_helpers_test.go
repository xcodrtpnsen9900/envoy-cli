package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyPromotionNoOverwrite(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "old"}

	promoted, skipped := applyPromotion(src, dst, nil, false)

	if dst["A"] != "old" {
		t.Errorf("A should not be overwritten, got %q", dst["A"])
	}
	if dst["B"] != "2" {
		t.Errorf("B should be promoted, got %q", dst["B"])
	}
	if len(promoted) != 1 || promoted[0] != "B" {
		t.Errorf("expected promoted=[B], got %v", promoted)
	}
	if len(skipped) != 1 || skipped[0] != "A" {
		t.Errorf("expected skipped=[A], got %v", skipped)
	}
}

func TestApplyPromotionWithOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	promoted, skipped := applyPromotion(src, dst, nil, true)

	if dst["A"] != "new" {
		t.Errorf("expected A=new, got %q", dst["A"])
	}
	if len(promoted) != 1 {
		t.Errorf("expected 1 promoted, got %d", len(promoted))
	}
	if len(skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(skipped))
	}
}

func TestApplyPromotionFilterKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	promoted, _ := applyPromotion(src, dst, []string{"A", "C"}, false)

	if _, ok := dst["B"]; ok {
		t.Error("B should not be promoted")
	}
	if len(promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(promoted))
	}
}

func TestWriteEnvMap(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.env")

	m := map[string]string{"Z": "last", "A": "first"}
	if err := writeEnvMap(p, m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)
	if content != "A=first\nZ=last\n" {
		t.Errorf("unexpected content:\n%s", content)
	}
}

func TestFilterKeys(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux", "HELLO=world"}
	out := filterKeys(lines, []string{"FOO", "HELLO"})
	if len(out) != 2 {
		t.Errorf("expected 2 lines, got %d", len(out))
	}
}
