package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupAnnotateDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy", "profiles"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestSetAndGetAnnotation(t *testing.T) {
	dir := setupAnnotateDir(t)
	if err := setAnnotation(dir, "dev", "DB_HOST", "Primary database hostname"); err != nil {
		t.Fatalf("setAnnotation: %v", err)
	}
	note, err := getAnnotationValue(dir, "dev", "DB_HOST")
	if err != nil {
		t.Fatalf("getAnnotationValue: %v", err)
	}
	if note != "Primary database hostname" {
		t.Errorf("expected annotation %q, got %q", "Primary database hostname", note)
	}
}

func TestGetAnnotationMissingKey(t *testing.T) {
	dir := setupAnnotateDir(t)
	note, err := getAnnotationValue(dir, "dev", "MISSING_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note != "" {
		t.Errorf("expected empty annotation, got %q", note)
	}
}

func TestLoadAnnotationsEmpty(t *testing.T) {
	dir := setupAnnotateDir(t)
	annotations, err := loadAnnotations(dir, "staging")
	if err != nil {
		t.Fatalf("loadAnnotations: %v", err)
	}
	if len(annotations) != 0 {
		t.Errorf("expected empty map, got %d entries", len(annotations))
	}
}

func TestAnnotationCount(t *testing.T) {
	dir := setupAnnotateDir(t)
	_ = setAnnotation(dir, "prod", "API_KEY", "External API key")
	_ = setAnnotation(dir, "prod", "DB_PASS", "Database password")
	count, err := annotationCount(dir, "prod")
	if err != nil {
		t.Fatalf("annotationCount: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 annotations, got %d", count)
	}
}

func TestRemoveAnnotation(t *testing.T) {
	dir := setupAnnotateDir(t)
	_ = setAnnotation(dir, "dev", "PORT", "App port")
	_ = setAnnotation(dir, "dev", "HOST", "App host")
	if err := removeAnnotation(dir, "dev", "PORT"); err != nil {
		t.Fatalf("removeAnnotation: %v", err)
	}
	exists, err := annotationExists(dir, "dev", "PORT")
	if err != nil {
		t.Fatalf("annotationExists: %v", err)
	}
	if exists {
		t.Error("expected annotation to be removed")
	}
	exists, err = annotationExists(dir, "dev", "HOST")
	if err != nil {
		t.Fatalf("annotationExists: %v", err)
	}
	if !exists {
		t.Error("expected HOST annotation to remain")
	}
}

func TestSortedAnnotationKeys(t *testing.T) {
	annotations := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	keys := sortedAnnotationKeys(annotations)
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
