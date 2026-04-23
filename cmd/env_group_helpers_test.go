package cmd

import (
	"testing"
)

func TestGroupExists(t *testing.T) {
	dir := setupGroupDir(t)
	_ = addProfilesToGroup(dir, "prod", []string{"production"})

	exists, err := groupExists(dir, "prod")
	if err != nil {
		t.Fatalf("groupExists error: %v", err)
	}
	if !exists {
		t.Fatal("expected group 'prod' to exist")
	}

	exists, err = groupExists(dir, "missing")
	if err != nil {
		t.Fatalf("groupExists error: %v", err)
	}
	if exists {
		t.Fatal("expected group 'missing' to not exist")
	}
}

func TestProfilesInGroupNotFound(t *testing.T) {
	dir := setupGroupDir(t)
	_, err := profilesInGroup(dir, "ghost")
	if err == nil {
		t.Fatal("expected error for non-existent group")
	}
}

func TestRemoveProfileFromGroup(t *testing.T) {
	dir := setupGroupDir(t)
	_ = addProfilesToGroup(dir, "team", []string{"dev", "staging", "prod"})

	if err := removeProfileFromGroup(dir, "team", "staging"); err != nil {
		t.Fatalf("removeProfileFromGroup: %v", err)
	}
	profiles, _ := profilesInGroup(dir, "team")
	for _, p := range profiles {
		if p == "staging" {
			t.Fatal("expected 'staging' to be removed from group")
		}
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles remaining, got %d", len(profiles))
	}
}

func TestRemoveProfileFromNonExistentGroup(t *testing.T) {
	dir := setupGroupDir(t)
	err := removeProfileFromGroup(dir, "nope", "dev")
	if err == nil {
		t.Fatal("expected error removing from non-existent group")
	}
}

func TestGroupNamesEmptyDir(t *testing.T) {
	dir := setupGroupDir(t)
	names, err := groupNames(dir)
	if err != nil {
		t.Fatalf("groupNames: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected 0 group names, got %d", len(names))
	}
}
