package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/store"
)

func newRecapStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedRecap(t *testing.T, st *store.Store) {
	t.Helper()
	set := newEnvSet(t, "ALPHA=hello", "BETA=world", "GAMMA=")
	if err := st.Save("myapp", set, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdRecapSuccess(t *testing.T) {
	st := newRecapStore(t)
	seedRecap(t, st)

	var buf bytes.Buffer
	if err := CmdRecap(st, "myapp", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "myapp") {
		t.Error("expected project name in output")
	}
	if !strings.Contains(out, "ALPHA") || !strings.Contains(out, "BETA") {
		t.Error("expected key names in output")
	}
	if strings.Contains(out, "hello") || strings.Contains(out, "world") {
		t.Error("values must not appear in recap output")
	}
	if !strings.Contains(out, "empty") {
		t.Error("expected 'empty' marker for GAMMA")
	}
}

func TestCmdRecapEmptyProject(t *testing.T) {
	st := newRecapStore(t)
	set := newEnvSet(t)
	if err := st.Save("empty", set, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var buf bytes.Buffer
	if err := CmdRecap(st, "empty", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no variables") {
		t.Error("expected 'no variables' message")
	}
}

func TestCmdRecapNotFound(t *testing.T) {
	st := newRecapStore(t)
	var buf bytes.Buffer
	err := CmdRecap(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdRecapWrongPassphrase(t *testing.T) {
	st := newRecapStore(t)
	seedRecap(t, st)
	var buf bytes.Buffer
	err := CmdRecap(st, "myapp", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdRecapEmptyName(t *testing.T) {
	st := newRecapStore(t)
	var buf bytes.Buffer
	err := CmdRecap(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
