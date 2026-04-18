package cli

import (
	"testing"

	"github.com/envchain-go/internal/env"
)

// newEnvSet is a test helper that returns a fresh env.Set.
func newEnvSet() (*env.Set, error) {
	return env.NewSet(), nil
}

// mustEnvSet panics if the set cannot be created (used in table-driven tests).
func mustEnvSet(t *testing.T) *env.Set {
	t.Helper()
	s := env.NewSet()
	if s == nil {
		t.Fatal("env.NewSet returned nil")
	}
	return s
}
