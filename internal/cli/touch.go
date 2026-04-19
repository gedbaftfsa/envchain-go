package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/envchain-go/internal/store"
)

// CmdTouch updates the modification timestamp of a project by re-saving it
// with the same passphrase. This is useful to confirm a passphrase is still
// valid and to refresh any metadata tracked by the store.
func CmdTouch(st *store.Store, project, passphrase string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	if err := st.Save(project, passphrase, set); err != nil {
		return err
	}

	fmt.Fprintf(out, "touched %q at %s\n", project, time.Now().Format(time.RFC3339))
	return nil
}
