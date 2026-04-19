package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/envchain-go/internal/store"
)

// CmdHistory prints a chronological audit trail of snapshots for a project.
func CmdHistory(st *store.Store, project, passphrase string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	snaps, err := st.ListSnapshots(project)
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err)
	}

	if len(snaps) == 0 {
		fmt.Fprintf(out, "no history found for project %q\n", project)
		return nil
	}

	fmt.Fprintf(out, "history for project %q:\n", project)
	for _, name := range snaps {
		_, ts, err := splitSnapshot(name)
		if err != nil {
			fmt.Fprintf(out, "  %s\n", name)
			continue
		}
		t := time.Unix(ts, 0).UTC()
		fmt.Fprintf(out, "  %s  (%s)\n", name, t.Format(time.RFC3339))
	}
	return nil
}

// cmdHistoryMain is the entry point wired into Main.
func cmdHistoryMain(args []string, st *store.Store, pass func(string) (string, error)) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envchain history <project>")
		return 1
	}
	project := args[0]
	passphrase, err := pass(project)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := CmdHistory(st, project, passphrase, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
