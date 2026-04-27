package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdRecap prints a summary of all keys in a project without revealing values.
// It shows key names, whether they are set, and their approximate value length.
func CmdRecap(st *store.Store, project, passphrase string, out io.Writer) error {
	if strings.TrimSpace(project) == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(out, "project %q has no variables\n", project)
		return nil
	}

	sort.Strings(keys)

	fmt.Fprintf(out, "project: %s (%d variable(s))\n", project, len(keys))
	fmt.Fprintln(out, strings.Repeat("-", 40))

	for _, k := range keys {
		v, _ := set.Get(k)
		length := len(v)
		status := "set"
		if length == 0 {
			status = "empty"
		}
		fmt.Fprintf(out, "  %-24s %s (~%d chars)\n", k, status, length)
	}

	return nil
}

func init() {
	recapUsage := "recap <project>  — summarise keys in a project without showing values"
	_ = recapUsage // referenced by Main dispatch
	_ = os.Stdout  // ensure os import used
}
