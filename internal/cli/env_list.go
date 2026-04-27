package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdEnvList prints the keys from a project that are also present in the
// current process environment, showing their live values alongside the stored
// ones. This lets you quickly see which stored variables are active in your
// shell session.
func CmdEnvList(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(w, "(no variables stored for %q)\n", project)
		return nil
	}

	sort.Strings(keys)

	fmt.Fprintf(w, "%-30s %-15s %s\n", "KEY", "STATUS", "LIVE VALUE")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 70))

	for _, k := range keys {
		live, present := os.LookupEnv(k)
		stored, _ := set.Get(k)

		var status string
		var display string

		switch {
		case !present:
			status = "missing"
			display = ""
		case live == stored:
			status = "match"
			display = truncate(live, 40)
		default:
			status = "differs"
			display = truncate(live, 40)
		}

		fmt.Fprintf(w, "%-30s %-15s %s\n", k, status, display)
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
