package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envchain-go/internal/store"
)

// CmdStats prints statistics about all projects in the store.
func CmdStats(st *store.Store, passphrase string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	projects, err := st.ListProjects()
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}

	if len(projects) == 0 {
		fmt.Fprintln(w, "No projects found.")
		return nil
	}

	sort.Strings(projects)

	totalVars := 0
	fmt.Fprintf(w, "%-24s  %s\n", "PROJECT", "VARS")
	fmt.Fprintf(w, "%-24s  %s\n", "-------", "----")

	for _, name := range projects {
		set, err := st.Load(name, passphrase)
		if err != nil {
			return fmt.Errorf("stats: load %q: %w", name, err)
		}
		count := len(set.Keys())
		totalVars += count
		fmt.Fprintf(w, "%-24s  %d\n", name, count)
	}

	fmt.Fprintf(w, "\nTotal: %d project(s), %d variable(s)\n", len(projects), totalVars)
	return nil
}
