package cli

import (
	"bytes"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newCheckStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedCheck(t *testing.T, st *store.Store, project, pass string, entries map[string]string) {
	t.Helper()
	set := env.NewSet()
	for k, v := range entries {
		if err := set.Put(k, v); err != nil {
			t.Fatalf("set.Put: %v", err)
		}
	}
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestCmdCheckClean(t *testing.T) {
	st := newCheckStore(t)
	seedCheck(t, st, "myapp", "pass", map[string]string{"FOO": "bar", "BAZ": "qux"})
	var buf bytes.Buffer
	if err := CmdCheck(st, "pass", "myapp", nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("passed")) {
		t.Errorf("expected 'passed' in output, got: %s", buf.String())
	}
}

func TestCmdCheckEmptyValue(t *testing.T) {
	st := newCheckStore(t)
	seedCheck(t, st, "myapp", "pass", map[string]string{"FOO": "", "BAZ": "qux"})
	var buf bytes.Buffer
	err := CmdCheck(st, "pass", "myapp", nil, &buf)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if !bytes.Contains(buf.Bytes(), []byte("empty value")) {
		t.Errorf("expected 'empty value' in output, got: %s", buf.String())
	}
}

func TestCmdCheckMissingRequired(t *testing.T) {
	st := newCheckStore(t)
	seedCheck(t, st, "myapp", "pass", map[string]string{"FOO": "bar"})
	var buf bytes.Buffer
	err := CmdCheck(st, "pass", "myapp", []string{"FOO", "MISSING_KEY"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	if !bytes.Contains(buf.Bytes(), []byte("MISSING_KEY")) {
		t.Errorf("expected MISSING_KEY in output, got: %s", buf.String())
	}
}

func TestCmdCheckNotFound(t *testing.T) {
	st := newCheckStore(t)
	var buf bytes.Buffer
	if err := CmdCheck(st, "pass", "ghost", nil, &buf); err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdCheckEmptyProject(t *testing.T) {
	st := newCheckStore(t)
	var buf bytes.Buffer
	if err := CmdCheck(st, "pass", "", nil, &buf); err == nil {
		t.Fatal("expected error for empty project name")
	}
}
