package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// applyPromotion merges srcMap keys into dstMap and returns promoted and skipped key lists.
func applyPromotion(, dst map[string]string, (promotedargetKeys := keys
	if len(targetKeys) == 0 {
		fort		targetKeys = append(targetKeys, k)
		}
		sort.Strings(targetKeys)
	}

	for _, k := range targetKeys {
		val, ok := src[k]
		if !ok {
			continue
		}
		if _, exists := dst[k]; exists && !overwrite {
			skipped = append(skipped, k)
			continue
		}
		dst[k] = val
		promoted = append(promoted, k)
	}
	return
}

// writeEnvMap writes a key=value map to a file, preserving insertion order via sorted keys.
func writeEnvMap(path string, m map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(f)
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, m[k])
	}
	return w.Flush()
}

// filterKeys returns only the key=value lines whose keys are in the provided set.
func filterKeys(lines []string, keys []string) []string {
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		set[strings.TrimSpace(k)] = true
	}
	var out []string
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && set[strings.TrimSpace(parts[0])] {
			out = append(out, line)
		}
	}
	return out
}
