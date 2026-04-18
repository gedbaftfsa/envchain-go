package cli

import (
	"fmt"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

// CmdClone copies a project from one store path to another store (or the same
// store under a new name), optionally re-encrypting with a different passphrase.
func CmdClone(
	src store.Store,
	srcProject string,
	srcPass string,
	dst store.Store,
	dstProject string,
	dstPass string,
) error {
	if srcProject == "" || dstProject == "" {
		return fmt.Errorf("source and destination project names must not be empty")
	}

	set, err := src.Load(srcProject, srcPass)
	if err != nil {
		return fmt.Errorf("clone: load source %q: %w", srcProject, err)
	}

	if err := dst.Save(dstProject, dstPass, set); err != nil {
		return fmt.Errorf("clone: save destination %q: %w", dstProject, err)
	}

	fmt.Printf("Cloned %q → %q\n", srcProject, dstProject)
	return nil
}

// CmdCloneSameStore is a convenience wrapper when src and dst are the same store.
func CmdCloneSameStore(
	s store.Store,
	srcProject, dstProject, passphrase string,
) error {
	set, err := s.Load(srcProject, passphrase)
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	// Build a fresh set so we don't share internal map references.
	newSet := env.NewSet()
	for _, k := range set.Keys() {
		v, _ := set.Get(k)
		_ = newSet.Put(k, v)
	}

	if err := s.Save(dstProject, passphrase, newSet); err != nil {
		return fmt.Errorf("clone: save %q: %w", dstProject, err)
	}

	fmt.Printf("Cloned %q → %q\n", srcProject, dstProject)
	return nil
}
