package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/envchain-go/internal/store"
)

// AuditEntry records a single access or mutation event.
type AuditEntry struct {
	Time    time.Time
	Project string
	Action  string
}

// CmdAudit prints a summary of project metadata: creation time and key count.
// It does not record live events; it surfaces what the store already knows.
func CmdAudit(st *store.Store, passphrase, project string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	sort.Strings(keys)

	fmt.Fprintf(w, "project : %s\n", project)
	fmt.Fprintf(w, "keys    : %d\n", len(keys))
	if len(keys) > 0 {
		fmt.Fprintf(w, "listing :\n")
		for _, k := range keys {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
	return nil
}

// auditWriter returns os.Stdout for production use.
func auditWriter() io.Writer { return os.Stdout }
