package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/your-org/envchain-go/internal/env"
	"github.com/your-org/envchain-go/internal/store"
)

// CmdEnvDiff compares a stored project's variables against the current process
// environment, reporting keys that are missing, extra, or have changed values.
func CmdEnvDiff(st *store.Store, project, passphrase string, w io.Writer) error {
	if strings.TrimSpace(project) == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	current := env.FromProcess()

	stored := make(map[string]string)
	for _, k := range set.Keys() {
		v, _ := set.Get(k)
		stored[k] = v
	}

	var missing, extra, changed []string

	for k, sv := range stored {
		cv, ok := current[k]
		if !ok {
			missing = append(missing, k)
		} else if cv != sv {
			changed = append(changed, k)
		}
	}

	for k := range current {
		if _, ok := stored[k]; !ok {
			extra = append(extra, k)
		}
	}

	sort.Strings(missing)
	sort.Strings(extra)
	sort.Strings(changed)

	if len(missing)+len(extra)+len(changed) == 0 {
		fmt.Fprintln(w, "no differences")
		return nil
	}

	for _, k := range missing {
		fmt.Fprintf(w, "- %s (missing from environment)\n", k)
	}
	for _, k := range changed {
		fmt.Fprintf(w, "~ %s (value differs)\n", k)
	}
	for _, k := range extra {
		fmt.Fprintf(w, "+ %s (not in project)\n", k)
	}

	return nil
}
