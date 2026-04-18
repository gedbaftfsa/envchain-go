package cli

import (
	"path/filepath"
	"testing"

	"github.com/envchain-go/internal/store"
)

// storeFromDir creates a *store.Store rooted at dir, failing the test on error.
func storeFromDir(t *testing.T, dir string) *store.Store {
	t.Helper()
	st, err := store.New(filepath.Join(dir, "envchain.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}
