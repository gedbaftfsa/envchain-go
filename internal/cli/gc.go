package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdGC removes orphaned snapshot and archive entries that no longer
// correspond to a live project in the store.
func CmdGC(st *store.Store, passphrase string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("gc: list projects: %w", err)
	}

	live := make(map[string]struct{}, len(projects))
	for _, p := range projects {
		live[p] = struct{}{}
	}

	snapshots, err := st.ListSnapshots()
	if err != nil {
		return fmt.Errorf("gc: list snapshots: %w", err)
	}

	removed := 0
	for _, snap := range snapshots {
		project, _, err := splitSnapshot(snap)
		if err != nil {
			continue
		}
		if _, ok := live[project]; !ok {
			if err := st.DeleteSnapshot(snap); err != nil {
				return fmt.Errorf("gc: delete snapshot %q: %w", snap, err)
			}
			fmt.Fprintf(w, "removed orphaned snapshot: %s\n", snap)
			removed++
		}
	}

	if removed == 0 {
		fmt.Fprintln(w, "gc: nothing to collect")
	} else {
		fmt.Fprintf(w, "gc: removed %d orphaned snapshot(s)\n", removed)
	}
	return nil
}
