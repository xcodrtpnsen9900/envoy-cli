package cmd

import (
	"fmt"
	"sort"
)

// groupExists returns true if the named group exists in the project.
func groupExists(dir, groupName string) (bool, error) {
	groups, err := loadGroups(dir)
	if err != nil {
		return false, err
	}
	_, ok := groups[groupName]
	return ok, nil
}

// profilesInGroup returns the sorted list of profiles belonging to a group.
func profilesInGroup(dir, groupName string) ([]string, error) {
	groups, err := loadGroups(dir)
	if err != nil {
		return nil, err
	}
	g, ok := groups[groupName]
	if !ok {
		return nil, fmt.Errorf("group %q not found", groupName)
	}
	out := make([]string, len(g.Profiles))
	copy(out, g.Profiles)
	sort.Strings(out)
	return out, nil
}

// groupNames returns a sorted list of all group names.
func groupNames(dir string) ([]string, error) {
	groups, err := loadGroups(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(groups))
	for k := range groups {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}

// removeProfileFromGroup removes a single profile from a group.
func removeProfileFromGroup(dir, groupName, profile string) error {
	groups, err := loadGroups(dir)
	if err != nil {
		return err
	}
	g, ok := groups[groupName]
	if !ok {
		return fmt.Errorf("group %q not found", groupName)
	}
	updated := g.Profiles[:0]
	for _, p := range g.Profiles {
		if p != profile {
			updated = append(updated, p)
		}
	}
	g.Profiles = updated
	groups[groupName] = g
	return saveGroups(dir, groups)
}
