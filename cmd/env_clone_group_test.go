package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupCloneGroupDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	envoyDir := filepath.Join(root, ".envoy")
	if err := os.MkdirAll(envoyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeCloneGroupProfile(t *testing.T, root, name, content string) {
	t.Helper()
	p := filepath.Join(root, ".envoy", name+".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCloneGroupBasic(t *testing.T) {
	root := setupCloneGroupDir(t)
	writeCloneGroupProfile(t, root, "dev", "FOO=bar\n")
	writeCloneGroupProfile(t, root, "staging", "FOO=baz\n")

	_ = saveGroups(root, map[string][]string{"mygroup": {"dev", "staging"}})

	cloned, err := cloneGroup(root, "mygroup", "-backup", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cloned) != 2 {
		t.Fatalf("expected 2 cloned profiles, got %d", len(cloned))
	}

	for _, name := range []string{"dev-backup", "staging-backup"} {
		p := filepath.Join(root, ".envoy", name+".env")
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected cloned profile %q to exist", name)
		}
	}
}

func TestCloneGroupNonExistentGroup(t *testing.T) {
	root := setupCloneGroupDir(t)
	_ = saveGroups(root, map[string][]string{})

	_, err := cloneGroup(root, "ghost", "-copy", false)
	if err == nil {
		t.Fatal("expected error for non-existent group")
	}
}

func TestCloneGroupAlreadyExistsNoOverwrite(t *testing.T) {
	root := setupCloneGroupDir(t)
	writeCloneGroupProfile(t, root, "dev", "KEY=val\n")
	writeCloneGroupProfile(t, root, "dev-copy", "KEY=old\n")

	_ = saveGroups(root, map[string][]string{"g": {"dev"}})

	_, err := cloneGroup(root, "g", "-copy", false)
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestCloneGroupOverwrite(t *testing.T) {
	root := setupCloneGroupDir(t)
	writeCloneGroupProfile(t, root, "dev", "KEY=new\n")
	writeCloneGroupProfile(t, root, "dev-copy", "KEY=old\n")

	_ = saveGroups(root, map[string][]string{"g": {"dev"}})

	cloned, err := cloneGroup(root, "g", "-copy", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cloned) != 1 || cloned[0] != "dev-copy" {
		t.Fatalf("unexpected cloned list: %v", cloned)
	}

	data, _ := os.ReadFile(filepath.Join(root, ".envoy", "dev-copy.env"))
	if string(data) != "KEY=new\n" {
		t.Errorf("overwritten file has unexpected content: %q", string(data))
	}
}

func TestCloneGroupMissingSourceProfile(t *testing.T) {
	root := setupCloneGroupDir(t)
	// group references a profile that doesn't exist on disk
	_ = saveGroups(root, map[string][]string{"g": {"missing"}})

	_, err := cloneGroup(root, "g", "-copy", false)
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}
