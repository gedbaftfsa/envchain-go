package cli

import (
	"fmt"
	"time"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdSnapshot saves a timestamped copy of a project's env set.
func CmdSnapshot(st *store.Store, project, passphrase string) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("load %q: %w", project, err)
	}

	tag := time.Now().UTC().Format("20060102T150405Z")
	snapName := fmt.Sprintf("%s@%s", project, tag)

	if err := st.Save(snapName, passphrase, set); err != nil {
		return fmt.Errorf("save snapshot %q: %w", snapName, err)
	}

	fmt.Printf("snapshot saved: %s\n", snapName)
	return nil
}

// CmdRestoreSnapshot copies a snapshot back to the project name.
func CmdRestoreSnapshot(st *store.Store, snapName, passphrase string) error {
	if snapName == "" {
		return fmt.Errorf("snapshot name required")
	}

	set, err := st.Load(snapName, passphrase)
	if err != nil {
		return fmt.Errorf("load snapshot %q: %w", snapName, err)
	}

	project, _, ok := splitSnapshot(snapName)
	if !ok {
		return fmt.Errorf("%q does not look like a snapshot (expected project@timestamp)", snapName)
	}

	if err := st.Save(project, passphrase, set); err != nil {
		return fmt.Errorf("restore to %q: %w", project, err)
	}

	fmt.Printf("restored %s -> %s\n", snapName, project)
	return nil
}

// CmdListSnapshots prints all snapshots for a project.
func CmdListSnapshots(st *store.Store, project string) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	projects, err := st.ListProjects()
	if err != nil {
		return err
	}

	prefix := project + "@"
	found := false
	for _, p := range projects {
		if len(p) > len(prefix) && p[:len(prefix)] == prefix {
			fmt.Println(p)
			found = true
		}
	}
	if !found {
		fmt.Printf("no snapshots found for %q\n", project)
	}
	return nil
}

func splitSnapshot(name string) (project, tag string, ok bool) {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '@' {
			return name[:i], name[i+1:], true
		}
	}
	return "", "", false
}

var _ = env.NewSet // keep import
