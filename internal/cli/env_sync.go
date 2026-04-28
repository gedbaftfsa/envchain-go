package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdEnvSync reconciles the stored environment set for a project against the
// current process environment. Keys present in the store but absent from the
// process environment are reported as "missing"; keys present in the process
// environment but absent from the store are reported as "extra". When the
// --apply flag is provided the missing keys are written into the store using
// their current process values.
func CmdEnvSync(st *store.Store, project, passphrase string, apply bool, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	// Build a lookup of what the store currently holds.
	storeKeys := make(map[string]string)
	for _, k := range set.Keys() {
		v, _ := set.Get(k)
		storeKeys[k] = v
	}

	// Build a lookup of the current process environment.
	processEnv := env.FromProcess()

	var missing []string // in store, not in process
	var extra []string   // in process, not in store

	for k := range storeKeys {
		if _, ok := processEnv[k]; !ok {
			missing = append(missing, k)
		}
	}
	for k := range processEnv {
		if _, ok := storeKeys[k]; !ok {
			extra = append(extra, k)
		}
	}

	sort.Strings(missing)
	sort.Strings(extra)

	if len(missing) == 0 && len(extra) == 0 {
		fmt.Fprintln(w, "store and process environment are in sync")
		return nil
	}

	if len(missing) > 0 {
		fmt.Fprintf(w, "missing from process environment (%d):\n", len(missing))
		for _, k := range missing {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(extra) > 0 {
		fmt.Fprintf(w, "extra in process environment (%d):\n", len(extra))
		for _, k := range extra {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if !apply {
		return nil
	}

	// --apply: promote extra process-env keys into the store.
	if len(extra) == 0 {
		fmt.Fprintln(w, "nothing to apply")
		return nil
	}

	for _, k := range extra {
		val := os.Getenv(k)
		if err := set.Put(k, val); err != nil {
			return fmt.Errorf("put %s: %w", k, err)
		}
	}

	if err := st.Save(project, passphrase, set); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	fmt.Fprintf(w, "applied %d key(s) from process environment into project %q\n", len(extra), project)
	return nil
}
