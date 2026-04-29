package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdShell prints a shell-sourceable script that exports all variables
// for the given project. Unlike `run`, this does not exec a subprocess;
// it is intended to be eval'd in the current shell:
//
//	 eval "$(envchain shell myproject)"
func CmdShell(st *store.Store, passphrase, project string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	sort.Strings(keys)

	for _, k := range keys {
		v, _ := es.Get(k)
		fmt.Fprintf(w, "export %s=%s\n", k, shellQuote(v))
	}
	return nil
}

// CmdUnshell prints a shell-sourceable script that unsets all variables
// for the given project in the current shell:
//
//	 eval "$(envchain unshell myproject)"
func CmdUnshell(st *store.Store, passphrase, project string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "unset %s\n", k)
	}
	return nil
}

// shellScriptHeader returns a comment header for generated scripts.
func shellScriptHeader(project string) string {
	return fmt.Sprintf("# envchain-go: generated for project %q\n", project)
}

func init() {
	if !strings.Contains(os.Getenv("ENVCHAIN_NO_SHELL"), "1") {
		// intentionally blank — hook point for future shell-detection logic
		_ = shellScriptHeader
	}
}
