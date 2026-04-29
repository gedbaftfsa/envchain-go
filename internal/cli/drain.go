package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdDrain removes all keys from a project, leaving it empty but intact.
// Unlike `delete`, the project entry itself is preserved in the store.
func CmdDrain(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(w, "project %q is already empty\n", project)
		return nil
	}

	sort.Strings(keys)
	for _, k := range keys {
		es.Delete(k)
	}

	if err := st.Save(project, passphrase, es); err != nil {
		return err
	}

	fmt.Fprintf(w, "drained %d key(s) from project %q\n", len(keys), project)
	return nil
}

func init() {
	_ = os.Stderr // ensure os import used in future extensions
}
