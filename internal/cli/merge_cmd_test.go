package cli

import (
	"bytes"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
)

func seedMerge(t *testing.T, st interface{ Save(string, string, *env.Set) error }, project, pass string, pairs ...string) {
	t.Helper()
	s := mustEnvSet(t, pairs...)
	if err := st.Save(project, pass, s); err != nil {
		t.Fatalf("seedMerge: %v", err)
	}
}

func TestCmdMergeSkip(t *testing.T) {
	st := newTempStore(t)
	seedMerge(t, st, "src", "pass", "A=1", "B=2")
	seedMerge(t, st, "dst", "pass", "A=old")

	var buf bytes.Buffer
	if err := CmdMerge(st, "src", "dst", "pass", false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, _ := st.Load("dst", "pass")
	v, _ := dst.Get("A")
	if v != "old" {
		t.Fatalf("expected A=old (skip mode), got %s", v)
	}
	v2, ok := dst.Get("B")
	if !ok || v2 != "2" {
		t.Fatalf("expected B=2 to be added")
	}
}

func TestCmdMergeOverwrite(t *testing.T) {
	st := newTempStore(t)
	seedMerge(t, st, "src", "pass", "A=new")
	seedMerge(t, st, "dst", "pass", "A=old")

	var buf bytes.Buffer
	if err := CmdMerge(st, "src", "dst", "pass", true, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, _ := st.Load("dst", "pass")
	v, _ := dst.Get("A")
	if v != "new" {
		t.Fatalf("expected A=new (overwrite mode), got %s", v)
	}
}

func TestCmdMergeSameName(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdMerge(st, "x", "x", "pass", false, &buf)
	if err == nil {
		t.Fatal("expected error for same src/dst")
	}
}

func TestCmdMergeSrcNotFound(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdMerge(st, "missing", "dst", "pass", false, &buf)
	if err == nil {
		t.Fatal("expected error for missing src")
	}
}
