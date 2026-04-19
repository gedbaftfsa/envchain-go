package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newTouchStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedTouch(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	set := env.NewSet()
	set.Put("KEY", "value")
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdTouchSuccess(t *testing.T) {
	st := newTouchStore(t)
	seedTouch(t, st, "myproject", "pass")

	var buf bytes.Buffer
	if err := CmdTouch(st, "myproject", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "touched") {
		t.Errorf("expected 'touched' in output, got: %s", buf.String())
	}
}

func TestCmdTouchWrongPassphrase(t *testing.T) {
	st := newTouchStore(t)
	seedTouch(t, st, "myproject", "correct")

	var buf bytes.Buffer
	err := CmdTouch(st, "myproject", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdTouchNotFound(t *testing.T) {
	st := newTouchStore(t)

	var buf bytes.Buffer
	err := CmdTouch(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdTouchEmptyProject(t *testing.T) {
	st := newTouchStore(t)

	var buf bytes.Buffer
	err := CmdTouch(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
