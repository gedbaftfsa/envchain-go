package cli

import (
	"bytes"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newRedactStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedRedact(t *testing.T, st *store.Store) {
	t.Helper()
	es := env.NewSet()
	_ = es.Put("API_KEY", "supersecret")
	_ = es.Put("DB_PASS", "hunter2")
	if err := st.Save("myapp", "pass", es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdRedactReplacesSecrets(t *testing.T) {
	st := newRedactStore(t)
	seedRedact(t, st)

	input := "connecting with key=supersecret and pass=hunter2 done"
	var buf bytes.Buffer
	if err := CmdRedact(st, "myapp", "pass", input, &buf); err != nil {
		t.Fatalf("CmdRedact: %v", err)
	}

	got := buf.String()
	if contains(got, "supersecret") || contains(got, "hunter2") {
		t.Errorf("secrets not redacted: %q", got)
	}
	if !contains(got, "***REDACTED***") {
		t.Errorf("expected placeholder in output: %q", got)
	}
}

func TestCmdRedactNoSecrets(t *testing.T) {
	st := newRedactStore(t)
	seedRedact(t, st)

	input := "nothing sensitive here"
	var buf bytes.Buffer
	if err := CmdRedact(st, "myapp", "pass", input, &buf); err != nil {
		t.Fatalf("CmdRedact: %v", err)
	}
	if buf.String() != input {
		t.Errorf("expected unchanged output, got %q", buf.String())
	}
}

func TestCmdRedactWrongPassphrase(t *testing.T) {
	st := newRedactStore(t)
	seedRedact(t, st)

	var buf bytes.Buffer
	err := CmdRedact(st, "myapp", "wrong", "text", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdRedactEmptyProject(t *testing.T) {
	st := newRedactStore(t)
	var buf bytes.Buffer
	err := CmdRedact(st, "", "pass", "text", &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdRedactNotFound(t *testing.T) {
	st := newRedactStore(t)
	var buf bytes.Buffer
	err := CmdRedact(st, "ghost", "pass", "text", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

// contains is a small helper to avoid importing strings in test files.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
