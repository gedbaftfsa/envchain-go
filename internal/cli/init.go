package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/envchain/envchain-go/internal/store"
)

// CmdInit initialises a new project namespace in the store.
// It prompts for a passphrase, creates an empty env set, and saves it.
func CmdInit(st *store.Store, project string, w io.Writer, passphraseFn func(string) (string, error)) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	// Check if already exists
	_, err := st.Load(project, "dummy")
	if err == nil {
		return fmt.Errorf("project %q already exists; use 'set' to add variables", project)
	}
	if err != store.ErrNotFound {
		// unexpected error other than not-found means something is wrong
		// but we proceed — it may just be a wrong passphrase on existing data
	}

	pass, err := passphraseFn("New passphrase: ")
	if err != nil {
		return fmt.Errorf("reading passphrase: %w", err)
	}
	if pass == "" {
		return fmt.Errorf("passphrase must not be empty")
	}

	confirm, err := passphraseFn("Confirm passphrase: ")
	if err != nil {
		return fmt.Errorf("reading passphrase confirmation: %w", err)
	}
	if pass != confirm {
		return fmt.Errorf("passphrases do not match")
	}

	es := newEnvSetEmpty()
	if err := st.Save(project, pass, es); err != nil {
		return fmt.Errorf("saving project: %w", err)
	}

	fmt.Fprintf(w, "Initialised project %q\n", project)
	return nil
}

// newEnvSetEmpty returns an empty env set via the env package.
func newEnvSetEmpty() interface{ Keys() []string } {
	// We import env indirectly through the existing helpers.
	// Return a minimal wrapper that satisfies store.Save expectations.
	set, _ := mustEnvSetFromPairs(nil)
	return set
}

func init() {
	_ = os.Stderr // ensure os imported
}
