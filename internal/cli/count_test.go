package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newCountStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedCount(t *testing.T, st *store.Store, pass, project string, pairs ...string) {
	t.Helper()
	es := newEnvSet(t, pairs...)
	if err := st.Save(project, pass, es); err != nil {
		t.Fatalf("Save %q: %v", project, err)
	}
}

func TestCmdCountAll(t *testing.T) {
	st := newCountStore(t)
	pass := "s3cr3t"
	seedCount(t, st, pass, "alpha", "A=1", "B=2", "C=3")
	seedCount(t, st, pass, "beta", "X=10")

	var buf bytes.Buffer
	if err := CmdCount(st, pass, nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected count 3 for alpha, got: %s", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected 'beta' in output, got: %s", out)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected count 1 for beta, got: %s", out)
	}
}

func TestCmdCountSingleProject(t *testing.T) {
	st := newCountStore(t)
	pass := "s3cr3t"
	seedCount(t, st, pass, "myproject", "FOO=bar", "BAZ=qux")

	var buf bytes.Buffer
	if err := CmdCount(st, pass, []string{"myproject"}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "myproject") {
		t.Errorf("expected 'myproject' in output, got: %s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2, got: %s", out)
	}
}

func TestCmdCountProjectNotFound(t *testing.T) {
	st := newCountStore(t)
	pass := "s3cr3t"
	seedCount(t, st, pass, "exists", "K=V")

	var buf bytes.Buffer
	err := CmdCount(st, pass, []string{"ghost"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestCmdCountEmptyStore(t *testing.T) {
	st := newCountStore(t)

	var buf bytes.Buffer
	if err := CmdCount(st, "pass", nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "no projects") {
		t.Errorf("expected 'no projects' message, got: %s", buf.String())
	}
}
