package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newEnvListStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedEnvList(t *testing.T, st *store.Store) {
	t.Helper()
	set := env.NewSet()
	set.Put("DB_HOST", "localhost")
	set.Put("DB_PORT", "5432")
	set.Put("API_KEY", "secret")
	if err := st.Save("myproject", set, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdEnvListMissingAndPresent(t *testing.T) {
	st := newEnvListStore(t)
	seedEnvList(t, st)

	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "9999") // differs from stored value
	// API_KEY intentionally not set

	var buf bytes.Buffer
	if err := CmdEnvList(st, "myproject", "pass", &buf); err != nil {
		t.Fatalf("CmdEnvList: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "match") {
		t.Error("expected at least one 'match' status")
	}
	if !strings.Contains(out, "differs") {
		t.Error("expected at least one 'differs' status")
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected at least one 'missing' status")
	}
}

func TestCmdEnvListEmptyProject(t *testing.T) {
	st := newEnvListStore(t)
	set := env.NewSet()
	if err := st.Save("empty", set, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var buf bytes.Buffer
	if err := CmdEnvList(st, "empty", "pass", &buf); err != nil {
		t.Fatalf("CmdEnvList: %v", err)
	}
	if !strings.Contains(buf.String(), "no variables") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestCmdEnvListNotFound(t *testing.T) {
	st := newEnvListStore(t)
	var buf bytes.Buffer
	err := CmdEnvList(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdEnvListEmptyName(t *testing.T) {
	st := newEnvListStore(t)
	var buf bytes.Buffer
	err := CmdEnvList(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
