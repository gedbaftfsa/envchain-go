package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/your-org/envchain-go/internal/store"
)

// CmdProtect marks specific keys in a project as protected, preventing
// accidental deletion or overwrite via set/unset without --force.
func CmdProtect(st *store.Store, project string, passphrase string, keys []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("at least one key must be specified")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	protectedRaw, _ := set.Get("__protected__")
	pinned := parseProtected(protectedRaw)

	added := 0
	for _, k := range keys {
		if k == "" {
			continue
		}
		if !pinned[k] {
			pinned[k] = true
			added++
		}
	}

	set.Put("__protected__", joinProtected(pinned))
	if err := st.Save(project, passphrase, set); err != nil {
		return err
	}

	fmt.Fprintf(w, "protected %d key(s) in project %q\n", added, project)
	return nil
}

// CmdUnprotect removes protection from specific keys in a project.
func CmdUnprotect(st *store.Store, project string, passphrase string, keys []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("at least one key must be specified")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	protectedRaw, _ := set.Get("__protected__")
	pinned := parseProtected(protectedRaw)

	for _, k := range keys {
		delete(pinned, k)
	}

	set.Put("__protected__", joinProtected(pinned))
	if err := st.Save(project, passphrase, set); err != nil {
		return err
	}

	fmt.Fprintf(w, "unprotected %d key(s) in project %q\n", len(keys), project)
	return nil
}

// CmdListProtected prints all protected keys for a project.
func CmdListProtected(st *store.Store, project string, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	protectedRaw, _ := set.Get("__protected__")
	pinned := parseProtected(protectedRaw)

	if len(pinned) == 0 {
		fmt.Fprintf(w, "no protected keys in project %q\n", project)
		return nil
	}

	keys := make([]string, 0, len(pinned))
	for k := range pinned {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintln(w, k)
	}
	return nil
}

func parseProtected(raw string) map[string]bool {
	m := make(map[string]bool)
	if raw == "" {
		return m
	}
	for _, k := range strings.Split(raw, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			m[k] = true
		}
	}
	return m
}

func joinProtected(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}
