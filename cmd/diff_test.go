package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffIdenticalProfiles(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	addProfile("base", ".env")
	p := profilePath("base")
	os.WriteFile(p, []byte("KEY=value\nFOO=bar\n"), 0644)

	addProfile("copy", ".env")
	os.WriteFile(profilePath("copy"), []byte("KEY=value\nFOO=bar\n"), 0644)

	// Should not panic; identical profiles produce no diff output
	valsA, _, err := parseEnvFile(profilePath("base"))
	if err != nil {
		t.Fatal(err)
	}
	valsB, _, err := parseEnvFile(profilePath("copy"))
	if err != nil {
		t.Fatal(err)
	}
	if valsA["KEY"] != valsB["KEY"] {
		t.Error("expected KEY to match")
	}
}

func TestDiffDifferentValues(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	addProfile("dev", ".env")
	os.WriteFile(profilePath("dev"), []byte("KEY=dev_val\nONLY_DEV=1\n"), 0644)

	addProfile("prod", ".env")
	os.WriteFile(profilePath("prod"), []byte("KEY=prod_val\nONLY_PROD=1\n"), 0644)

	valsA, keysA, _ := parseEnvFile(profilePath("dev"))
	valsB, _, _ := parseEnvFile(profilePath("prod"))

	changed := false
	for _, k := range keysA {
		if valsA[k] != valsB[k] {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("expected differences between dev and prod")
	}
}

func TestParseEnvFileSkipsComments(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.env")
	os.WriteFile(f, []byte("# comment\nKEY=val\n\nFOO=bar\n"), 0644)

	vals, keys, err := parseEnvFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	if vals["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", vals["KEY"])
	}
}

func TestParseEnvFileMissing(t *testing.T) {
	_, _, err := parseEnvFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
