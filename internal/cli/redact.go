package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdRedact reads stdin (or a provided string) and replaces any known secret
// values for the given project with "***REDACTED***". Useful for sanitising
// log output before sharing.
func CmdRedact(st *store.Store, project, passphrase string, input string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	if len(keys) == 0 {
		_, err = fmt.Fprint(out, input)
		return err
	}

	result := input
	for _, k := range keys {
		v, ok := es.Get(k)
		if !ok || v == "" {
			continue
		}
		result = strings.ReplaceAll(result, v, "***REDACTED***")
	}

	_, err = fmt.Fprint(out, result)
	return err
}

// CmdRedactFile reads the named file, redacts secrets, and writes the result
// to out. Pass "-" as path to read from stdin.
func CmdRedactFile(st *store.Store, project, passphrase, path string, out io.Writer) error {
	var raw []byte
	var err error

	if path == "-" {
		raw, err = io.ReadAll(os.Stdin)
	} else {
		raw, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	return CmdRedact(st, project, passphrase, string(raw), out)
}
