package cli

import (
	"fmt"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdRename renames a project namespace from oldName to newName.
// The passphrase is required to decrypt and re-encrypt the data.
func CmdRename(st *store.Store, oldName, newName, passphrase string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("rename: both old and new project names must be provided")
	}
	if oldName == newName {
		return fmt.Errorf("rename: old and new project names are identical")
	}

	set, err := st.Load(oldName, passphrase)
	if err != nil {
		return fmt.Errorf("rename: could not load project %q: %w", oldName, err)
	}

	// Check destination does not already exist.
	_, err = st.Load(newName, passphrase)
	if err == nil {
		return fmt.Errorf("rename: project %q already exists", newName)
	}

	if err := st.Save(newName, passphrase, set); err != nil {
		return fmt.Errorf("rename: could not save project %q: %w", newName, err)
	}

	if err := st.Delete(oldName); err != nil {
		// Best-effort rollback.
		_ = st.Delete(newName)
		return fmt.Errorf("rename: could not remove old project %q: %w", oldName, err)
	}

	fmt.Printf("Renamed project %q → %q\n", oldName, newName)
	return nil
}
