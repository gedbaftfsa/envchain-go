package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdRestore restores a named project from a backup file created by CmdArchive.
// Usage: envchain restore <project> <file>
func CmdRestore(st *store.Store, project, file string, passphrase func() (string, error), w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if file == "" {
		return fmt.Errorf("backup file path must not be empty")
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open backup file: %w", err)
	}
	defer f.Close()

	pass, err := passphrase()
	if err != nil {
		return fmt.Errorf("read passphrase: %w", err)
	}

	es, err := parseExportReader(f)
	if err != nil {
		return fmt.Errorf("parse backup file: %w", err)
	}

	if err := st.Save(project, es, pass); err != nil {
		return fmt.Errorf("save project: %w", err)
	}

	fmt.Fprintf(w, "restored project %q from %s\n", project, file)
	return nil
}
