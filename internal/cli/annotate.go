package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdAnnotate sets or displays a free-text annotation (description) for a
// project's environment variable set.
//
// Usage:
//
//	envchain annotate <project>              # print current annotation
//	envchain annotate <project> <text...>    # set annotation
func CmdAnnotate(st *store.Store, passphrase, project string, args []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	// Read-only: no extra args → print existing annotation.
	if len(args) == 0 {
		annotation := set.Annotation()
		if annotation == "" {
			fmt.Fprintf(w, "(no annotation)\n")
		} else {
			fmt.Fprintf(w, "%s\n", annotation)
		}
		return nil
	}

	// Write: join remaining args as the new annotation.
	newAnnotation := strings.Join(args, " ")
	set.SetAnnotation(newAnnotation)

	if err := st.Save(project, passphrase, set); err != nil {
		return err
	}

	fmt.Fprintf(w, "annotation updated for %q\n", project)
	return nil
}

func init() {
	_ = os.Stdout // ensure os import is used in non-test builds
}
