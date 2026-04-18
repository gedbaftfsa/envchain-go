package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdListProjects prints all project names that have stored env sets.
func CmdListProjects(st *store.Store, args []string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	projects, err := st.ListProjects()
	if err != nil {
		return fmt.Errorf("list projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Fprintln(w, "(no projects found)")
		return nil
	}

	sort.Strings(projects)
	for _, p := range projects {
		fmt.Fprintln(w, p)
	}
	return nil
}
