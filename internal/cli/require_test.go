package cli

import (
	"bytes"
	"testing"

	"github.com/user/envchain-go/internal/env"
	"github.com/user/envchain-go/internal/store"
)

func newRequireStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedRequire(t *testing.T, st *store.Store) {
	t.Helper()
	es := env.NewSet()
	es.Put("DATABASE_URL", "postgres://localhost/mydb")
	es.Put("SECRET_KEY", "s3cr3t")
	es.Put("EMPTY_VAL", "")
	if err := st.Save("myapp", "pass", es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdRequireAllPresent(t *testing.T) {
	st := newRequireStore(t)
	seedRequire(t, st)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "myapp", []string{"DATABASE_URL", "SECRET_KEY"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got == "" {
		t.Error("expected ok message, got empty output")
	}
}

func TestCmdRequireMissingKey(t *testing.T) {
	st := newRequireStore(t)
	seedRequire(t, st)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "myapp", []string{"DATABASE_URL", "MISSING_KEY"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if got := buf.String(); got == "" {
		t.Error("expected missing key output, got empty")
	}
}

func TestCmdRequireEmptyValue(t *testing.T) {
	st := newRequireStore(t)
	seedRequire(t, st)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "myapp", []string{"EMPTY_VAL"}, &buf)
	if err == nil {
		t.Fatal("expected error for empty value, got nil")
	}
}

func TestCmdRequireWrongPassphrase(t *testing.T) {
	st := newRequireStore(t)
	seedRequire(t, st)
	var buf bytes.Buffer
	err := CmdRequire(st, "wrong", "myapp", []string{"DATABASE_URL"}, &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestCmdRequireProjectNotFound(t *testing.T) {
	st := newRequireStore(t)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "ghost", []string{"KEY"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
}

func TestCmdRequireEmptyProject(t *testing.T) {
	st := newRequireStore(t)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "", []string{"KEY"}, &buf)
	if err == nil {
		t.Fatal("expected error for empty project name, got nil")
	}
}

func TestCmdRequireNoKeys(t *testing.T) {
	st := newRequireStore(t)
	seedRequire(t, st)
	var buf bytes.Buffer
	err := CmdRequire(st, "pass", "myapp", []string{}, &buf)
	if err == nil {
		t.Fatal("expected error when no keys specified, got nil")
	}
}
