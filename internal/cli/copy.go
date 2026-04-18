package cli

import (
	"fmt"

	"github.com/user/envchain-go/internal/env"
	"github.com/user/envchain-go/internal/store"
)

// CmdCopy copies all variables from srcProject into dstProject.
// Existing keys in dstProject are overwritten; extra keys are preserved.
// If overwrite is false, existing keys in dstProject are kept.
func CmdCopy(st *store.Store, srcProject, srcPass, dstProject, dstPass string, overwrite bool) error {
	src, err := st.Load(srcProject, srcPass)
	if err != nil {
		return fmt.Errorf("copy: load source %q: %w", srcProject, err)
	}

	dst, err := loadOrNew(st, dstProject, dstPass)
	if err != nil {
		return fmt.Errorf("copy: load dest %q: %w", dstProject, err)
	}

	merged := env.Merge(dst, src, overwrite)

	if err := st.Save(dstProject, dstPass, merged); err != nil {
		return fmt.Errorf("copy: save dest %q: %w", dstProject, err)
	}
	return nil
}
