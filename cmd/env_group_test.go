package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupGroupDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestAddProfilesToGroup(t *testing.T) {
	dir := setupGroupDir(t)
	if err := addProfilesToGroup(dir, "staging", []string{"dev", "qa"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profiles, err := profilesInGroup(dir, "staging")
	if err != nil {
		t.Fatalf("profilesInGroup: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestAddProfilesToGroupDeduplicates(t *testing.T) {
	dir := setupGroupDir(t)
	_ = addProfilesToGroup(dir, "g1", []string{"alpha", "beta"})
	_ = addProfilesToGroup(dir, "g1", []string{"alpha", "gamma"})
	profiles, _ := profilesInGroup(dir, "g1")
	if len(profiles) != 3 {
		t.Fatalf("expected 3 unique profiles, got %d", len(profiles))
	}
}

func TestDeleteGroup(t *testing.T) {
	dir := setupGroupDir(t)
	_ = addProfilesToGroup(dir, "mygroup", []string{"dev"})
	if err := deleteGroup(dir, "mygroup"); err != nil {
		t.Fatalf("deleteGroup: %v", err)
	}
	exists, _ := groupExists(dir, "mygroup")
	if exists {
		t.Fatal("expected group to be deleted")
	}
}

func TestDeleteNonExistentGroup(t *testing.T) {
	dir := setupGroupDir(t)
	if err := deleteGroup(dir, "nope"); err == nil {
		t.Fatal("expected error deleting non-existent group")
	}
}

func TestGroupNamesReturnsAllGroups(t *testing.T) {
	dir := setupGroupDir(t)
	_ = addProfilesToGroup(dir, "z-group", []string{"x"})
	_ = addProfilesToGroup(dir, "a-group", []string{"y"})
	names, err := groupNames(dir)
	if err != nil {
		t.Fatalf("groupNames: %v", err)
	}
	if len(names) != 2 || names[0] != "a-group" {
		t.Fatalf("expected sorted group names, got %v", names)
	}
}

func TestLoadGroupsEmptyDir(t *testing.T) {
	dir := setupGroupDir(t)
	groups, err := loadGroups(dir)
	if err != nil {
		t.Fatalf("loadGroups: %v", err)
	}
	if len(groups) != 0 {
		t.Fatalf("expected empty groups, got %d", len(groups))
	}
}
