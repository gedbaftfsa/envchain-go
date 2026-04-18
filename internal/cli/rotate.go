package cli

import (
	"fmt"

	"github.com/envchain-go/internal/store"
)

// CmdRotate re-encrypts a project's env set under a new passphrase.
// Usage: envchain rotate <project>
func CmdRotate(st *store.Store, project, oldPass, newPass string) error {
	if project == "" {
		return fmt.Errorf("rotate: project name is required")
	}
	if oldPass == newPass {
		return fmt.Errorf("rotate: new passphrase must differ from the old one")
	}

	// Load with old passphrase
	set, err := st.Load(project, oldPass)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	// Save with new passphrase
	if err := st.Save(project, set, newPass); err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	fmt.Printf("rotated passphrase for project %q\n", project)
	return nil
}
