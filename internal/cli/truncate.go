package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdTruncate removes all but the N most recent snapshots for a project.
func CmdTruncate(st *store.Store, project, passphrase string, keep int, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}
	if keep < 1 {
		return fmt.Errorf("keep must be at least 1")
	}

	snaps, err := st.ListSnapshots(project)
	if err != nil {
		return err
	}

	if len(snaps) <= keep {
		fmt.Fprintf(w, "nothing to truncate (%d snapshot(s), keeping %d)\n", len(snaps), keep)
		return nil
	}

	// Snapshots are stored oldest-first; drop the oldest ones.
	toRemove := snaps[:len(snaps)-keep]
	for _, name := range toRemove {
		key := project + ":" + name
		if err := st.Delete(key); err != nil {
			return fmt.Errorf("delete snapshot %q: %w", name, err)
		}
		fmt.Fprintf(w, "removed snapshot %s\n", name)
	}
	return nil
}

func init() {
	registerCommand("truncate", func(args []string, st *store.Store, pass func(string) (string, error), w io.Writer) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: truncate <project> <keep>")
		}
		project := args[0]
		keep, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("keep must be an integer: %w", err)
		}
		passphrase, err := pass(project)
		if err != nil {
			return err
		}
		_ = passphrase
		return CmdTruncate(st, project, passphrase, keep, os.Stdout)
	})
}
