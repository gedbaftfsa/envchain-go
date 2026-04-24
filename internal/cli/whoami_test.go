package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envchain-go/internal/env"
	"github.com/your-org/envchain-go/internal/store"
)

func newWhoamiStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func seedWhoami(t *testing.T, s *store.Store, project, pass string, pairs ...string) {
	t.Helper()
	set := env.NewSet()
	for i := 0; i+1 < len(pairs); i += 2 {
		if err := set.Put(pairs[i], pairs[i+1]); err != nil {
			t.Fatalf("set.Put: %v", err)
		}
	}
	if err := s.Save(project, set, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestCmdWhoamiFound(t *testing.T) {
	s := newWhoamiStore(t)
	seedWhoami(t, s, "alpha", "pass", "DB_URL", "postgres://localhost/alpha")
	seedWhoami(t, s, "beta", "pass", "DB_URL", "postgres://localhost/beta", "API_KEY", "secret")
	seedWhoami(t, s, "gamma", "pass", "API_KEY", "other")

	var buf bytes.Buffer
	if err := CmdWhoami(s, "pass", "DB_URL", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output, got: %s", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected beta in output, got: %s", out)
	}
	if strings.Contains(out, "gamma") {
		t.Errorf("gamma should not appear for DB_URL, got: %s", out)
	}
}

func TestCmdWhoamiNotFound(t *testing.T) {
	s := newWhoamiStore(t)
	seedWhoami(t, s, "alpha", "pass", "FOO", "bar")

	var buf bytes.Buffer
	if err := CmdWhoami(s, "pass", "MISSING_KEY", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "not found") {
		t.Errorf("expected 'not found' message, got: %s", out)
	}
}

func TestCmdWhoamiEmptyKey(t *testing.T) {
	s := newWhoamiStore(t)
	var buf bytes.Buffer
	err := CmdWhoami(s, "pass", "", &buf)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCmdWhoamiWrongPassphrase(t *testing.T) {
	s := newWhoamiStore(t)
	seedWhoami(t, s, "alpha", "correct", "KEY", "val")

	var buf bytes.Buffer
	err := CmdWhoami(s, "wrong", "KEY", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
