package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdScatter writes each key=value from a project into individual files
// inside a target directory. Each file is named after the key and contains
// only the value, making it easy to consume secrets from file-based systems
// (e.g. Docker secrets, Kubernetes-style mounts).
func CmdScatter(st *store.Store, project, dir string, passphrase string, w io.Writer) error {
	if strings.TrimSpace(project) == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("target directory must not be empty")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := es.Keys()
	if len(keys) == 0 {
		fmt.Fprintf(w, "project %q has no keys — nothing written\n", project)
		return nil
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("cannot create directory %q: %w", dir, err)
	}

	sort.Strings(keys)
	for _, k := range keys {
		v, _ := es.Get(k)
		path := fmt.Sprintf("%s/%s", strings.TrimRight(dir, "/"), k)
		if err := os.WriteFile(path, []byte(v), 0o600); err != nil {
			return fmt.Errorf("failed to write %q: %w", path, err)
		}
		fmt.Fprintf(w, "wrote %s\n", path)
	}
	return nil
}

func init() {
	registerCommand("scatter", "scatter <project> <dir>  — write each key to a file inside <dir>")
}
