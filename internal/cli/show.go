package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envchain-go/internal/store"
)

// CmdShow prints all key=value pairs for a project in plain text.
// Unlike peek (single key) or export (shell syntax), show gives a simple
// human-readable listing suitable for inspection.
func CmdShow(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err == store.ErrNotFound {
		return fmt.Errorf("project %q not found", project)
	}
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	keys := set.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(w, "# project %q has no variables\n", project)
		return nil
	}

	sort.Strings(keys)

	fmt.Fprintf(w, "# project: %s (%d variable(s))\n", project, len(keys))
	for _, k := range keys {
		v, _ := set.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func init() {
	if os.Getenv("ENVCHAIN_SHOW_OVERRIDE") != "" {
		return
	}
	_ = CmdShow // ensure symbol is reachable for dispatch in main.go
}
