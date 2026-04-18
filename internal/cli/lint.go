package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdLint checks a project's env set for common issues such as duplicate keys,
// empty values, or keys that shadow well-known system variables.
func CmdLint(st *store.Store, project, passphrase string, w io.Writer) error {
	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(w, "project %q has no variables\n", project)
		return nil
	}

	issue := false
	shadowed := shadowedVars()

	for _, k := range keys {
		v, _ := set.Get(k)
		if strings.TrimSpace(v) == "" {
			fmt.Fprintf(w, "WARN  %s: value is empty or whitespace-only\n", k)
			issue = true
		}
		if _, ok := shadowed[strings.ToUpper(k)]; ok {
			fmt.Fprintf(w, "WARN  %s: shadows a well-known system variable\n", k)
			issue = true
		}
	}

	if !issue {
		fmt.Fprintf(w, "OK    %s: no issues found (%d variable(s))\n", project, len(keys))
	}
	return nil
}

func shadowedVars() map[string]struct{} {
	well := []string{"PATH", "HOME", "USER", "SHELL", "TERM", "LANG", "PWD", "TMPDIR", "TMP", "TEMP"}
	m := make(map[string]struct{}, len(well))
	for _, v := range well {
		m[v] = struct{}{}
	}
	// also add current process env keys that are already set
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[strings.ToUpper(parts[0])] = struct{}{}
		}
	}
	return m
}
