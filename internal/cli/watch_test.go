package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newWatchStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(filepath.Join(dir, "watch.db"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedWatch(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	set := env.NewSet()
	set.Put("WATCH_VAR", "hello")
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdWatchRunsCommand(t *testing.T) {
	st := newWatchStore(t)
	seedWatch(t, st, "proj", "pass")

	var buf bytes.Buffer
	err := CmdWatch(st, "proj", "pass", []string{"env"}, &buf)
	if err != nil {
		t.Fatalf("CmdWatch: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "[watch] started pid") {
		t.Errorf("expected start message, got: %s", output)
	}
}

func TestCmdWatchInjectsEnv(t *testing.T) {
	st := newWatchStore(t)
	seedWatch(t, st, "proj", "pass")

	// Use printenv to verify the variable is injected.
	var buf bytes.Buffer
	_ = os.Setenv("WATCH_VAR", "original")
	defer os.Unsetenv("WATCH_VAR")

	err := CmdWatch(st, "proj", "pass", []string{"printenv", "WATCH_VAR"}, &buf)
	if err != nil {
		t.Fatalf("CmdWatch: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected WATCH_VAR=hello in output, got: %s", buf.String())
	}
}

func TestCmdWatchEmptyArgs(t *testing.T) {
	st := newWatchStore(t)
	var buf bytes.Buffer
	err := CmdWatch(st, "proj", "pass", []string{}, &buf)
	if err == nil {
		t.Fatal("expected error for empty args")
	}
}

func TestCmdWatchProjectNotFound(t *testing.T) {
	st := newWatchStore(t)
	var buf bytes.Buffer
	err := CmdWatch(st, "missing", "pass", []string{"echo", "hi"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}
