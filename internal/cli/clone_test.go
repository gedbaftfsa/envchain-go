package cli

import (
	"testing"

	"github.com/envchain-go/internal/store"
)

func seedClone(t *testing.T, s store.Store, project, pass string) {
	t.Helper()
	if err := CmdSet(s, project, pass, "FOO=bar"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(s, project, pass, "BAZ=qux"); err != nil {
		t.Fatal(err)
	}
}

func TestCmdCloneSameStore(t *testing.T) {
	s := newTempStore(t)
	const pass = "passphrase"
	seedClone(t, s, "src", pass)

	if err := CmdCloneSameStore(s, "src", "dst", pass); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	set, err := s.Load("dst", pass)
	if err != nil {
		t.Fatalf("load dst: %v", err)
	}

	v, ok := set.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q ok=%v", v, ok)
	}
	v, ok = set.Get("BAZ")
	if !ok || v != "qux" {
		t.Errorf("expected BAZ=qux, got %q ok=%v", v, ok)
	}
}

func TestCmdCloneSrcNotFound(t *testing.T) {
	s := newTempStore(t)
	err := CmdCloneSameStore(s, "ghost", "dst", "pass")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCmdCloneEmptyNames(t *testing.T) {
	s := newTempStore(t)
	if err := CmdClone(s, "", "pass", s, "dst", "pass"); err == nil {
		t.Error("expected error for empty src project")
	}
	if err := CmdClone(s, "src", "pass", s, "", "pass"); err == nil {
		t.Error("expected error for empty dst project")
	}
}

func TestCmdCloneCrossStore(t *testing.T) {
	src := newTempStore(t)
	dst := newTempStore(t)
	const pass = "secret"
	seedClone(t, src, "proj", pass)

	if err := CmdClone(src, "proj", pass, dst, "proj-copy", pass); err != nil {
		t.Fatalf("cross-store clone: %v", err)
	}

	set, err := dst.Load("proj-copy", pass)
	if err != nil {
		t.Fatalf("load from dst: %v", err)
	}
	if _, ok := set.Get("FOO"); !ok {
		t.Error("FOO missing in cloned project")
	}
}
