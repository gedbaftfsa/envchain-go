package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdFmt normalises the stored env-set for a project: keys are sorted
// alphabetically and re-saved with the same passphrase.
func CmdFmt(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	sort.Strings(keys)

	sorted := env.NewSet()
	for _, k := range keys {
		v, _ := es.Get(k)
		_ = sorted.Put(k, v)
	}

	if err := st.Save(project, passphrase, sorted); err != nil {
		return err
	}

	fmt.Fprintf(w, "formatted %s (%d keys)\n", project, len(keys))
	return nil
}

// CmdFmtDiff prints what would change without writing, using a simple
// before/after key-order comparison.
func CmdFmtDiff(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	before := es.Keys()
	sorted := make([]string, len(before))
	copy(sorted, before)
	sort.Strings(sorted)

	same := strings.Join(before, ",") == strings.Join(sorted, ",")
	if same {
		fmt.Fprintf(w, "%s is already sorted\n", project)
		return nil
	}

	fmt.Fprintf(w, "would reorder keys for %s:\n", project)
	for i, k := range sorted {
		fmt.Fprintf(w, "  [%d] %s\n", i+1, k)
	}
	return nil
}
