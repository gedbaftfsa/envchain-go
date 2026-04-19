package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/envchain-go/internal/store"
)

// CmdDoctor checks the health of the envchain store and environment.
func CmdDoctor(st *store.Store, out io.Writer) error {
	ok := true

	// Check store directory exists
	dir := st.Dir()
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		fmt.Fprintf(out, "[FAIL] store directory not found: %s\n", dir)
		ok = false
	} else {
		fmt.Fprintf(out, "[ OK ] store directory: %s\n", dir)
	}

	// Check directory permissions
	if info, err := os.Stat(dir); err == nil {
		mode := info.Mode().Perm()
		if mode&0o077 != 0 {
			fmt.Fprintf(out, "[WARN] store directory is world/group readable (mode %o)\n", mode)
			ok = false
		} else {
			fmt.Fprintf(out, "[ OK ] store directory permissions: %o\n", mode)
		}
	}

	// Count projects
	entries, err := filepath.Glob(filepath.Join(dir, "*.bin"))
	if err != nil {
		entries = nil
	}
	fmt.Fprintf(out, "[ OK ] projects found: %d\n", len(entries))

	if !ok {
		return fmt.Errorf("doctor found issues with your envchain setup")
	}
	return nil
}
