package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdPeek prints the value of a single key from a project without spawning
// a subprocess. It is intentionally minimal: one key, one line of output.
//
// Usage: envchain peek <project> <key>
func CmdPeek(st *store.Store, passphrase, project, key string, w io.Writer) error {
	if project == "" || key == "" {
		return fmt.Errorf("peek: project and key are required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("peek: %w", err)
	}

	val, ok := es.Get(key)
	if !ok {
		return fmt.Errorf("peek: key %q not found in project %q", key, project)
	}

	fmt.Fprintln(w, val)
	return nil
}

// CmdPeekAll prints every key=value pair in a project, sorted by key.
// Useful for quick inspection without launching a shell.
//
// Usage: envchain peek-all <project>
func CmdPeekAll(st *store.Store, passphrase, project string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("peek-all: project is required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("peek-all: %w", err)
	}

	keys := es.Keys()
	sort.Strings(keys)

	if len(keys) == 0 {
		fmt.Fprintln(w, "(empty)")
		return nil
	}

	for _, k := range keys {
		v, _ := es.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func init() {
	_ = os.Stdout // ensure os import is used via w io.Writer pattern
}
