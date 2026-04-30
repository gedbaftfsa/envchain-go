package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdNamespace lists all unique namespace prefixes found across project keys.
// Keys are considered namespaced when they contain an underscore (e.g. DB_HOST → DB).
func CmdNamespace(st *store.Store, project string, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	if len(keys) == 0 {
		fmt.Fprintln(w, "(no variables)")
		return nil
	}

	seen := make(map[string]int)
	for _, k := range keys {
		ns := namespaceOf(k)
		seen[ns]++
	}

	namespaces := make([]string, 0, len(seen))
	for ns := range seen {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	for _, ns := range namespaces {
		fmt.Fprintf(w, "%-30s %d key(s)\n", ns, seen[ns])
	}
	return nil
}

// CmdNamespaceKeys lists all keys belonging to a given namespace prefix.
func CmdNamespaceKeys(st *store.Store, project, namespace, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	prefix := strings.ToUpper(namespace) + "_"
	keys := es.Keys()
	var matched []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			matched = append(matched, k)
		}
	}

	if len(matched) == 0 {
		fmt.Fprintf(w, "no keys found for namespace %q\n", namespace)
		return nil
	}

	sort.Strings(matched)
	for _, k := range matched {
		fmt.Fprintln(w, k)
	}
	return nil
}

func namespaceOf(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

func init() {
	registerCommand("namespace", "list key namespaces for a project")
}
