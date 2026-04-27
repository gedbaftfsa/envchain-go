package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newRestoreStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func writeBackupFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "backup.env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write backup file: %v", err)
	}
	return p
}

func TestCmdRestoreSuccess(t *testing.T) {
	st := newRestoreStore(t)
	content := "FOO=bar\nBAZ=qux\n"
	file := writeBackupFile(t, content)

	var buf bytes.Buffer
	err := CmdRestore(st, "myproject", file, func() (string, error) { return "secret", nil }, &buf)
	if err != nil {
		t.Fatalf("CmdRestore: %v", err)
	}

	es, err := st.Load("myproject", "secret")
	if err != nil {
		t.Fatalf("Load after restore: %v", err)
	}
	if v, _ := es.Get("FOO"); v != "bar" {
		t.Errorf("expected FOO=bar, got %q", v)
	}
	if v, _ := es.Get("BAZ"); v != "qux" {
		t.Errorf("expected BAZ=qux, got %q", v)
	}
}

func TestCmdRestoreEmptyProject(t *testing.T) {
	st := newRestoreStore(t)
	file := writeBackupFile(t, "A=1\n")
	var buf bytes.Buffer
	err := CmdRestore(st, "", file, func() (string, error) { return "x", nil }, &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdRestoreMissingFile(t *testing.T) {
	st := newRestoreStore(t)
	var buf bytes.Buffer
	err := CmdRestore(st, "proj", "/nonexistent/path.env", func() (string, error) { return "x", nil }, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCmdRestoreEmptyFilePath(t *testing.T) {
	st := newRestoreStore(t)
	var buf bytes.Buffer
	err := CmdRestore(st, "proj", "", func() (string, error) { return "x", nil }, &buf)
	if err == nil {
		t.Fatal("expected error for empty file path")
	}
}
