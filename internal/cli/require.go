package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envchain-go/internal/store"
)

// CmdRequire checks that all specified keys exist and have non-empty values
// in the given project, exiting non-zero if any are missing.
//
// Usage: envchain require <project> <KEY1> [KEY2 ...]
func CmdRequire(st *store.Store, passphrase, project string, keys []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("at least one key must be specified")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	var missing []string
	for _, k := range keys {
		v, ok := es.Get(k)
		if !ok || strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		for _, k := range missing {
			fmt.Fprintf(w, "missing: %s\n", k)
		}
		return fmt.Errorf("%d required key(s) missing or empty in project %q", len(missing), project)
	}

	fmt.Fprintf(w, "ok: all %d required key(s) present in %q\n", len(keys), project)
	return nil
}

func init() {
	_ = os.Stderr // ensure os imported
}
