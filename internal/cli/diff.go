package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

// CmdDiff prints the keys that differ between two projects.
func CmdDiff(st *store.Store, passphrase, projectA, projectB string, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	setA, err := st.Load(projectA, passphrase)
	if err != nil {
		return fmt.Errorf("project %q: %w", projectA, err)
	}

	setB, err := st.Load(projectB, passphrase)
	if err != nil {
		return fmt.Errorf("project %q: %w", projectB, err)
	}

	keysA := keySet(setA)
	keysB := keySet(setB)

	onlyA, onlyB, both := diffKeys(keysA, keysB)

	for _, k := range onlyA {
		fmt.Fprintf(out, "< %s\n", k)
	}
	for _, k := range onlyB {
		fmt.Fprintf(out, "> %s\n", k)
	}
	for _, k := range both {
		vA, _ := setA.Get(k)
		vB, _ := setB.Get(k)
		if vA != vB {
			fmt.Fprintf(out, "~ %s\n", k)
		}
	}
	return nil
}

func keySet(s *env.Set) map[string]struct{} {
	m := make(map[string]struct{})
	for _, k := range s.Keys() {
		m[k] = struct{}{}
	}
	return m
}

func diffKeys(a, b map[string]struct{}) (onlyA, onlyB, both []string) {
	for k := range a {
		if _, ok := b[k]; ok {
			both = append(both, k)
		} else {
			onlyA = append(onlyA, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			onlyB = append(onlyB, k)
		}
	}
	sort.Strings(onlyA)
	sort.Strings(onlyB)
	sort.Strings(both)
	return
}
