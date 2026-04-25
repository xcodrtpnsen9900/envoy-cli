package cmd

import "strings"

// accessCount returns the total number of access log entries.
func accessCount(root string) int {
	entries, err := readAccessLog(root)
	if err != nil {
		return 0
	}
	return len(entries)
}

// accessCountForProfile returns the number of log entries for a specific profile.
func accessCountForProfile(root, profile string) int {
	entries, err := readAccessLog(root)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if strings.EqualFold(e.Profile, profile) {
			count++
		}
	}
	return count
}

// lastAccessEntry returns the most recent access log entry, or nil if none.
func lastAccessEntry(root string) *AccessEntry {
	entries, err := readAccessLog(root)
	if err != nil || len(entries) == 0 {
		return nil
	}
	last := entries[len(entries)-1]
	return &last
}

// clearAccessLog removes all access log entries.
func clearAccessLog(root string) error {
	return writeAccessLog(root, nil)
}

func writeAccessLog(root string, entries []AccessEntry) error {
	import_json_data, err := marshalAccessEntries(entries)
	if err != nil {
		return err
	}
	return writeFile(accessLogPath(root), import_json_data)
}
