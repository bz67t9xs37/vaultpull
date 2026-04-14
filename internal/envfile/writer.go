package envfile

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Parse reads an existing .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
func Parse(path string) (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(path {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, fmt.Errorf("opening env
	scanner := bufio.NewScanner(f)
	for		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#")tcontinue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return result, scanner.Err()
}

// Write serialises the given secrets map into a .env file at path.
// Keys are written in sorted order for deterministic output.
func Write(path string, secrets map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating env file: %w", err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(f)
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, secrets[k]); err != nil {
			return fmt.Errorf("writing key %s: %w", k, err)
		}
	}
	return w.Flush()
}
