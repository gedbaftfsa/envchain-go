package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newArchiveStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedArchive(t *testing.T, st *store.Store, pass string) {
	t.Helper()
	for _, proj := range []string{"alpha", "beta"} {
		env := newEnvSet(t, "FOO=bar", "BAZ=qux")
		if err := st.Save(proj, pass, env); err != nil {
			t.Fatalf("save %s: %v", proj, err)
		}
	}
}

func TestCmdArchiveCreatesFiles(t *testing.T) {
	st := newArchiveStore(t)
	seedArchive(t, st, "secret")
	dest := filepath.Join(t.TempDir(), "out")
	var buf bytes.Buffer
	if err := CmdArchive(st, "secret", dest, &buf); err != nil {
		t.Fatalf("CmdArchive: %v", err)
	}
	for _, proj := range []string{"alpha", "beta"} {
		p := filepath.Join(dest, proj+".env")
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("missing archive file %s: %v", p, err)
		}
		if !strings.Contains(string(data), "FOO=bar") {
			t.Errorf("%s: expected FOO=bar in content", proj)
		}
	}
	if !strings.Contains(buf.String(), "archived alpha") {
		t.Errorf("expected output to mention alpha")
	}
}

func TestCmdArchiveWrongPassphrase(t *testing.T) {
	st := newArchiveStore(t)
	seedArchive(t, st, "secret")
	dest := t.TempDir()
	var buf bytes.Buffer
	if err := CmdArchive(st, "wrong", dest, &buf); err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdRestoreArchive(t *testing.T) {
	src := newArchiveStore(t)
	seedArchive(t, src, "pass")
	dest := filepath.Join(t.TempDir(), "arch")
	var buf bytes.Buffer
	if err := CmdArchive(src, "pass", dest, &buf); err != nil {
		t.Fatalf("archive: %v", err)
	}

	dst := newArchiveStore(t)
	var buf2 bytes.Buffer
	if err := CmdRestoreArchive(dst, "newpass", dest, &buf2); err != nil {
		t.Fatalf("restore-archive: %v", err)
	}
	env, err := dst.Load("alpha", "newpass")
	if err != nil {
		t.Fatalf("load after restore: %v", err)
	}
	v, ok := env.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q ok=%v", v, ok)
	}
}

func TestCmdArchiveEmpty(t *testing.T) {
	st := newArchiveStore(t)
	dest := t.TempDir()
	var buf bytes.Buffer
	if err := CmdArchive(st, "pass", dest, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no projects") {
		t.Errorf("expected 'no projects' message")
	}
}
