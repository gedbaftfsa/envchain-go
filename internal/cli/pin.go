package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdPin marks one or more keys in a project as "pinned", storing them in a
// special __pinned__ meta-key so that CmdLint and CmdDiff can highlight them.
func CmdPin(st *store.Store, project, passphrase string, keys []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}
	if len(keys) == 0 {
		return fmt.Errorf("at least one key required")
	}
	for _, k := range keys {
		if strings.TrimSpace(k) == "" {
			return fmt.Errorf("key must not be blank")
		}
	}

	es, err := loadOrNew(st, project, passphrase)
	if err != nil {
		return err
	}

	// Retrieve existing pinned set.
	existing := ""
	if v, ok := es.Get("__pinned__"); ok {
		existing = v
	}
	pinned := parsePinned(existing)
	for _, k := range keys {
		pinned[k] = struct{}{}
	}
	es.Put("__pinned__", joinPinned(pinned))

	if err := st.Save(project, passphrase, es); err != nil {
		return err
	}
	fmt.Fprintf(w, "pinned %d key(s) in project %q\n", len(keys), project)
	return nil
}

// CmdUnpin removes keys from the pinned set of a project.
func CmdUnpin(st *store.Store, project, passphrase string, keys []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}
	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}
	existing := ""
	if v, ok := es.Get("__pinned__"); ok {
		existing = v
	}
	pinned := parsePinned(existing)
	for _, k := range keys {
		delete(pinned, k)
	}
	if len(pinned) == 0 {
		es.Delete("__pinned__")
	} else {
		es.Put("__pinned__", joinPinned(pinned))
	}
	if err := st.Save(project, passphrase, es); err != nil {
		return err
	}
	fmt.Fprintf(w, "unpinned %d key(s) in project %q\n", len(keys), project)
	return nil
}

// CmdListPinned prints the pinned keys for a project.
func CmdListPinned(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}
	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}
	v, ok := es.Get("__pinned__")
	if !ok || v == "" {
		fmt.Fprintf(w, "no pinned keys in project %q\n", project)
		return nil
	}
	for k := range parsePinned(v) {
		fmt.Fprintln(w, k)
	}
	return nil
}

func parsePinned(s string) map[string]struct{} {
	m := map[string]struct{}{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			m[p] = struct{}{}
		}
	}
	return m
}

func joinPinned(m map[string]struct{}) string {
	parts := make([]string, 0, len(m))
	for k := range m {
		parts = append(parts, k)
	}
	_ = os.Stderr // suppress unused import
	return strings.Join(parts, ",")
}
