package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdEnvCopy copies individual keys from one project to another.
// Usage: envchain env-copy <src-project> <dst-project> KEY [KEY...]
func CmdEnvCopy(st *store.Store, passphrase, srcProject, dstProject string, keys []string, overwrite bool, out io.Writer) error {
	if strings.TrimSpace(srcProject) == "" || strings.TrimSpace(dstProject) == "" {
		return fmt.Errorf("source and destination project names must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("at least one key must be specified")
	}

	srcSet, err := st.Load(srcProject, passphrase)
	if err != nil {
		return fmt.Errorf("load source %q: %w", srcProject, err)
	}

	dstSet, err := st.Load(dstProject, passphrase)
	if err != nil {
		dstSet = newEnvSetEmpty()
	}

	copied := 0
	skipped := 0
	for _, key := range keys {
		val, ok := srcSet.Get(key)
		if !ok {
			fmt.Fprintf(out, "warning: key %q not found in %q, skipping\n", key, srcProject)
			continue
		}
		_, exists := dstSet.Get(key)
		if exists && !overwrite {
			fmt.Fprintf(out, "skip: %q already exists in %q (use --overwrite to replace)\n", key, dstProject)
			skipped++
			continue
		}
		if err := dstSet.Put(key, val); err != nil {
			return fmt.Errorf("put key %q: %w", key, err)
		}
		copied++
	}

	if copied == 0 {
		fmt.Fprintf(out, "no keys copied (%d skipped)\n", skipped)
		return nil
	}

	if err := st.Save(dstProject, passphrase, dstSet); err != nil {
		return fmt.Errorf("save destination %q: %w", dstProject, err)
	}
	fmt.Fprintf(out, "copied %d key(s) from %q to %q (%d skipped)\n", copied, srcProject, dstProject, skipped)
	return nil
}
