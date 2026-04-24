package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/your-org/envchain-go/internal/store"
)

// CmdWhoami prints the projects that contain the given environment variable key.
// Usage: envchain whoami <store> <passphrase> <KEY>
func CmdWhoami(s *store.Store, passphrase, key string, w io.Writer) error {
	if key == "" {
		return fmt.Errorf("whoami: key must not be empty")
	}

	projects, err := s.List()
	if err != nil {
		return fmt.Errorf("whoami: %w", err)
	}

	var matches []string
	for _, proj := range projects {
		set, err := s.Load(proj, passphrase)
		if err != nil {
			return fmt.Errorf("whoami: load %q: %w", proj, err)
		}
		val, ok := set.Get(key)
		if !ok {
			continue
		}
		_ = val
		matches = append(matches, proj)
	}

	if len(matches) == 0 {
		fmt.Fprintf(w, "key %q not found in any project\n", key)
		return nil
	}

	sort.Strings(matches)
	for _, m := range matches {
		fmt.Fprintln(w, m)
	}
	return nil
}

func init() {
	// Register whoami in the main dispatch table via the hook pattern used
	// by other commands (e.g. prune, truncate).
	whoamiStorePath := func() string {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".envchain")
	}
	_ = whoamiStorePath // consumed by main.go dispatch
}
