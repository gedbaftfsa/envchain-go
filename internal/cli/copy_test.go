package cli

import (
	"testing"

	"github.com/user/envchain-go/internal/env"
)

func seedProject2(t *testing.T, st interface{ Save(string, string, *env.Set) error }, project, pass string, kvs map[string]string) {
	t.Helper()
	set := env.NewSet()
	for k, v := range kvs {
		set.Put(k, v)
	}
	if err := st.Save(project, pass, set); err != nil {
		t.Fatal(err)
	}
}

func TestCmdCopyOverwrite(t *testing.T) {
	st := newTempStore(t)
	seedProject2(t, st, "src", "p", map[string]string{"A": "1", "B": "2"})
	seedProject2(t, st, "dst", "p", map[string]string{"A": "old", "C": "3"})

	if err := CmdCopy(st, "src", "p", "dst", "p", true); err != nil {
		t.Fatal(err)
	}

	loaded, _ := st.Load("dst", "p")
	if v, _ := loaded.Get("A"); v != "1" {
		t.Errorf("A = %q, want 1", v)
	}
	if v, _ := loaded.Get("B"); v != "2" {
		t.Errorf("B = %q, want 2", v)
	}
	if v, _ := loaded.Get("C"); v != "3" {
		t.Errorf("C = %q, want 3", v)
	}
}

func TestCmdCopyNoOverwrite(t *testing.T) {
	st := newTempStore(t)
	seedProject2(t, st, "src", "p", map[string]string{"A": "new", "B": "2"})
	seedProject2(t, st, "dst", "p", map[string]string{"A": "old"})

	if err := CmdCopy(st, "src", "p", "dst", "p", false); err != nil {
		t.Fatal(err)
	}

	loaded, _ := st.Load("dst", "p")
	if v, _ := loaded.Get("A"); v != "old" {
		t.Errorf("A = %q, want old (no overwrite)", v)
	}
	if v, _ := loaded.Get("B"); v != "2" {
		t.Errorf("B = %q, want 2", v)
	}
}

func TestCmdCopySrcNotFound(t *testing.T) {
	st := newTempStore(t)
	err := CmdCopy(st, "missing", "p", "dst", "p", true)
	if err == nil {
		t.Fatal("expected error for missing source project")
	}
}
