package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdSearch searches for a key across all projects in the store.
func CmdSearch(st *store.Store, passphrase, query string) error {
	if query == "" {
		return fmt.Errorf("search: query must not be empty")
	}

	projects, err := st.ListProjects()
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	type match struct {
		project string
		key     string
	}

	var matches []match

	for _, proj := range projects {
		set, err := st.Load(proj, passphrase)
		if err != nil {
			return fmt.Errorf("search: load %q: %w", proj, err)
		}
		for _, k := range set.Keys() {
			if strings.Contains(strings.ToLower(k), strings.ToLower(query)) {
				matches = append(matches, match{proj, k})
			}
		}
	}

	if len(matches) == 0 {
		fmt.Printf("no keys matching %q found\n", query)
		return nil
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].project != matches[j].project {
			return matches[i].project < matches[j].project
		}
		return matches[i].key < matches[j].key
	})

	for _, m := range matches {
		fmt.Printf("%s\t%s\n", m.project, m.key)
	}
	return nil
}
