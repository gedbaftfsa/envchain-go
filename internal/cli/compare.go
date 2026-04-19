package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdCompare prints a side-by-side key comparison of two projects,
// showing which keys are missing or present in each.
func CmdCompare(st *store.Store, passphrase, projectA, projectB string, w io.Writer) error {
	if projectA == "" || projectB == "" {
		return fmt.Errorf("compare: two project names required")
	}

	setA, err := st.Load(projectA, passphrase)
	if err != nil {
		return fmt.Errorf("compare: loading %q: %w", projectA, err)
	}
	setB, err := st.Load(projectB, passphrase)
	if err != nil {
		return fmt.Errorf("compare: loading %q: %w", projectB, err)
	}

	keysA := keySet(setA.Keys())
	keysB := keySet(setB.Keys())

	all := make(map[string]struct{})
	for k := range keysA {
		all[k] = struct{}{}
	}
	for k := range keysB {
		all[k] = struct{}{}
	}

	sorted := make([]string, 0, len(all))
	for k := range all {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	fmt.Fprintf(w, "%-40s %-10s %-10s\n", "KEY", projectA, projectB)
	fmt.Fprintf(w, "%-40s %-10s %-10s\n", "---", "---", "---")
	for _, k := range sorted {
		_, inA := keysA[k]
		_, inB := keysB[k]
		colA, colB := "missing", "missing"
		if inA {
			colA = "present"
		}
		if inB {
			colB = "present"
		}
		fmt.Fprintf(w, "%-40s %-10s %-10s\n", k, colA, colB)
	}
	return nil
}

func init() {
	_ = os.Stdout // ensure os import used
}
