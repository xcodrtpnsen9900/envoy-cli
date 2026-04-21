package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupBackupDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("ENVOY_PROJECT_DIR", tmp)
	envoyDir := filepath.Join(tmp, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return tmp
}

func writeBackupProfile(t *testing.T, tmp, name, content string) {
	t.Helper()
	p := filepath.Join(tmp, ".envoy", name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestBackupProfile(t *testing.T) {
	tmp := setupBackupDir(t)
	writeBackupProfile(t, tmp, "staging", "KEY=value\n")

	if err := backupProfile("staging", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if backupCount("staging") != 1 {
		t.Errorf("expected 1 backup, got %d", backupCount("staging"))
	}
}

func TestBackupNonExistentProfile(t *testing.T) {
	setupBackupDir(t)
	err := backupProfile("ghost", "")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestRestoreBackup(t *testing.T) {
	tmp := setupBackupDir(t)
	writeBackupProfile(t, tmp, "prod", "ENV=production\n")

	if err := backupProfile("prod", ""); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Overwrite profile
	writeBackupProfile(t, tmp, "prod", "ENV=changed\n")

	latest := latestBackup("prod")
	if latest == "" {
		t.Fatal("expected a backup to exist")
	}

	if err := restoreBackup("prod", latest); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	data, err := os.ReadFile(profilePath("prod"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "ENV=production\n" {
		t.Errorf("expected restored content, got %q", string(data))
	}
}

func TestPurgeOldBackups(t *testing.T) {
	tmp := setupBackupDir(t)
	writeBackupProfile(t, tmp, "dev", "X=1\n")

	for i := 0; i < 4; i++ {
		if err := backupProfile("dev", ""); err != nil {
			t.Fatalf("backup %d failed: %v", i, err)
		}
	}

	if err := purgeOldBackups("dev", 2); err != nil {
		t.Fatalf("purge failed: %v", err)
	}

	if backupCount("dev") != 2 {
		t.Errorf("expected 2 backups after purge, got %d", backupCount("dev"))
	}
}

func TestLatestBackupNoBackups(t *testing.T) {
	setupBackupDir(t)
	result := latestBackup("nonexistent")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
