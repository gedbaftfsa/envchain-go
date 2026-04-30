package cli

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newEnvCopyStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedEnvCopy(t *testing.T, st *store.Store, project, pass string, pairs map[string]string) {
	t.Helper()
	set := env.NewSet()
	for k, v := range pairs {
		if err := set.Put(k, v); err != nil {
			t.Fatalf("put %s: %v", k, err)
		}
	}
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("save %s: %v", project, err)
	}
}

func TestCmdEnvCopyCopiesKeys(t *testing.T) {
	st := newEnvCopyStore(t)
	pass := "secret"
	seedEnvCopy(t, st, "src", pass, map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "abc"})

	var buf bytes.Buffer
	if err := CmdEnvCopy(st, pass, "src", "dst", []string{"DB_HOST", "DB_PORT"}, false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dstSet, err := st.Load("dst", pass)
	if err != nil {
		t.Fatalf("load dst: %v", err)
	}
	if v, ok := dstSet.Get("DB_HOST"); !ok || v != "localhost" {
		t.Errorf("DB_HOST: got %q, want %q", v, "localhost")
	}
	if v, ok := dstSet.Get("DB_PORT"); !ok || v != "5432" {
		t.Errorf("DB_PORT: got %q, want %q", v, "5432")
	}
	if _, ok := dstSet.Get("API_KEY"); ok {
		t.Error("API_KEY should not have been copied")
	}
}

func TestCmdEnvCopyNoOverwrite(t *testing.T) {
	st := newEnvCopyStore(t)
	pass := "secret"
	seedEnvCopy(t, st, "src", pass, map[string]string{"KEY": "new-value"})
	seedEnvCopy(t, st, "dst", pass, map[string]string{"KEY": "old-value"})

	var buf bytes.Buffer
	if err := CmdEnvCopy(st, pass, "src", "dst", []string{"KEY"}, false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dstSet, _ := st.Load("dst", pass)
	if v, _ := dstSet.Get("KEY"); v != "old-value" {
		t.Errorf("expected old-value preserved, got %q", v)
	}
}

func TestCmdEnvCopyOverwrite(t *testing.T) {
	st := newEnvCopyStore(t)
	pass := "secret"
	seedEnvCopy(t, st, "src", pass, map[string]string{"KEY": "new-value"})
	seedEnvCopy(t, st, "dst", pass, map[string]string{"KEY": "old-value"})

	var buf bytes.Buffer
	if err := CmdEnvCopy(st, pass, "src", "dst", []string{"KEY"}, true, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dstSet, _ := st.Load("dst", pass)
	if v, _ := dstSet.Get("KEY"); v != "new-value" {
		t.Errorf("expected new-value, got %q", v)
	}
}

func TestCmdEnvCopyMissingKey(t *testing.T) {
	st := newEnvCopyStore(t)
	pass := "secret"
	seedEnvCopy(t, st, "src", pass, map[string]string{"REAL": "val"})

	var buf bytes.Buffer
	if err := CmdEnvCopy(st, pass, "src", "dst", []string{"REAL", "GHOST"}, false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("GHOST")) {
		t.Errorf("expected warning about missing key GHOST, got: %s", out)
	}
}

func TestCmdEnvCopyEmptyKeys(t *testing.T) {
	st := newEnvCopyStore(t)
	var buf bytes.Buffer
	err := CmdEnvCopy(st, "pass", "src", "dst", []string{}, false, &buf)
	if err == nil {
		t.Error("expected error for empty keys list")
	}
}

func TestCmdEnvCopyEmptyProject(t *testing.T) {
	st := newEnvCopyStore(t)
	var buf bytes.Buffer
	err := CmdEnvCopy(st, "pass", "", "dst", []string{"KEY"}, false, &buf)
	if err == nil {
		t.Error("expected error for empty source project")
	}
}
