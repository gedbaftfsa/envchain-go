package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/envchain-go/internal/store"
)

// CmdVerify checks that a project's encrypted data can be successfully
// decrypted with the provided passphrase, without printing any values.
func CmdVerify(st *store.Store, project, passphrase string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}
	_, err := st.Load(project, passphrase)
	if err == store.ErrNotFound {
		return fmt.Errorf("project %q not found", project)
	}
	if err != nil {
		return fmt.Errorf("verify failed for %q: %w", project, err)
	}
	fmt.Fprintf(out, "OK: project %q verified successfully\n", project)
	return nil
}

func init() {
	_ = os.Stderr // ensure os import used
}
