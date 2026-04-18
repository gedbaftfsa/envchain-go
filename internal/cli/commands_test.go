package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/envchain-go/internal/store"
)

const testPass = "hunter2"

func newTempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(filepath.Join(dir, "envchain"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestCmdSet(t *testing.T) {
	st := newTempStore(t)
	if err := CmdSet(st, testPass, "myapp", []string{"FOO=bar", "BAZ=qux"}); err != nil {
		t.Fatalf("CmdSet: %v", err)
	}
	set, err := st.Load("myapp", testPass)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if v, _ := set.Get("FOO"); v != "bar" {
		t.Errorf("FOO = %q, want bar", v)
	}
	if v, _ := set.Get("BAZ"); v != "qux" {
		t.Errorf("BAZ = %q, want qux", v)
	}
}

func TestCmdSetInvalidEntry(t *testing.T) {
	st := newTempStore(t)
	if err := CmdSet(st, testPass, "myapp", []string{"NOEQUALS"}); err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestCmdUnset(t *testing.T) {
	st := newTempStore(t)
	_ = CmdSet(st, testPass, "proj", []string{"A=1", "B=2"})
	if err := CmdUnset(st, testPass, "proj", []string{"A"}); err != nil {
		t.Fatalf("CmdUnset: %v", err)
	}
	set, _ := st.Load("proj", testPass)
	if _, ok := set.Get("A"); ok {
		t.Error("A should have been removed")
	}
	if v, _ := set.Get("B"); v != "2" {
		t.Errorf("B = %q, want 2", v)
	}
}

func TestCmdList(t *testing.T) {
	st := newTempStore(t)
	_ = CmdSet(st, testPass, "proj", []string{"X=hello"})

	r, w, _ := os.Pipe()
	err := CmdList(st, testPass, "proj", false, w)
	w.Close()
	if err != nil {
		t.Fatalf("CmdList: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.String() != "X\n" {
		t.Errorf("output = %q, want \"X\\n\"", buf.String())
	}
}

func TestCmdDelete(t *testing.T) {
	st := newTempStore(t)
	_ = CmdSet(st, testPass, "proj", []string{"K=V"})
	if err := CmdDelete(st, "proj"); err != nil {
		t.Fatalf("CmdDelete: %v", err)
	}
	_, err := st.Load("proj", testPass)
	if err == nil {
		t.Error("expected error after deletion")
	}
}
