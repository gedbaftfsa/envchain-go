package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envchain-go/internal/env"
	"github.com/user/envchain-go/internal/store"
)

// CmdExport prints environment variables for a project in shell-sourceable format.
// Format: export KEY=VALUE
func CmdExport(project, passphrase string, st *store.Store, w io.Writer) error {
	set, err := st.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("load project %q: %w", project, err)
	}

	keys := set.Keys()
	sort.Strings(keys)

	for _, k := range keys {
		v, _ := set.Get(k)
		fmt.Fprintf(w, "export %s=%s\n", k, shellQuote(v))
	}
	return nil
}

// CmdImport reads KEY=VALUE lines from r and stores them into the project.
func CmdImport(project, passphrase string, st *store.Store, r io.Reader, raw []byte) error {
	set, err := loadOrNew(project, passphrase, st)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, err := env.ParseEntry(line)
		if err != nil {
			return fmt.Errorf("parse line %q: %w", line, err)
		}
		if err := set.Put(k, v); err != nil {
			return fmt.Errorf("put %q: %w", k, err)
		}
	}

	return st.Save(project, passphrase, set)
}

// shellQuote wraps a value in single quotes, escaping existing single quotes.
func shellQuote(s string) string {
	if !strings.ContainsAny(s, " \t\n\"'\\$`!&|;<>(){}#~") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\'''") + "'"
}

// init registers export/import in main usage; actual wiring is in main.go.
var _ = os.Stdout
