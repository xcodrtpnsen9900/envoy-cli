package cmd

import (
	"fmt"
	"os"
)

// clearAuditLog removes all entries from the audit log.
func clearAuditLog(dir string) error {
	path := auditLogPath(dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

// auditEntryCount returns the number of entries in the audit log.
func auditEntryCount(dir string) (int, error) {
	entries, err := readAuditLog(dir)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// lastAuditEntry returns the most recent audit log entry, or empty string.
func lastAuditEntry(dir string) (string, error) {
	entries, err := readAuditLog(dir)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", nil
	}
	return entries[len(entries)-1], nil
}

// auditAction writes an audit entry and prints a confirmation.
func auditAction(dir, action, profile string) {
	writeAuditEntry(dir, action, profile)
	fmt.Printf("Audit: recorded '%s' for profile '%s'\n", action, profile)
}
