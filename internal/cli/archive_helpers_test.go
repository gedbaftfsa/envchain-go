package cli

import (
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// mustStoreDir creates a store rooted at dir, failing the test on error.
func mustStoreDir(t *testing.T, dir string) *store.Store {
	t.Helper()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("mustStoreDir: %v", err)
	}
	return st
}

// passphraseReader returns a passphrase-reader func that always returns pass.
func passphraseReader(pass string) func(string) (string, error) {
	return func(_ string) (string, error) {
		return pass, nil
	}
}
