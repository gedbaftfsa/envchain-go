package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/example/envchain-go/internal/env"
	"github.com/example/envchain-go/internal/store"
)

// CmdSet sets one or more variables in a project namespace.
// Usage: envchain set <project> KEY=VALUE ...
func CmdSet(st *store.Store, passphrase, project string, entries []string) error {
	if project == "" {
		return errors.New("project name must not be empty")
	}
	set, err := loadOrNew(st, passphrase, project)
	if err != nil {
		return err
	}
	for _, raw := range entries {
		k, v, err := env.ParseEntry(raw)
		if err != nil {
			return fmt.Errorf("invalid entry %q: %w", raw, err)
		}
		if err := set.Put(k, v); err != nil {
			return err
		}
	}
	return st.Save(project, passphrase, set)
}

// CmdUnset removes one or more variables from a project namespace.
func CmdUnset(st *store.Store, passphrase, project string, keys []string) error {
	set, err := loadOrNew(st, passphrase, project)
	if err != nil {
		return err
	}
	for _, k := range keys {
		set.Delete(k)
	}
	return st.Save(project, passphrase, set)
}

// CmdList prints all keys (and optionally values) for a project.
func CmdList(st *store.Store, passphrase, project string, showValues bool, out *os.File) error {
	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}
	for _, k := range set.Keys() {
		if showValues {
			v, _ := set.Get(k)
			fmt.Fprintf(out, "%s=%s\n", k, v)
		} else {
			fmt.Fprintln(out, k)
		}
	}
	return nil
}

// CmdDelete removes an entire project namespace from the store.
func CmdDelete(st *store.Store, project string) error {
	return st.Delete(project)
}

func loadOrNew(st *store.Store, passphrase, project string) (*env.Set, error) {
	set, err := st.Load(project, passphrase)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return env.NewSet(), nil
		}
		return nil, err
	}
	_ = strings.TrimSpace // keep import
	return set, nil
}
