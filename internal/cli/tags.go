package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdTags lists all unique key names (tags) across all projects in the store,
// optionally filtered by a prefix.
//
// Usage: envchain tags [<store-path>] [--prefix=<prefix>]
func CmdTags(st *store.Store, passphrase, prefix string, out fmt.Stringer) error {
	return cmdTags(st, passphrase, prefix, out)
}

func cmdTags(st *store.Store, passphrase, prefix string, out interface{ WriteString(string) (int, error) }) error {
	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("tags: list projects: %w", err)
	}

	seen := make(map[string]struct{})
	for _, proj := range projects {
		es, err := st.Load(proj, passphrase)
		if err != nil {
			return fmt.Errorf("tags: load %q: %w", proj, err)
		}
		for _, k := range es.Keys() {
			if prefix == "" || strings.HasPrefix(k, prefix) {
				seen[k] = struct{}{}
			}
		}
	}

	tags := make([]string, 0, len(seen))
	for k := range seen {
		tags = append(tags, k)
	}
	sort.Strings(tags)

	for _, t := range tags {
		out.WriteString(t + "\n") //nolint:errcheck
	}
	return nil
}
