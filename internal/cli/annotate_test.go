package cli

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newAnnotateStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedAnnotate(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	set := env.NewSet()
	_ = set.Put("KEY", "val")
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("seedAnnotate Save: %v", err)
	}
}

func TestCmdAnnotateNoAnnotation(t *testing.T) {
	st := newAnnotateStore(t)
	seedAnnotate(t, st, "proj", "secret")

	var buf bytes.Buffer
	if err := CmdAnnotate(st, "secret", "proj", nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "(no annotation)\n" {
		t.Errorf("expected '(no annotation)', got %q", got)
	}
}

func TestCmdAnnotateSetAndGet(t *testing.T) {
	st := newAnnotateStore(t)
	seedAnnotate(t, st, "proj", "secret")

	var buf bytes.Buffer
	if err := CmdAnnotate(st, "secret", "proj", []string{"Production", "creds"}, &buf); err != nil {
		t.Fatalf("set annotation: %v", err)
	}

	buf.Reset()
	if err := CmdAnnotate(st, "secret", "proj", nil, &buf); err != nil {
		t.Fatalf("get annotation: %v", err)
	}
	if got := buf.String(); got != "Production creds\n" {
		t.Errorf("expected 'Production creds', got %q", got)
	}
}

func TestCmdAnnotateWrongPassphrase(t *testing.T) {
	st := newAnnotateStore(t)
	seedAnnotate(t, st, "proj", "secret")

	var buf bytes.Buffer
	err := CmdAnnotate(st, "wrong", "proj", nil, &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestCmdAnnotateEmptyProject(t *testing.T) {
	st := newAnnotateStore(t)

	var buf bytes.Buffer
	err := CmdAnnotate(st, "secret", "", nil, &buf)
	if err == nil {
		t.Fatal("expected error for empty project, got nil")
	}
}

func TestCmdAnnotateNotFound(t *testing.T) {
	st := newAnnotateStore(t)

	var buf bytes.Buffer
	err := CmdAnnotate(st, "secret", "ghost", nil, &buf)
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
}
