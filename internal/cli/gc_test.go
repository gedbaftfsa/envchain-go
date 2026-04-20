package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newGCStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedGC(t *testing.T, st *store.Store, project, passphrase string) {
	t.Helper()
	es := newEnvSet(t, "KEY=val")
	if err := st.Save(project, passphrase, es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdGCNothingToDo(t *testing.T) {
	st := newGCStore(t)
	seedGC(t, st, "myproject", "pass")

	// create a snapshot that belongs to the live project
	var buf bytes.Buffer
	if err := CmdSnapshot(st, "myproject", "pass", "snap1", &buf); err != nil {
		t.Fatalf("CmdSnapshot: %v", err)
	}

	buf.Reset()
	if err := CmdGC(st, "", &buf); err != nil {
		t.Fatalf("CmdGC: %v", err)
	}
	out := buf.String()
	if out != "gc: nothing to collect\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestCmdGCRemovesOrphans(t *testing.T) {
	st := newGCStore(t)
	seedGC(t, st, "alive", "pass")

	// snapshot for a project we will then delete
	var buf bytes.Buffer
	if err := CmdSnapshot(st, "alive", "pass", "before-delete", &buf); err != nil {
		t.Fatalf("CmdSnapshot: %v", err)
	}

	// delete the project so the snapshot becomes orphaned
	if err := st.Delete("alive"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	buf.Reset()
	if err := CmdGC(st, "", &buf); err != nil {
		t.Fatalf("CmdGC: %v", err)
	}
	out := buf.String()
	expected := fmt.Sprintf("removed orphaned snapshot: alive@before-delete\ngc: removed 1 orphaned snapshot(s)\n")
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestCmdGCEmptyStore(t *testing.T) {
	st := newGCStore(t)
	var buf bytes.Buffer
	if err := CmdGC(st, "", &buf); err != nil {
		t.Fatalf("CmdGC: %v", err)
	}
	if buf.String() != "gc: nothing to collect\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
