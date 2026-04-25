package cmd

import "sort"

// resolveAlias returns the profile name for the given alias, or the input
// unchanged if no alias is found.
func resolveAlias(dir, nameOrAlias string) (string, error) {
	aliases, err := loadAliases(dir)
	if err != nil {
		return "", err
	}
	if target, ok := aliases[nameOrAlias]; ok {
		return target, nil
	}
	return nameOrAlias, nil
}

// aliasesForProfile returns all alias names that point to the given profile.
func aliasesForProfile(dir, profile string) ([]string, error) {
	aliases, err := loadAliases(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	for alias, target := range aliases {
		if target == profile {
			result = append(result, alias)
		}
	}
	sort.Strings(result)
	return result, nil
}

// aliasCount returns the total number of defined aliases.
func aliasCount(dir string) (int, error) {
	aliases, err := loadAliases(dir)
	if err != nil {
		return 0, err
	}
	return len(aliases), nil
}

// aliasExists reports whether the given alias name is defined.
func aliasExists(dir, alias string) (bool, error) {
	aliases, err := loadAliases(dir)
	if err != nil {
		return false, err
	}
	_, ok := aliases[alias]
	return ok, nil
}
