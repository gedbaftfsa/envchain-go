package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newShareStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedShare(t *testing.T, st *store.Store) {
	t.Helper()
	set := env.NewSet()
	set.Put("FOO", "bar")
	set.Put("BAZ", "qux")
	if err := st.Save("shareproj", "pass", set); err != nil {
		t.Fatalf("seed save: %v", err)
	}
}

func TestCmdShareOutput(t *testing.T) {
	st := newShareStore(t)
	seedShare(t, st)

	var buf bytes.Buffer
	if err := CmdShare(st, "shareproj", "pass", &buf); err != nil {
		t.Fatalf("CmdShare: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "# envchain-share project=shareproj") {
		t.Errorf("missing header, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=") || !strings.Contains(out, "BAZ=") {
		t.Errorf("missing variable lines, got:\n%s", out)
	}
}

func TestCmdShareEmptyProject(t *testing.T) {
	st := newShareStore(t)
	set := env.NewSet()
	_ = st.Save("empty", "pass", set)

	var buf bytes.Buffer
	err := CmdShare(st, "empty", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdShareNotFound(t *testing.T) {
	st := newShareStore(t)
	var buf bytes.Buffer
	err := CmdShare(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdReceiveRoundTrip(t *testing.T) {
	src := newShareStore(t)
	seedShare(t, src)

	var blob bytes.Buffer
	if err := CmdShare(src, "shareproj", "pass", &blob); err != nil {
		t.Fatalf("CmdShare: %v", err)
	}

	dst := newShareStore(t)
	var out bytes.Buffer
	if err := CmdReceive(dst, &blob, "", "newpass", &out); err != nil {
		t.Fatalf("CmdReceive: %v", err)
	}

	set, err := dst.Load("shareproj", "newpass")
	if err != nil {
		t.Fatalf("dst Load: %v", err)
	}
	if v, _ := set.Get("FOO"); v != "bar" {
		t.Errorf("FOO = %q, want bar", v)
	}
	if v, _ := set.Get("BAZ"); v != "qux" {
		t.Errorf("BAZ = %q, want qux", v)
	}
}

func TestCmdReceiveOverrideName(t *testing.T) {
	src := newShareStore(t)
	seedShare(t, src)

	var blob bytes.Buffer
	_ = CmdShare(src, "shareproj", "pass", &blob)

	dst := newShareStore(t)
	var out bytes.Buffer
	if err := CmdReceive(dst, &blob, "renamed", "p2", &out); err != nil {
		t.Fatalf("CmdReceive override: %v", err)
	}

	if _, err := dst.Load("renamed", "p2"); err != nil {
		t.Errorf("expected project 'renamed' to exist: %v", err)
	}
}

func TestCmdReceiveNoProjectName(t *testing.T) {
	dst := newShareStore(t)
	var out bytes.Buffer
	err := CmdReceive(dst, strings.NewReader("FOO=bar\n"), "", "p", &out)
	if err == nil {
		t.Fatal("expected error when no project name available")
	}
}
