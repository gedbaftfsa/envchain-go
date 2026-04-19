package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/envchain/envchain-go/internal/env"
	"github.com/envchain/envchain-go/internal/store"
)

// CmdMerge merges variables from srcProject into dstProject.
// By default existing keys in dst are not overwritten unless --overwrite is set.
func CmdMerge(st *store.Store, srcProject, dstProject, passphrase string, overwrite bool, w io.Writer) error {
	if srcProject == "" || dstProject == "" {
		return fmt.Errorf("source and destination project names are required")
	}
	if srcProject == dstProject {
		return fmt.Errorf("source and destination project must differ")
	}

	srcSet, err := st.Load(srcProject, passphrase)
	if err != nil {
		return fmt.Errorf("load source %q: %w", srcProject, err)
	}

	dstSet, err := loadOrNew(st, dstProject, passphrase)
	if err != nil {
		return fmt.Errorf("load destination %q: %w", dstProject, err)
	}

	mode := env.MergeSkip
	if overwrite {
		mode = env.MergeOverwrite
	}

	added, skipped := env.Merge(dstSet, srcSet, mode)

	if err := st.Save(dstProject, passphrase, dstSet); err != nil {
		return fmt.Errorf("save destination %q: %w", dstProject, err)
	}

	fmt.Fprintf(w, "merged %q → %q: %d added, %d skipped\n", srcProject, dstProject, added, skipped)
	return nil
}

func init() {
	_ = os.Stderr // ensure os import used
}
