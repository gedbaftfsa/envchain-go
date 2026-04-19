package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdPrune removes all projects from the store that have no environment
// variables set, printing each removed project name to w.
func CmdPrune(st *store.Store, passphrase string, w io.Writer) error {
	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("prune: list projects: %w", err)
	}

	removed := 0
	for _, proj := range projects {
		es, err := st.Load(proj, passphrase)
		if err != nil {
			return fmt.Errorf("prune: load %q: %w", proj, err)
		}
		if len(es.Keys()) == 0 {
			if err := st.Delete(proj); err != nil {
				return fmt.Errorf("prune: delete %q: %w", proj, err)
			}
			fmt.Fprintf(w, "pruned %s\n", proj)
			removed++
		}
	}

	if removed == 0 {
		fmt.Fprintln(w, "nothing to prune")
	}
	return nil
}

func init() {
	_ = os.Stderr // ensure os imported for future use
}
