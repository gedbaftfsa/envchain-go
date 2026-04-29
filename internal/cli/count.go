package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdCount prints the number of keys stored in each project, or a specific
// project if one is provided. Output is sorted alphabetically by project name.
func CmdCount(st *store.Store, passphrase string, args []string, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("count: list projects: %w", err)
	}

	// Filter to a single project if an argument was provided.
	if len(args) > 0 {
		target := args[0]
		found := false
		for _, p := range projects {
			if p == target {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("count: project %q not found", target)
		}
		projects = []string{target}
	}

	if len(projects) == 0 {
		fmt.Fprintln(out, "no projects found")
		return nil
	}

	sort.Strings(projects)

	for _, name := range projects {
		es, err := st.Load(name, passphrase)
		if err != nil {
			return fmt.Errorf("count: load %q: %w", name, err)
		}
		fmt.Fprintf(out, "%-30s %d\n", name, len(es.Keys()))
	}

	return nil
}

func init() {
	registerCommand("count", "count", "Print key count per project", func(st *store.Store, pass string, args []string) error {
		return CmdCount(st, pass, args, nil)
	})
}
