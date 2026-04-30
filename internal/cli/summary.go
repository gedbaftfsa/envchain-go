package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdSummary prints a human-readable summary of all projects in the store,
// showing the project name, number of keys, and whether any keys are pinned
// or protected.
func CmdSummary(st *store.Store, passphrase string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("summary: list projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Fprintln(w, "no projects found")
		return nil
	}

	sort.Strings(projects)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROJECT\tKEYS\tPINNED\tPROTECTED")

	for _, name := range projects {
		set, err := st.Load(name, passphrase)
		if err != nil {
			return fmt.Errorf("summary: load %q: %w", name, err)
		}

		keys := set.Keys()
		pinnedCount := countMeta(set.Meta("pinned"))
		protectedCount := countMeta(set.Meta("protected"))

		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\n", name, len(keys), pinnedCount, protectedCount)
	}

	return tw.Flush()
}

func countMeta(val string) int {
	if val == "" {
		return 0
	}
	parts := splitNonEmpty(val, ",")
	return len(parts)
}

func splitNonEmpty(s, sep string) []string {
	var out []string
	for _, p := range splitString(s, sep) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
