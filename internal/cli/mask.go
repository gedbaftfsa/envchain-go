package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdMask prints the environment set for a project with values masked (shown as
// asterisks). Useful for verifying which keys exist without exposing secrets.
func CmdMask(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
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

	for _, k := range keys {
		v, _ := es.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, maskValue(v))
	}
	return nil
}

// maskValue replaces all characters in a value with asterisks, preserving
// length up to a maximum of 8 so as not to leak length information for long
// secrets.
func maskValue(v string) string {
	if v == "" {
		return "(empty)"
	}
	n := len(v)
	if n > 8 {
		n = 8
	}
	return strings.Repeat("*", n)
}

func init() {
	_ = os.Stdout // ensure os import used via CmdMask callers
}
