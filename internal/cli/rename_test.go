package cli

import (
	"testing"

	"github.com/envchain/envchain-go/internal/env"
)

func TestCmdRenameSuccess(t *testing.T) {
	st := newTempStore(t)
	const pass = "secret"

	seedProject2(t, st, "alpha", pass)

	if err := CmdRename(st, "alpha", "beta", pass); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// old name should be gone
	if _, err := st.Load("alpha", pass); err == nil {
		t.Fatal("expected old project to be deleted")
	}

	// new name should exist with same data
	set, err := st.Load("beta", pass)
	if err != nil {
		t.Fatalf("could not load renamed project: %v", err)
	}
	if v, _ := set.Get("KEY"); v != "val" {
		t.Fatalf("expected KEY=val, got %q", v)
	}
}

func TestCmdRenameSameName(t *testing.T) {
	st := newTempStore(t)
	seedProject2(t, st, "alpha", "secret")

	err := CmdRename(st, "alpha", "alpha", "secret")
	if err == nil {
		t.Fatal("expected error for identical names")
	}
}

func TestCmdRenameDestExists(t *testing.T) {
	st := newTempStore(t)
	const pass = "secret"
	seedProject2(t, st, "alpha", pass)
	seedProject2(t, st, "beta", pass)

	err := CmdRename(st, "alpha", "beta", pass)
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestCmdRenameSrcNotFound(t *testing.T) {
	st := newTempStore(t)

	err := CmdRename(st, "ghost", "new", "secret")
	if err == nil {
		t.Fatal("expected error for missing source project")
	}
}

func TestCmdRenameEmptyArgs(t *testing.T) {
	st := newTempStore(t)

	if err := CmdRename(st, "", "new", "secret"); err == nil {
		t.Fatal("expected error for empty old name")
	}
	if err := CmdRename(st, "old", "", "secret"); err == nil {
		t.Fatal("expected error for empty new name")
	}
}

// helper shared with copy_test.go already defines seedProject2; keep local if needed.
func init() {
	// ensure env package is used
	_ = env.NewSet
}
