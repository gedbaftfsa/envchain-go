package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdEnvCheck compares the keys stored in a project against the keys
// currently present in the process environment and reports any that are
// missing or that differ from the stored value.
//
// Usage: envchain env-check <project>
func CmdEnvCheck(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	if len(keys) == 0 {
		fmt.Fprintln(w, "no keys stored for project")
		return nil
	}

	var missing, mismatched []string

	for _, k := range keys {
		stored, _ := es.Get(k)
		live, ok := os.LookupEnv(k)
		if !ok {
			missing = append(missing, k)
			continue
		}
		if live != stored {
			mismatched = append(mismatched, k)
		}
	}

	if len(missing) == 0 && len(mismatched) == 0 {
		fmt.Fprintln(w, "environment matches stored project")
		return nil
	}

	if len(missing) > 0 {
		fmt.Fprintf(w, "missing from environment: %s\n", strings.Join(missing, ", "))
	}
	if len(mismatched) > 0 {
		fmt.Fprintf(w, "value mismatch: %s\n", strings.Join(mismatched, ", "))
	}
	return fmt.Errorf("environment does not match stored project")
}
