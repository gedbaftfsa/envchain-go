package cli

import (
	"fmt"
	"io"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdRekey re-encrypts all projects in the store with a new passphrase.
// Unlike rotate (single project), rekey operates on every project at once.
func CmdRekey(st *store.Store, oldPass, newPass string, out io.Writer) error {
	if oldPass == "" || newPass == "" {
		return fmt.Errorf("rekey: passphrase must not be empty")
	}
	if oldPass == newPass {
		return fmt.Errorf("rekey: new passphrase must differ from old")
	}

	projects, err := st.List()
	if err != nil {
		return fmt.Errorf("rekey: list projects: %w", err)
	}
	if len(projects) == 0 {
		fmt.Fprintln(out, "rekey: no projects found")
		return nil
	}

	var rekeyed []string
	for _, name := range projects {
		set, err := st.Load(name, oldPass)
		if err != nil {
			return fmt.Errorf("rekey: load %q: %w", name, err)
		}
		if err := st.Save(name, set, newPass); err != nil {
			return fmt.Errorf("rekey: save %q: %w", name, err)
		}
		rekeyed = append(rekeyed, name)
	}

	fmt.Fprintf(out, "rekeyed %d project(s):\n", len(rekeyed))
	for _, name := range rekeyed {
		fmt.Fprintf(out, "  %s\n", name)
	}
	return nil
}
